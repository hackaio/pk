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
package cli

import (
	"fmt"
	"github.com/hackaio/pk/sqlite"
	"github.com/spf13/cobra"
)

var (
	initCmd   *cobra.Command
	addCmd    *cobra.Command
	getCmd    *cobra.Command
	listCmd   *cobra.Command
	updateCmd *cobra.Command
	deleteCmd *cobra.Command
	loginCmd  *cobra.Command
)

func configureCommands() {

	_, err := sqlite.NewStore(appDBDir)
	if err != nil {
		panic(err)
	}



	initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize pk",
		Long:  `this command should be run the first time to set up pk on your machine`,

		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete account",
		Long:  `delete an account that was previously recorded.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
		},
	}

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "add new account",
		Long:  `add new account details`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get called")
		},
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "get details of account",
		Long:  `get all details of account`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get called")
		},
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "update username/password",
		Long:  `updates details of the account. You can only update password and username`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("update called")
		},
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list all accounts",
		Long:  `list all details of all accounts`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get called")
		},
	}

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login to the app",
		Long:  `login to the platform. generate token.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("update called")
		},
	}

}
