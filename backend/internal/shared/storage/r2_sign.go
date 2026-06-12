package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sort"
	"strings"
	"time"
)

const amzContentSHA = "X-Amz-Content-Sha256"
const amzDate = "X-Amz-Date"

func signR2Request(req *http.Request, options R2Options, body []byte) {
	now := time.Now().UTC()
	bodyHash := hashHex(body)
	req.Header.Set(amzContentSHA, bodyHash)
	req.Header.Set(amzDate, now.Format("20060102T150405Z"))
	req.Header.Set("Authorization", authorizationHeader(req, options, now, bodyHash))
}

func authorizationHeader(
	req *http.Request,
	options R2Options,
	now time.Time,
	bodyHash string,
) string {
	date := now.Format("20060102")
	signedHeaders := signedHeaderNames(req)
	credentialScope := date + "/" + r2Region + "/" + r2Service + "/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		req.Header.Get(amzDate),
		credentialScope,
		hashHex([]byte(canonicalRequest(req, signedHeaders, bodyHash))),
	}, "\n")
	signature := hex.EncodeToString(hmacSHA(signingKey(options.SecretKey, date), stringToSign))
	return "AWS4-HMAC-SHA256 Credential=" + options.AccessKeyID + "/" +
		credentialScope + ", SignedHeaders=" + strings.Join(signedHeaders, ";") +
		", Signature=" + signature
}

func canonicalRequest(req *http.Request, headers []string, bodyHash string) string {
	return strings.Join([]string{
		req.Method,
		req.URL.EscapedPath(),
		req.URL.RawQuery,
		canonicalHeaders(req, headers),
		strings.Join(headers, ";"),
		bodyHash,
	}, "\n")
}

func signedHeaderNames(req *http.Request) []string {
	names := []string{"host"}
	for name := range req.Header {
		names = append(names, strings.ToLower(name))
	}
	sort.Strings(names)
	return names
}

func canonicalHeaders(req *http.Request, headers []string) string {
	lines := make([]string, 0, len(headers))
	for _, name := range headers {
		value := req.URL.Host
		if name != "host" {
			value = req.Header.Get(name)
		} else if req.Host != "" {
			value = req.Host
		}
		lines = append(lines, name+":"+cleanHeaderValue(value))
	}
	return strings.Join(lines, "\n") + "\n"
}

func signingKey(secret, date string) []byte {
	dateKey := hmacSHA([]byte("AWS4"+secret), date)
	regionKey := hmacSHA(dateKey, r2Region)
	serviceKey := hmacSHA(regionKey, r2Service)
	return hmacSHA(serviceKey, "aws4_request")
}

func hmacSHA(key []byte, value string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))
	return mac.Sum(nil)
}

func hashHex(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func cleanHeaderValue(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
