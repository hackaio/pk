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

package test

import (
	"context"
	"github.com/hackaio/pk"
	"reflect"
	"testing"
)

func TestJSONReader(t *testing.T) {
	var tests []struct {
		name     string
		filename string
		want     interface{}
		wantErr  bool
	}

	tests = append(tests, struct {
		name     string
		filename string
		want     interface{}
		wantErr  bool
	}{name: "test number of accounts", filename: "accounts.json", want:1000 , wantErr: false})

	for index, tt := range tests {

		//test length of array loaded from accounts.json
		if index == 0 {
			t.Run(tt.name, func(t *testing.T) {
				c := pk.NewJsonReader()
				gotRes, err := c.Read(context.Background(), tt.filename)
				if (err != nil) != tt.wantErr {
					t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(len(gotRes), tt.want) {
					t.Errorf("Read() gotRes = %v, want %v", gotRes, tt.want)
				}
			})
		}

	}
}
