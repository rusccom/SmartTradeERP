package storage

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const unsignedPayload = "UNSIGNED-PAYLOAD"

func (s *R2Store) PresignPut(input PresignPutRequest) (PresignedPut, error) {
	if input.Key == "" || input.ContentType == "" || input.Expires <= 0 {
		return PresignedPut{}, errors.New("invalid presign request")
	}
	req, err := s.newPresignRequest(input)
	if err != nil {
		return PresignedPut{}, err
	}
	now := time.Now().UTC()
	signR2Query(req, s.options, input.Expires, now)
	return PresignedPut{
		URL:       req.URL.String(),
		Method:    http.MethodPut,
		Headers:   map[string]string{"Content-Type": input.ContentType},
		ExpiresAt: now.Add(input.Expires),
	}, nil
}

func (s *R2Store) newPresignRequest(input PresignPutRequest) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPut, s.objectURL(input.Key), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", input.ContentType)
	return req, nil
}

func signR2Query(req *http.Request, options R2Options, expires time.Duration, now time.Time) {
	date := now.Format("20060102")
	headers := signedHeaderNames(req)
	scope := date + "/" + r2Region + "/" + r2Service + "/aws4_request"
	values := req.URL.Query()
	values.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	values.Set("X-Amz-Credential", options.AccessKeyID+"/"+scope)
	values.Set("X-Amz-Date", now.Format("20060102T150405Z"))
	values.Set("X-Amz-Expires", strconv.Itoa(int(expires.Seconds())))
	values.Set("X-Amz-SignedHeaders", strings.Join(headers, ";"))
	req.URL.RawQuery = values.Encode()
	signature := r2QuerySignature(req, options.SecretKey, date, scope, headers)
	values.Set("X-Amz-Signature", signature)
	req.URL.RawQuery = values.Encode()
}

func r2QuerySignature(
	req *http.Request,
	secret string,
	date string,
	scope string,
	headers []string,
) string {
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		req.URL.Query().Get("X-Amz-Date"),
		scope,
		hashHex([]byte(canonicalRequest(req, headers, unsignedPayload))),
	}, "\n")
	return hex.EncodeToString(hmacSHA(signingKey(secret, date), stringToSign))
}
