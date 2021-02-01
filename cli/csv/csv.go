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

package csv

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/hackaio/pk"
	"io"
	"os"
)

var (
	_ pk.BulkReader = (*reader)(nil)
	_ pk.BulkWriter = (*writer)(nil)
)

type reader struct{}

func NewReader() pk.BulkReader {
	return &reader{}
}

func (r *reader) Read(ctx context.Context, fileName string) (res []pk.Account, err error) {
	csvFile, err := os.Open(fileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			err = nil
		}

		acc := pk.Account{
			Name:     line[0],
			UserName: line[1],
			Email:    line[2],
			Password: line[3],
		}

		res = append(res, acc)

	}

	return res, err
}

type writer struct{}

func NewWriter() pk.BulkWriter {
	return &writer{}
}

func (w *writer) Write(ctx context.Context, request pk.FileWriterReq) error {

	fileName := request.FileName
	fileNameExt := fmt.Sprintf("%v.csv", fileName)
	file, err := os.Create(fileNameExt)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	as := request.Accounts

	for _, acc := range as {
		var record []string

		record = []string{acc.Name, acc.UserName, acc.Email, acc.Password, acc.Created}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return err
}
