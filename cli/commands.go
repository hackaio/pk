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
	"bytes"
	"context"
	"fmt"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

const (
	minPasswordLen = 6
)

var (
	debugMessage = "not yet implemented"
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
	DB
)

type CommandFunc func(cmd *cobra.Command, args []string)

type commander struct {
	keeper pk.PasswordKeeper
}

type Commands struct {
	Init   *cobra.Command
	Login  *cobra.Command
	Add    *cobra.Command
	Get    *cobra.Command
	Delete *cobra.Command
	Update *cobra.Command
	DB     *cobra.Command
	List   *cobra.Command
}

func MakeAllCommands(comm commander) Commands {
	return Commands{
		Init:   makeInitCommand(comm),
		Login:  makeLoginCommand(comm),
		Add:    makeAddCommand(comm),
		Get:    makeGetCommand(comm),
		Delete: makeDeleteCommand(comm),
		Update: makeUpdateCommand(comm),
		DB:     makeDBCommand(comm),
		List:   makeListCommand(comm),
	}
}

func (comm *commander) runInitCommand() CommandFunc {
	return func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		email, err := cmd.Flags().GetString("email")

		if err != nil {
			logError(err)
			os.Exit(1)
		}

		if username == "" || email == "" {
			err1 := errors.New("username and email not specified")
			logUsage(cmd.Example)
			logError(err1)
			os.Exit(1)
		}

		fmt.Println("Enter password: ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			logError(err)
			os.Exit(1)
		}
		fmt.Println("Enter password again: ")

		password1, err := terminal.ReadPassword(0)

		if err != nil {
			logError(err)
			os.Exit(1)
		}

		if !bytes.Equal(password, password1) {
			err1 := errors.New("password mismatch")
			logError(err1)
			os.Exit(1)
		}

		if string(password) == "" || len(string(password)) < minPasswordLen {
			err1 := errors.New("password length should be >= 6 chars")
			logError(err1)
			os.Exit(1)
		}

		cs := comm.keeper.CredStore()

		err = cs.Set(pk.AppName, username, string(password))

		if err != nil {
			err1 := errors.New(fmt.Sprintf("could not save token due to: %v", err))
			logError(err1)
			os.Exit(1)
		}
		request := pk.RegisterRequest{
			Username: username,
			Email:    email,
			Password: string(password),
		}
		errResponse := comm.keeper.Register(context.Background(), request)

		if errResponse.Err != nil {
			logError(errResponse.Err)
			return
		}

		logOK()
		return
	}
}

func (comm *commander) runLoginCommand() CommandFunc {
	return func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")

		if err != nil || username == "" {
			logError(err)
			os.Exit(1)
		}

		cs := comm.keeper.CredStore()

		var password string

		password, err = cs.Get(pk.AppName, username)

		if err != nil || password == "" {

			if err == keyring.ErrNotFound {
				fmt.Println("Enter password: ")
				passwordBytes, err := terminal.ReadPassword(0)
				if err != nil {
					logError(err)
					os.Exit(1)
				}
				fmt.Println("Enter password again: ")

				passwordBytes1, err := terminal.ReadPassword(0)

				if err != nil {
					logError(err)
					os.Exit(1)
				}

				if !bytes.Equal(passwordBytes, passwordBytes1) {
					err1 := errors.New("password mismatch")
					logError(err1)
					os.Exit(1)
				}

				if string(passwordBytes) == "" || len(string(passwordBytes)) < minPasswordLen {
					err1 := errors.New("password length should be >= 6 chars")
					logError(err1)
					os.Exit(1)
				}

				password = string(passwordBytes)


				//fixme
				_ = cs.Set(pk.AppName, username, password)

				return
			} else {
				logError(err)
				os.Exit(1)
			}
		}

		request := pk.LoginRequest{
			UserName: username,
			Password: password,
		}
		response := comm.keeper.Login(context.Background(), request)
		if response.Err != nil {
			os.Exit(1)
		}

		err = cs.Set(pk.AppName, "token", response.Token)

		if err != nil {
			err1 := errors.New(fmt.Sprintf("could not save token due to: %v", err))
			logError(err1)
			os.Exit(1)
		}

		logMessage("token", response.Token)
		return
	}

}

