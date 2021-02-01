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

package json

import (
	"context"
	cli2 "github.com/hackaio/pk/cli"
	"github.com/hackaio/pk/v0.1.0/pk"
)

var (
	_ cli2.BulkReader = (*reader)(nil)
	_ cli2.BulkWriter = (*writer)(nil)
)

type reader struct {}

func NewReader() cli2.BulkReader {
	return &reader{}
}

func (r *reader) Read(ctx context.Context, fileName string) (res []pk.Account, err error) {
	panic("implement me")
}

type writer struct {}

func NewWriter() cli2.BulkWriter {
	return &writer{}
}

func (w *writer) Write(ctx context.Context, request cli2.FileWriterReq) error {
	panic("implement me")
}

