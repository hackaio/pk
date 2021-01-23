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

package credstore

import (
	"log"
	"testing"
)

func TestSetGet(t *testing.T) {
	url := "https://github.com/hackaio/pk"
	credStore := New()
	err := credStore.Set("pk-password-manager", url, "user", "password")

	if err != nil {
		panic(err)
	}
	user, secret, err := credStore.Get("pk-password-manager", url)
	if err == nil {
		if user != "user" {
			t.Errorf("Expecting user, got %s", user)
		}

		if secret != "password" {
			t.Errorf("Expecting password, got %s", secret)
		}
	} else {
		log.Println("got error:", err)
	}
}

