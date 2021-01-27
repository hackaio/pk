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


import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	_ BulkReader = (*jsonReader)(nil)
	
	_ BulkReader = (*csvReader)(nil)
)

type BulkReader interface {
	Read(ctx context.Context, fileName string)(res []Account,err error)
}

type jsonReader struct {}

func NewJsonReader() BulkReader {
	return &jsonReader{}
}

type accounts struct {
	Accounts []Account `json:"accounts"`
}

func (j jsonReader) Read(ctx context.Context, fileName string) (res []Account, err error) {
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		return res, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'accounts' which we defined above
	err = json.Unmarshal(byteValue, &res)

	if err != nil{
		return res, err
	}

	return res,err
}

type csvReader struct {}

func NewCsvReader() BulkReader {
	return &csvReader{}
}

func (c csvReader) Read(ctx context.Context, fileName string) (res []Account, err error) {
	panic("implement me")
}