func (comm *commander) runAddCommand() CommandFunc {
	return func(cmd *cobra.Command, args []string) {

		cs := comm.keeper.CredStore()
		username, err := cmd.Flags().GetString("username")
		email, err := cmd.Flags().GetString("email")
		password, err := cmd.Flags().GetString("password")
		name, err := cmd.Flags().GetString("name")
		token, err := cs.Get(pk.AppName, "token")

		if err != nil {
			logError(err)
			return
		}

		if username == "" || email == "" || password == "" ||
			name == "" || token == "" {
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

		res := comm.keeper.Add(context.Background(), request)

		if res.Err != nil {
			logError(res.Err)
			return
		}

		logOK()
		return
	}

}

func (comm *commander) runGetCommand() CommandFunc {

	cs := comm.keeper.CredStore()
	return func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		name, err := cmd.Flags().GetString("name")
		token, err := cs.Get(pk.AppName, "token")

		if err != nil {
			logError(err)
			os.Exit(1)
		}

		if username == "" || name == "" || token == "" {
			logUsage(cmd.Example)
			os.Exit(1)
		}
		request := pk.GetRequest{
			Token:    token,
			Name:     name,
			UserName: username,
		}

		response := comm.keeper.Get(context.Background(), request)

		if response.Err != nil {
			logError(response.Err)
		}

		logMessage("username", response.Email)
		logMessage("password", response.Password)

	}

}

func (comm *commander) runDeleteCommand() CommandFunc {
	return func(cmd *cobra.Command, args []string) {
		logMessage("error", debugMessage)
	}
}

func (comm *commander) runListCommand() CommandFunc {
	return func(cmd *cobra.Command, args []string) {
		logMessage("error", debugMessage)
	}
}

func (comm *commander) runUpdateCommand() CommandFunc {
	return func(cmd *cobra.Command, args []string) {
		logMessage("error", debugMessage)
	}
}

func (comm *commander) runDBCommand() CommandFunc {
	return func(cmd *cobra.Command, args []string) {
		logMessage("error", debugMessage)
	}
}

func (comm *commander) RunCommand(command Command) CommandFunc {

	switch command {

	case Init:
		return comm.runInitCommand()

	case Login:
		return comm.runLoginCommand()

	case Add:
		return comm.runAddCommand()

	case Get:
		return comm.runGetCommand()

	case Delete:
		return comm.runDeleteCommand()

	case List:
		return comm.runListCommand()

	case Update:
		return comm.runUpdateCommand()

	case DB:
		return comm.runDBCommand()

	default:
		return func(cmd *cobra.Command, args []string) {
			logUsage("this should not happen")
		}
	}
}

func makeInitCommand(comm commander) *cobra.Command {
	var initCmd = &cobra.Command{
		Use:     "init",
		Short:   "set up pk",
		Long:    `init should be run firstly before anything after installation`,
		Example: "pk init -u <username> -e <email>",
		Run:     comm.RunCommand(Init),
	}

	initCmd.Flags().BoolP("credstore", "c", true, "store credentials")

	return initCmd
}

func makeLoginCommand(comm commander) *cobra.Command {
	// loginCmd represents the login command
	var loginCmd = &cobra.Command{
		Use:     "login",
		Short:   "generate auth token",
		Example: "pk login -u <username>",
		Long:    `generates a jwt token string after the user has supplied username along side master password`,
		Run:     comm.RunCommand(Login),
	}

	return loginCmd

}

func makeAddCommand(comm commander) *cobra.Command {
	// addCmd represents the add command
	var addCmd = &cobra.Command{
		Use:     "add",
		Short:   "add new details to db",
		Example: "pk add <token> <name> <username> <email> <password>",
		Long:    `provide name,username,email and password to add new acc`,
		Run:     comm.RunCommand(Add),
	}

	return addCmd

}

func makeGetCommand(comm commander) *cobra.Command {
	// getCmd represents the get command
	var getCmd = &cobra.Command{
		Use:     "get",
		Short:   "get account",
		Example: "pk get -n <name> -u <username>",
		Long:    `return the details of a single account with specified username`,
		Run:     comm.RunCommand(Get),
	}

	return getCmd
}

func makeListCommand(comm commander) *cobra.Command {
	// listCmd represents the get command
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "retrieve all accounts",
		Example: "pk list",
		Long:    `list all accounts details`,
		Run:     comm.RunCommand(List),
	}

	return listCmd
}

func makeDeleteCommand(comm commander) *cobra.Command {
	// deleteCmd represents the get command
	var deleteCmd = &cobra.Command{
		Use:     "delete",
		Short:   "delete account",
		Example: "pk delete -n <name> -u <username>",
		Long:    `delete account details by specifying username and name`,
		Run:     comm.RunCommand(Delete),
	}

	return deleteCmd
}

func makeUpdateCommand(comm commander) *cobra.Command {
	// updateCmd represents the get command
	var updateCmd = &cobra.Command{
		Use:     "update",
		Short:   "update account",
		Example: "pk update -n <name> -u <username>",
		Long:    `update account details by specifying username and name`,
		Run:     comm.RunCommand(Update),
	}

	return updateCmd
}

func makeDBCommand(comm commander) *cobra.Command {
	// dbCmd represents the get command
	var dbCmd = &cobra.Command{
		Use:     "db",
		Short:   "db management command",
		Example: "pk db",
		Long:    `manages pk database`,
		Run:     comm.RunCommand(DB),
	}

	return dbCmd
}
