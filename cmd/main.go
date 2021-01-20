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
	"fmt"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/cli"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

func main() {

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("can not find home: %v\n",err)
		os.Exit(1)
	}

	appHomePath := filepath.Join(home,pk.AppDir)
	appCredPath := filepath.Join(appHomePath,pk.CredDir)
	appDBPath := filepath.Join(appHomePath,pk.DBDir)
	err = os.MkdirAll(appCredPath, 0777)
	err = os.MkdirAll(appDBPath, 0777)
	if err != nil {
		fmt.Printf("can not create dir: %v\n",err)
		os.Exit(1)
	}

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
