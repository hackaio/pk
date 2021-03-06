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

package keyring

import (
	"github.com/hackaio/pk"
	"github.com/zalando/go-keyring"
)

type secrets struct{}

var _ pk.SecretsRepository = (*secrets)(nil)

func New() pk.SecretsRepository {
	return &secrets{}
}

func (s *secrets) Set(service, user, password string) error {
	return keyring.Set(service, user, password)
}

func (s *secrets) Get(service, user string) (string, error) {
	return keyring.Get(service, user)
}

func (s *secrets) Delete(service, user string) error {
	return keyring.Delete(service, user)
}
