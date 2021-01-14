/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/hackaio/pp"
	"github.com/hackaio/pp/bcrypt"
	"github.com/hackaio/pp/cli"
	"github.com/hackaio/pp/memstore"
)

func main() {

	hasher := bcrypt.New()
	store,err := memstore.New()
	if err != nil {
		panic(err)
	}

	ppInstance := pp.NewInstance(store,hasher)
	ppInstance = pp.AuthMiddleware(ppInstance)
	//ppInstance.Add()
	cli.Execute()
}
