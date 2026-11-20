package personal

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/integrations/allegro/connection"
)

type deviceReply struct {
	Code     string  `json:"device_code"`
	URL      string  `json:"verification_uri_complete"`
	Interval seconds `json:"interval"`
	Lifetime seconds `json:"expires_in"`
}

type seconds int

func (s *seconds) UnmarshalJSON(data []byte) error {
	value, err := strconv.Atoi(strings.Trim(string(data), "\""))
	*s = seconds(value)
	return err
}

type invitation struct {
	ID       string    `json:"id"`
	URL      string    `json:"verification_url"`
	Interval int       `json:"interval"`
	Expiry   time.Time `json:"expires_at"`
}

type pendingGrant struct {
	Code string                 `json:"code,omitempty"`
	Pair *connection.Credential `json:"pair,omitempty"`
}

func (s *Service) requestCode(ctx context.Context) (deviceReply, error) {
	var reply deviceReply
	form := url.Values{"client_id": {s.cfg.ClientID}, "scope": {connection.Scopes}}
	if err := s.app.Send(ctx, "device", form, true, &reply); err != nil {
		return reply, err
	}
	address, err := url.Parse(reply.URL)
	if err != nil || address.Scheme != "https" || address.Host != "allegro.pl" || address.User != nil || reply.Code == "" {
		return reply, connection.Remote
	}
	if reply.Interval < 1 || reply.Interval > 300 || reply.Lifetime <= reply.Interval || reply.Lifetime > 86400 {
		return reply, connection.Remote
	}
	return reply, nil
}

func (s *Service) start(ctx context.Context, who actor, version string) (invitation, error) {
	var next invitation
	err := s.DB.WithTenantTx(ctx, who.tenant, func(tx pgx.Tx) error {
		if err := s.checkTarget(ctx, tx, who, version); err != nil {
			return err
		}
		reply, err := s.requestCode(ctx)
		if err != nil {
			return err
		}
		next = invitation{uuid.NewString(), reply.URL, int(reply.Interval), time.Now().UTC().Add(time.Duration(reply.Lifetime) * time.Second)}
		return s.recordDevice(ctx, tx, deviceRequest{who: who, id: next.ID, version: version, expiry: next.Expiry, interval: next.Interval}, reply.Code)
	})
	return next, err
}

func (s *Service) deviceBinding(req deviceRequest) string {
	return "allegro:device:" + req.who.tenant + ":" + req.who.owner + ":" + s.cfg.ClientID + ":" + s.cfg.SellerID + ":" + req.id
}

const insertDevice = `INSERT INTO platform.allegro_device_requests
	(tenant_id,owner_id,request_id,connection_version,encrypted_code,expires_at,poll_after,poll_seconds)
	VALUES ($1,$2,$3,$4,$5,$6,now()+$7::int*interval '1 second',$7)
	ON CONFLICT (tenant_id) DO UPDATE SET owner_id=EXCLUDED.owner_id,request_id=EXCLUDED.request_id,
	connection_version=EXCLUDED.connection_version,encrypted_code=EXCLUDED.encrypted_code,expires_at=EXCLUDED.expires_at,
	poll_after=EXCLUDED.poll_after,poll_seconds=EXCLUDED.poll_seconds`

func (s *Service) recordDevice(ctx context.Context, tx pgx.Tx, req deviceRequest, code string) error {
	data, err := s.Vault.Encrypt(s.deviceBinding(req), pendingGrant{Code: code})
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insertDevice, req.who.tenant, req.who.owner, req.id, req.version, data, req.expiry, req.interval)
	return err
}

func Cancel(ctx context.Context, tx pgx.Tx, tenant string) error {
	_, err := tx.Exec(ctx, `DELETE FROM platform.allegro_device_requests WHERE tenant_id=$1`, tenant)
	return err
}
