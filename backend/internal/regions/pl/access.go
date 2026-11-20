package pl

import (
	"context"
	"strconv"
	"sync"
	"time"

	"smarterp/backend/internal/regions/pl/ksef"
	"smarterp/backend/internal/shared/tenant"
)

type cached struct {
	session *ksef.Session
	until   time.Time
}
type sessionCache struct {
	mu      sync.Mutex
	entries map[string]cached
}

func (c *sessionCache) get(key string) *ksef.Session {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := c.entries[key]
	if time.Now().Before(entry.until) {
		return entry.session
	}
	return nil
}
func (c *sessionCache) put(key string, s *ksef.Session) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = map[string]cached{}
	}
	now := time.Now()
	for k, v := range c.entries {
		if !now.Before(v.until) {
			delete(c.entries, k)
		}
	}
	c.entries[key] = cached{s, now.Add(2 * time.Minute)}
}

func (h *Handler) session(ctx context.Context, d *delivery) (*ksef.Session, error) {
	var p *Seller
	var sealed []byte
	err := h.run(ctx, func(s scope) error {
		var err error
		p, err = readSeller(s.ctx, s.tx, s.tenant)
		if err != nil {
			return err
		}
		if p == nil || p.NIP != d.SellerNIP || p.Environment != d.Environment || !p.HasToken {
			return problem("The configured KSeF account must match this invoice.")
		}
		return s.tx.QueryRow(s.ctx, `SELECT encrypted_token FROM region_pl.sellers WHERE tenant_id=$1 AND revision=$2`, s.tenant, p.Revision).Scan(&sealed)
	})
	if err != nil {
		return nil, err
	}
	return h.openAccess(ctx, *p, sealed)
}

func (h *Handler) openAccess(ctx context.Context, p Seller, sealed []byte) (*ksef.Session, error) {
	tid := tenant.FromContext(ctx)
	key := tid + ":" + p.Environment + ":" + strconv.FormatInt(p.Revision, 10)
	if session := h.sessions.get(key); session != nil {
		return session, nil
	}
	secret, err := h.decrypt(tid, p, sealed)
	if err != nil {
		return nil, err
	}
	client, err := ksef.New(p.Environment)
	if err != nil {
		return nil, err
	}
	session, err := client.Authenticate(ctx, p.NIP, secret)
	if err == nil {
		h.sessions.put(key, session)
	}
	return session, err
}
