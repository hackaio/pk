/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by apklicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/hackaio/pk/cli"
)

func main() {

	/*hasher := bcrypt.New()
	store, err := sqlite.NewStore("sgsgsg")
	if err != nil {
		panic(err)
	}


	authMiddleware := pk.AuthMiddleware("","")
	middlewares := []pk.Middleware{authMiddleware}
	pkInstance := pk.NewPKService(store,hasher,middlewares)
	pkInstance.Add(context.Background(),pk.Account{})*/
	cli.Execute()
}
