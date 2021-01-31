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
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	_ BulkReader = (*jsonReaderWriter)(nil)
	
	_ BulkReader = (*csvReaderWriter)(nil)

	_ BulkWriter = (*jsonReaderWriter)(nil)

	_ BulkWriter = (*csvReaderWriter)(nil)
)

type FileWriterReq struct {
	Accounts []Account
	FileName string
	FileExt string
	FileDir string
}

type BulkReader interface {
	Read(ctx context.Context, fileName string)(res []Account,err error)
}

type BulkWriter interface {
	Write(ctx context.Context, request FileWriterReq) error
}

type BulkReaderWriter interface {
	BulkReader
	BulkWriter
}

type jsonReaderWriter struct {}

func (j jsonReaderWriter) Write(ctx context.Context, request FileWriterReq) error {
	panic("implement me")
}

func JsonReaderWriter() BulkReaderWriter {
	return &jsonReaderWriter{}
}

type accounts struct {
	Accounts []Account `json:"accounts"`
}

func (j jsonReaderWriter) Read(ctx context.Context, fileName string) (res []Account, err error) {
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

type csvReaderWriter struct {
	reader *csv.Reader
	writer *csv.Writer
}

func CSVReaderWriter() BulkReaderWriter {
	return &csvReaderWriter{}
}

func (c csvReaderWriter) Read(ctx context.Context, fileName string) (res []Account, err error) {
	csvFile,err := os.Open(fileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			err = nil
		}

		acc := Account{
			Name:     line[0],
			UserName: line[1],
			Email:    line[2],
			Password: line[3],
		}

		res = append(res,acc)

	}

	return res,err
}


func (c csvReaderWriter) Write(ctx context.Context, request FileWriterReq) error {

	fileName := request.FileName
	fileNameExt := fmt.Sprintf("%v.csv",fileName)
	file, err := os.Create(fileNameExt)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	as := request.Accounts

	for _, acc := range as{
		var record []string

		record = []string{acc.Name, acc.UserName, acc.Email, acc.Password, acc.Created}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return err

}





