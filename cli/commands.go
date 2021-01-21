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
package cli

import (
	"context"
	"github.com/hackaio/pk"
	"github.com/spf13/cobra"
	"os"
)

type Command int

const (
	Init Command = iota
	Login
	Add
	Get
	Delete
	List
	Update
)




type commander struct {
	keeper pk.PasswordKeeper
}

func (comm *commander) RunCommand(command Command) func(cmd *cobra.Command, args []string) {
	switch command {

	case Init:
		return func(cmd *cobra.Command, args []string) {
			username,err := cmd.Flags().GetString("username")
			email,err := cmd.Flags().GetString("email")
			password,err := cmd.Flags().GetString("password")

			if err != nil{
				logError(err)
				os.Exit(1)
			}

			if username == "" || email == "" || password == ""{
				logUsage(cmd.Example)
				os.Exit(1)
			}
			request := pk.RegisterRequest{
				Username: username,
				Email:    email,
				Password: password,
			}
			errResponse := comm.keeper.Register(context.Background(), request)

			if errResponse.Err != nil {
				logError(errResponse.Err)
				return
			}

			logOK()
			return
		}

	case Login:
		return func(cmd *cobra.Command, args []string) {
			username,err := cmd.Flags().GetString("username")
			password,err := cmd.Flags().GetString("password")

			if err != nil{
				logError(err)
				os.Exit(1)
			}

			if username == "" || password == ""{
				logUsage(cmd.Example)
				os.Exit(1)
			}
			request := pk.LoginRequest{
				UserName: username,
				Password: password,
			}
			response := comm.keeper.Login(context.Background(),request)
			if response.Err != nil {
				os.Exit(1)
			}

			logMessage("token",response.Token)
			return
		}

	case Add:
		return func(cmd *cobra.Command, args []string) {
			username,err := cmd.Flags().GetString("username")
			email,err := cmd.Flags().GetString("email")
			password,err := cmd.Flags().GetString("password")
			name,err := cmd.Flags().GetString("name")
			token,err := cmd.Flags().GetString("token")

			if err != nil {
				logError(err)
				return
			}

			if username == "" || email == "" || password == ""||
				name == "" || token == ""{
				logUsage(cmd.Example)
				os.Exit(1)
			}

			request := pk.AddRequest{
				Token:    token,
				Name:     name,
				UserName: username,
				Email:    email,
				Password: password,
			}

			res := comm.keeper.Add(context.Background(),request)

			if res.Err != nil {
				logError(res.Err)
				return
			}

			logOK()
			return
		}

	case Get:
		return func(cmd *cobra.Command, args []string) {
			username,err := cmd.Flags().GetString("username")
			name,err := cmd.Flags().GetString("name")
			token,err := cmd.Flags().GetString("token")

			if err != nil{
				logError(err)
				os.Exit(1)
			}

			if username == "" || name == "" || token == ""{
				logUsage(cmd.Example)
				os.Exit(1)
			}
			request := pk.GetRequest{
				Token:    token,
				Name:     name,
				UserName: username,
			}

			response := comm.keeper.Get(context.Background(),request)

			if response.Err != nil {
				logError(response.Err)
			}

			logMessage("username",response.Email)
			logMessage("password",response.Password)

		}

	default:
		return func(cmd *cobra.Command, args []string) {
			logUsage("this should not happen")
		}
	}
}

func NewInitCommand(comm commander)*cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize pk",
		Long: `init should be run firstly before anything after installation`,
		Example: "pk init --username <username> --email <email> --password <password>",
		Run: comm.RunCommand(Init),
	}
	
	return initCmd
}

func NewLoginCommand(comm commander)*cobra.Command{
	// loginCmd represents the login command
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "generate auth token",
		Example: "pk login -u <username> -p <password>",
		Long: `generates a jwt token string after the user has supplied username along side master password`,
		Run: comm.RunCommand(Login),
	}

	return loginCmd

}

func NewAddCmd(comm commander) *cobra.Command {
	// addCmd represents the add command
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "add new details to db",
		Example: "pk add <token> <name> <username> <email> <password>",
		Long: `provide name,username,email and password to add new acc`,
		Run: comm.RunCommand(Add),
	}

	return addCmd

}

func NewGetCommand(comm commander)*cobra.Command {
	// getCmd represents the get command
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "get account details",
		Example: "pk get -n <name> -u <username>",
		Long: `return the details of a single account with specified username`,
		Run: comm.RunCommand(Get),
	}

	return getCmd
}

