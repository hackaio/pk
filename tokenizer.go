/*
 * Copyright © 2021 PIUS ALFRED me.pius1102@gmail.com
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
