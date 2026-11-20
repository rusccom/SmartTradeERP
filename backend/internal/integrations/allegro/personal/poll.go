package personal

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/integrations/allegro/connection"
)

type deviceRequest struct {
	who               actor
	id, version       string
	expiry, pollAfter time.Time
	interval          int
	data              []byte
}

type progress struct {
	Connected  bool `json:"connected"`
	RetryAfter int  `json:"retry_after,omitempty"`
}

func loadDevice(ctx context.Context, tx pgx.Tx, who actor, id string) (deviceRequest, error) {
	req := deviceRequest{who: who, id: id}
	const query = `SELECT connection_version,expires_at,poll_after,poll_seconds,encrypted_code
		FROM platform.allegro_device_requests WHERE tenant_id=$1 AND owner_id=$2 AND request_id=$3 FOR UPDATE`
	err := tx.QueryRow(ctx, query, who.tenant, who.owner, id).Scan(&req.version, &req.expiry, &req.pollAfter, &req.interval, &req.data)
	if errors.Is(err, pgx.ErrNoRows) {
		return req, connection.State
	}
	return req, err
}

func (s *Service) poll(ctx context.Context, who actor, id string) (progress, error) {
	var next progress
	var terminal error
	err := s.DB.WithTenantTx(ctx, who.tenant, func(tx pgx.Tx) error {
		req, err := loadDevice(ctx, tx, who, id)
		if err != nil {
			return err
		}
		next, terminal = s.advance(ctx, tx, req)
		return settlePoll(ctx, tx, who.tenant, next, terminal)
	})
	if err != nil {
		return next, err
	}
	return next, terminal
}

func settlePoll(ctx context.Context, tx pgx.Tx, tenant string, next progress, err error) error {
	if err == nil && !next.Connected {
		return nil
	}
	var problem *connection.Problem
	if err != nil && !errors.As(err, &problem) {
		return err
	}
	return Cancel(ctx, tx, tenant)
}

func (s *Service) advance(ctx context.Context, tx pgx.Tx, req deviceRequest) (progress, error) {
	if !time.Now().Before(req.expiry) {
		return progress{}, connection.State
	}
	if err := s.checkTarget(ctx, tx, req.who, req.version); err != nil {
		return progress{}, err
	}
	if wait := time.Until(req.pollAfter); wait > 0 {
		return progress{RetryAfter: int(wait/time.Second) + 1}, nil
	}
	pair, err := s.pendingPair(ctx, &req)
	if err == nil {
		err = s.savePair(ctx, tx, req.who, pair, req.version)
	}
	if err != nil {
		return s.waitOrFail(ctx, tx, req, err)
	}
	return progress{Connected: true}, nil
}

func (s *Service) pendingPair(ctx context.Context, req *deviceRequest) (connection.Credential, error) {
	var saved pendingGrant
	key := s.deviceBinding(*req)
	if err := s.Vault.Decrypt(key, req.data, &saved); err != nil {
		return connection.Credential{}, err
	}
	if saved.Pair != nil {
		return *saved.Pair, nil
	}
	pair, err := s.exchangeDevice(ctx, saved.Code)
	if err != nil {
		return pair, err
	}
	req.data, err = s.Vault.Encrypt(key, pendingGrant{Pair: &pair})
	return pair, err
}

func (s *Service) exchangeDevice(ctx context.Context, code string) (connection.Credential, error) {
	form := url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:device_code"}, "device_code": {code}}
	return s.app.Grant(ctx, form, true)
}

func (s *Service) waitOrFail(ctx context.Context, tx pgx.Tx, req deviceRequest, cause error) (progress, error) {
	switch cause {
	case connection.Slow, connection.Limit, connection.Remote:
		req.interval = min(req.interval+5, 120)
	case connection.Pending:
	default:
		return progress{}, cause
	}
	const query = `UPDATE platform.allegro_device_requests SET encrypted_code=$2,poll_seconds=$3,
		poll_after=now()+$3::int*interval '1 second' WHERE tenant_id=$1`
	_, err := tx.Exec(ctx, query, req.who.tenant, req.data, req.interval)
	return progress{RetryAfter: req.interval}, err
}
