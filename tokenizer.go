package pk

import "time"

type Token struct {
	ID        string
	IssuerID  string
	Subject   string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func NewToken(id string) Token {
	return Token{
		ID:        id,
		IssuerID:  "pk123456789",
		Subject:   "pk-master-auth",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

// Tokenizer specifies API for encoding and decoding between string and Key.
type Tokenizer interface {
	// Issue converts Token to its string representation.
	Issue(token Token) (string, error)

	// Parse extracts Token data from string token.
	Parse(string) (Token, error)
}
