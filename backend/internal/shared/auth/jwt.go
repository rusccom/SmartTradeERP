package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	TenantID  string `json:"tenant_id,omitempty"`
	Role      string `json:"role"`
	Scope     string `json:"scope"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type TokenService struct {
	secret    []byte
	accessTTL time.Duration
}

func NewTokenService(secret string, accessTTL time.Duration) *TokenService {
	return &TokenService{secret: []byte(secret), accessTTL: accessTTL}
}

func (s *TokenService) Issue(userID, tenantID, role, scope string) (TokenResponse, error) {
	token, err := s.issueToken(userID, tenantID, role, scope, "access", s.accessTTL)
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{AccessToken: token}, nil
}

func (s *TokenService) Parse(tokenRaw string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenRaw, &Claims{}, s.keyFunc, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func (s *TokenService) issueToken(userID, tenantID, role, scope, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{}
	claims.UserID = userID
	claims.TenantID = tenantID
	claims.Role = role
	claims.Scope = scope
	claims.TokenType = tokenType
	claims.RegisteredClaims = buildRegistered(userID, now, ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func buildRegistered(subject string, now time.Time, ttl time.Duration) jwt.RegisteredClaims {
	claims := jwt.RegisteredClaims{}
	claims.Subject = subject
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(ttl))
	return claims
}

func (s *TokenService) keyFunc(_ *jwt.Token) (interface{}, error) {
	return s.secret, nil
}
