package ksef

import (
	"context"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

type ref struct {
	Reference string `json:"referenceNumber"`
}
type Upload struct {
	Reference string
	body      map[string]any
}
type Result struct {
	Reference string `json:"referenceNumber"`
	Hash      string `json:"invoiceHash"`
	Number    string `json:"ksefNumber"`
	Status    Status `json:"status"`
}

func (s *Session) Open(ctx context.Context, xml []byte) (*Upload, error) {
	cert, err := selectCert(s.certs, "SymmetricKeyEncryption")
	if err != nil {
		return nil, err
	}
	m, err := newMaterial()
	if err != nil {
		return nil, err
	}
	wrapped, err := cert.wrap(m.key)
	if err != nil {
		return nil, err
	}
	payload, err := m.invoice(xml)
	if err != nil {
		return nil, err
	}
	reference, err := s.openSession(ctx, wrapped, m.iv, cert.ID)
	return &Upload{reference, payload}, err
}

func (s *Session) openSession(ctx context.Context, key string, iv []byte, id string) (string, error) {
	body := map[string]any{"formCode": map[string]string{"systemCode": "FA (3)", "schemaVersion": "1-0E", "value": "FA"},
		"encryption": map[string]string{"encryptedSymmetricKey": key, "initializationVector": base64.StdEncoding.EncodeToString(iv), "publicKeyId": id}}
	var out ref
	err := s.client.request(ctx, call{"POST", "/sessions/online", s.token, body, &out})
	if err == nil && out.Reference == "" {
		return "", errors.New("KSeF did not open a session")
	}
	return out.Reference, err
}

func (s *Session) Send(ctx context.Context, u *Upload) (string, error) {
	var out ref
	path := "/sessions/online/" + url.PathEscape(u.Reference) + "/invoices"
	err := s.client.request(ctx, call{"POST", path, s.token, u.body, &out})
	if err == nil && out.Reference == "" {
		return "", errors.New("KSeF did not return an invoice reference")
	}
	return out.Reference, err
}

func (s *Session) Close(ctx context.Context, session string) error {
	return s.client.request(ctx, call{method: "POST", path: "/sessions/online/" + url.PathEscape(session) + "/close", token: s.token})
}

func (s *Session) Check(ctx context.Context, session, invoice string, xml []byte) (Result, error) {
	path := "/sessions/" + url.PathEscape(session) + "/invoices"
	if invoice == "" {
		return s.findByHash(ctx, path, xml)
	}
	var out Result
	err := s.client.request(ctx, call{method: "GET", path: path + "/" + url.PathEscape(invoice), token: s.token, out: &out})
	if err == nil && out.Hash != digest(xml) {
		return Result{}, errors.New("invoice hash does not match the submitted document")
	}
	return out, err
}

func (s *Session) findByHash(ctx context.Context, path string, xml []byte) (Result, error) {
	var list struct {
		Invoices []Result `json:"invoices"`
	}
	err := s.client.request(ctx, call{method: "GET", path: path, token: s.token, out: &list})
	if err != nil {
		return Result{}, err
	}
	hash := digest(xml)
	for _, item := range list.Invoices {
		if item.Hash == hash {
			return item, nil
		}
	}
	return s.checkEmpty(ctx, path)
}

func (s *Session) checkEmpty(ctx context.Context, path string) (Result, error) {
	var session struct {
		Status Status `json:"status"`
		Count  *int   `json:"invoiceCount"`
	}
	base := strings.TrimSuffix(path, "/invoices")
	if err := s.client.request(ctx, call{method: "GET", path: base, token: s.token, out: &session}); err != nil {
		return Result{}, err
	}
	if session.Status.Code == 440 || (session.Status.Code == 200 && session.Count != nil && *session.Count == 0) {
		return Result{Status: Status{Code: 405, Description: "KSeF closed this session without receiving an invoice."}}, nil
	}
	return Result{}, errors.New("KSeF has not confirmed this invoice in its session")
}
