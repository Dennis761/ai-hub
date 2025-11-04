package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	adminports "ai_hub.com/app/core/ports/adminports"
	"ai_hub.com/app/infra/config"
	"github.com/golang-jwt/jwt"
	str2dur "github.com/xhit/go-str2duration/v2"
)

var _ adminports.TokenIssuer = (*JwtTokenIssuer)(nil)

type JwtTokenIssuer struct {
	secret         string
	defaultExpires string
}

func NewJwtTokenIssuer(secret, defaultExpires string) *JwtTokenIssuer {
	return &JwtTokenIssuer{
		secret:         strings.TrimSpace(secret),
		defaultExpires: strings.TrimSpace(defaultExpires),
	}
}

func NewJwtTokenIssuerFromConfig() *JwtTokenIssuer {
	return NewJwtTokenIssuer(config.Env.JWTSecret, config.Env.JWTExpiresIn)
}

func (j *JwtTokenIssuer) Issue(payload adminports.TokenPayload, opts *adminports.TokenOptions) (string, error) {
	if j.secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	claims := jwt.MapClaims{
		"userID": payload.UserID,
		"roles":  payload.Roles,
	}
	if payload.Email != nil && *payload.Email != "" {
		claims["email"] = *payload.Email
	}

	expires := j.defaultExpires
	if opts != nil && strings.TrimSpace(opts.ExpiresIn) != "" {
		expires = strings.TrimSpace(opts.ExpiresIn)
	}
	if expires != "" {
		d, err := parseDurationMS(expires)
		if err != nil {
			return "", fmt.Errorf("invalid ExpiresIn %q: %w", expires, err)
		}
		claims["exp"] = time.Now().UTC().Add(d).Unix()
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(j.secret))
}

func (j *JwtTokenIssuer) Verify(tokenStr string) (*adminports.TokenPayload, error) {
	if strings.TrimSpace(tokenStr) == "" {
		return nil, errors.New("empty token")
	}
	if j.secret == "" {
		return nil, errors.New("jwt secret is empty")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, _ := claims["userID"].(string)
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid token payload: userID is empty")
	}

	var emailPtr *string
	if v, ok := claims["email"].(string); ok && strings.TrimSpace(v) != "" {
		email := v
		emailPtr = &email
	}

	var roles []string
	if arr, ok := claims["roles"].([]interface{}); ok {
		for _, it := range arr {
			if s, ok := it.(string); ok {
				roles = append(roles, s)
			}
		}
	} else if arrStr, ok := claims["roles"].([]string); ok {
		roles = append(roles, arrStr...)
	}

	return &adminports.TokenPayload{
		UserID: userID,
		Email:  emailPtr,
		Roles:  roles,
	}, nil
}

func parseDurationMS(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)

	if isNumeric(s) {
		return str2dur.ParseDuration(s + "ms")
	}
	return str2dur.ParseDuration(s)
}

func isNumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		if (s[i] < '0' || s[i] > '9') && !(i == 0 && s[i] == '-') {
			return false
		}
	}
	return len(s) > 0
}
