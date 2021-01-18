package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/pkg/errors"
	"time"
)

const issuerName = "pk.auth"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type claims struct {
	jwt.StandardClaims
	IssuerID string  `json:"issuer_id,omitempty"`
}

func (c claims) Valid() error {
	if c.Issuer != issuerName {
		return ErrInvalidToken
	}

	return c.StandardClaims.Valid()
}

func (c claims) toToken() pk.Token {
	token := pk.Token{
		ID:        c.Id,
		IssuerID:  c.IssuerID,
		Subject:   c.Subject,
		IssuedAt:  time.Unix(c.IssuedAt,0).UTC(),
	}

	if c.ExpiresAt != 0 {
		token.ExpiresAt = time.Unix(c.ExpiresAt, 0).UTC()
	}

	return token
}


type tokenizer struct {
	secret string
}

var _ pk.Tokenizer = (*tokenizer)(nil)

func NewTokenizer(secret string) pk.Tokenizer {
	return &tokenizer{secret: secret}
}

func (t tokenizer) Issue(token pk.Token) (string, error) {
	claims := claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   issuerName,
			Subject:  token.Subject,
			IssuedAt: token.IssuedAt.UTC().Unix(),
		},
		IssuerID: token.IssuerID,
	}

	if !token.ExpiresAt.IsZero() {
		claims.ExpiresAt = token.ExpiresAt.UTC().Unix()
	}
	if token.ID != "" {
		claims.Id = token.ID
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(t.secret))
}

func (t tokenizer) Parse(token string) (pk.Token, error) {
	c := claims{}
	_, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pk.ErrPermissionDenied
		}
		return []byte(t.secret), nil
	})

	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok && e.Errors == jwt.ValidationErrorExpired {

			return pk.Token{}, errors.Wrap(ErrTokenExpired, err)
		}
		return pk.Token{}, errors.Wrap(pk.ErrPermissionDenied, err)
	}

	return c.toToken(), nil

}


