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

package jwt

import (
	"fmt"
	"github.com/hackaio/pk"
	"testing"
)

func TestIssueToken(t *testing.T) {
	to := NewTokenizer("pk-jwt-tokenizer")
	token := pk.NewToken("piusalfred")
	tStr, err := to.Issue(token)
	if err != nil {
		panic(err)
	}

	fmt.Printf("token : %v\n", tStr)

	tokenRecovered, err := to.Parse(tStr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("id : %v issuedAt: %v", tokenRecovered.ID, tokenRecovered.IssuedAt)
}
