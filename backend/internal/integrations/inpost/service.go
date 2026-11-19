package inpost

import (
	"context"
	"crypto/cipher"
	"net/http"
)

type service struct {
	repo   repository
	client *http.Client
	cipher cipher.AEAD
}

func (s service) save(ctx context.Context, tenantID string, u update) (*connection, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}
	token, err := s.credential(ctx, tenantID, u)
	if err != nil {
		return nil, err
	}
	org, err := verifyOrganization(ctx, s.client, u.account, token.plain)
	if err != nil {
		return nil, err
	}
	c := &connection{account: u.account, OrganizationName: org.Name, Services: org.Services, Revision: u.Revision}
	return c, s.repo.save(ctx, tenantID, c, token.encrypted)
}

func (s service) credential(ctx context.Context, tenantID string, u update) (credential, error) {
	if s.cipher == nil {
		return credential{}, errVault
	}
	binding := []byte("inpost:" + tenantID + ":" + u.Environment + ":" + u.organizationID())
	if u.Token != "" {
		return sealToken(s.cipher, u.Token, binding)
	}
	encrypted, err := s.repo.token(ctx, tenantID, u)
	if err != nil {
		return credential{}, err
	}
	return openToken(s.cipher, encrypted, binding)
}
