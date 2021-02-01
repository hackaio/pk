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
	"github.com/hackaio/pk/cli/commands"
	"os"
	"path/filepath"

	"github.com/hackaio/pk"
	"github.com/hackaio/pk/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	minPasswordLen = 6
)

var (
	debugMessage = "not yet implemented"
)

type commander struct {
	keeper      pk.PasswordKeeper
	credentials CredStore
	csvReader   BulkReader
	csvWriter   BulkWriter
	jsonReader  BulkReader
	jsonWriter  BulkWriter
}

var _ commands.Runner = (*commander)(nil)

func NewCommandsRunner(keeper pk.PasswordKeeper, store CredStore,
	csvReader BulkReader, csvWriter BulkWriter, jsonReader BulkReader,
	jsonWriter BulkWriter) commands.Runner {
	return &commander{
		keeper:      keeper,
		credentials: store,
		csvReader:   csvReader,
		csvWriter:   csvWriter,
		jsonReader:  jsonReader,
		jsonWriter:  jsonWriter,
	}
}

//Commands a struct with all pk commands
type Commands struct {
	Init     *cobra.Command
	Register *cobra.Command
	Login    *cobra.Command
	Add      *cobra.Command
	Get      *cobra.Command
	Delete   *cobra.Command
	Update   *cobra.Command
	DB       *cobra.Command
	List     *cobra.Command
}

func MakeAllCommands(comm commander) Commands {
	return Commands{
		Init:     makeInitCommand(comm),
		Register: makeRegisterCommand(comm),
		Login:    makeLoginCommand(comm),
		Add:      makeAddCommand(comm),
		Get:      makeGetCommand(comm),
		Delete:   makeDeleteCommand(comm),
		Update:   makeUpdateCommand(comm),
		DB:       makeDBCommand(comm),
		List:     makeListCommand(comm),
	}
}

func (comm *commander) runInitCommand() commands.RunFunc {
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

		err = comm.credentials.Set(pk.AppName, username, string(password))

		if err != nil {
			err1 := errors.New(fmt.Sprintf("could not save token due to: %v", err))
			logError(err1)
			os.Exit(1)
		}

		err = comm.keeper.Register(context.Background(), username, email, string(password))

		if err != nil {
			logError(err)
			return
		}

		logOK()
		return
	}
}

func (comm *commander) runLoginCommand() commands.RunFunc {
	return func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")

		if err != nil || username == "" {
			logError(err)
			os.Exit(1)
		}

		var password string

		password, err = comm.credentials.Get(pk.AppName, username)

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
				_ = comm.credentials.Set(pk.AppName, username, password)

				return
			} else {
				logError(err)
				os.Exit(1)
			}
		}

		token, err := comm.keeper.Login(context.Background(), username, password)
		if err != nil {
			logError(err)
			os.Exit(1)
		}

		err = comm.credentials.Set(pk.AppName, "token", token)

		if err != nil {
			err1 := errors.New(fmt.Sprintf("could not save token due to: %v", err))
			logError(err1)
			os.Exit(1)
		}

		logMessage("token", token)
		return
	}

}

func (comm *commander) runAddCommand() commands.RunFunc {
	return func(cmd *cobra.Command, args []string) {

		username, err := cmd.Flags().GetString("username")
		fileName, err := cmd.Flags().GetString("file")
		email, err := cmd.Flags().GetString("email")
		password, err := cmd.Flags().GetString("password")
		name, err := cmd.Flags().GetString("name")
		token, err := comm.credentials.Get(pk.AppName, "token")

		if err != nil {
			logError(err)
			return
		}

		ctx := context.Background()

		fileNameAvailable := len(fileName) > 4

		var accounts []pk.Account

		if fileNameAvailable && len(token) > 1 {

			//check if its csv or json
			ext := filepath.Ext(fileName)
			if ext == ".json" {

				accounts, err := comm.jsonReader.Read(ctx, fileName)

				if err != nil {
					logError(err)
					os.Exit(1)
				}

				err = comm.keeper.AddAll(ctx, token, accounts)

				if err != nil {
					logError(err)
					os.Exit(1)
				}

				logOK()
				return

			} else if ext == ".csv" {

				accounts, err = comm.csvReader.Read(ctx, fileName)

				if err != nil {
					logError(err)
					os.Exit(1)
				}

				err = comm.keeper.AddAll(ctx, token, accounts)

				if err != nil {
					logError(err)
					os.Exit(1)
				}

				logOK()
				return

			} else {
				err1 := errors.New("parse json or csv files only")
				logError(err1)
				os.Exit(1)
			}

		} else if username == "" || email == "" || password == "" ||
			name == "" || token == "" {
			logUsage(cmd.Example)
			os.Exit(1)
		} else {

			account := pk.Account{
				Name:     name,
				UserName: username,
				Email:    email,
				Password: password,
			}

			err = comm.keeper.Add(context.Background(), token, account)

			if err != nil {
				logError(err)
				return
			}

			logOK()
			return
		}
	}
}

func (comm *commander) runGetCommand() commands.RunFunc {

	return func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		name, err := cmd.Flags().GetString("name")
		token, err := comm.credentials.Get(pk.AppName, "token")

		if err != nil {
			logError(err)
			os.Exit(1)
		}

		if username == "" || name == "" || token == "" {
			logUsage(cmd.Example)
			os.Exit(1)
		}

		response, err := comm.keeper.Get(context.Background(), token, name, username)

		if err != nil {
			logError(err)

		}

		logMessage("username", response.Email)
		logMessage("password", response.Password)

	}

}

func (comm *commander) runDeleteCommand() commands.RunFunc {
	return func(cmd *cobra.Command, args []string) {
		logMessage("error", debugMessage)
	}
}

func (comm *commander) runListCommand() commands.RunFunc {

	return func(cmd *cobra.Command, args []string) {
		limit, err := cmd.Flags().GetInt("limit")
		token, err := comm.credentials.Get(pk.AppName, "token")

		if err != nil {
			logError(err)
			os.Exit(1)
		}

		var argx map[string]interface{}
		accounts, err := comm.keeper.List(context.Background(), token, argx)

		if err != nil {
			logError(err)
			os.Exit(1)
		}

		out, err := cmd.Flags().GetString("out")
		format, err := cmd.Flags().GetString("format")
		dir, err := cmd.Flags().GetString("dir")

		/*	fileFormat, err := getUploadFormat(format)

			if err != nil {
				logError(err)
				os.Exit(1)
			}*/

		if limit == 0 {
			limit = len(accounts)
		}

		if out == "" && format == "" && dir == "" {
			logJSON(accounts[:limit])
		} else {
			if out == "" {
				out = "accounts"
			}

			if format == "" {
				format = "json"
			}

			if !(format == "json" || format == "csv") {
				errFormat := errors.New("invalid file format")
				logError(errFormat)
				os.Exit(1)

			}

			if dir == "" {
				path, err := os.Getwd()
				if err != nil {
					logError(err)
					os.Exit(1)
				}
				dir = path
			}

			req := FileWriterReq{
				Accounts: accounts[:limit],
				FileName: out,
				FileExt:  format,
				FileDir:  dir,
			}

			if format == "csv" {

				err := comm.csvWriter.Write(context.Background(), req)
				if err != nil {
					logError(err)
					os.Exit(1)
				}
			} else if format == "json" {

				err := comm.jsonWriter.Write(context.Background(), req)
				if err != nil {
					logError(err)
					os.Exit(1)
				}
			}
		}

	}
}

func (comm *commander) runUpdateCommand() commands.RunFunc {
	return func(cmd *cobra.Command, args []string) {
		logMessage("error", debugMessage)
	}
}

func (comm *commander) runDBCommand() commands.RunFunc {
	return func(cmd *cobra.Command, args []string) {
		logMessage("error", debugMessage)
	}
}

func (comm *commander) Run(command commands.Command) commands.RunFunc {

	switch command {

	case commands.Init:
		return comm.runInitCommand()

	case commands.Login:
		return comm.runLoginCommand()

	case commands.Add:
		return comm.runAddCommand()

	case commands.Get:
		return comm.runGetCommand()

	case commands.Delete:
		return comm.runDeleteCommand()

	case commands.List:
		return comm.runListCommand()

	case commands.Update:
		return comm.runUpdateCommand()

	case commands.DB:
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
		Run:     comm.Run(commands.Init),
	}

	initCmd.Flags().BoolP("credstore", "c", true, "store credentials")

	return initCmd
}

func makeRegisterCommand(comm commander) *cobra.Command {
	var regCmd = &cobra.Command{
		Use:     "register",
		Short:   "set up pk",
		Long:    `init should be run firstly before anything after installation`,
		Example: "pk init -u <username> -e <email>",
		Run:     comm.Run(commands.Register),
	}

	return regCmd
}

func makeLoginCommand(comm commander) *cobra.Command {
	// loginCmd represents the login command
	var loginCmd = &cobra.Command{
		Use:     "login",
		Short:   "generate auth token",
		Example: "pk login -u <username>",
		Long:    `generates a jwt token string after the user has supplied username along side master password`,
		Run:     comm.Run(commands.Login),
	}

	return loginCmd

}

func makeAddCommand(comm commander) *cobra.Command {
	// addCmd represents the add command
	var addCmd = &cobra.Command{
		Use:     "add",
		Short:   "add new details to db",
		Example: "pk add --file accounts.json",
		Long:    `provide name,username,email and password to add new acc`,
		Run:     comm.Run(commands.Add),
	}

	addCmd.PersistentFlags().StringP("file", "f", "", "json or csv accounts file")

	return addCmd

}

func makeGetCommand(comm commander) *cobra.Command {
	// getCmd represents the get command
	var getCmd = &cobra.Command{
		Use:     "get",
		Short:   "get account",
		Example: "pk get -n <name> -u <username>",
		Long:    `return the details of a single account with specified username`,
		Run:     comm.Run(commands.Get),
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
		Run:     comm.Run(commands.List),
	}

	listCmd.PersistentFlags().IntP("limit", "l", 0, "limits of accounts to list")
	listCmd.PersistentFlags().StringP("out", "o", "", "output filename")
	listCmd.PersistentFlags().StringP("format", "m", "", "output file format")
	listCmd.PersistentFlags().StringP("dir", "d", "", "output directory")

	return listCmd
}

func makeDeleteCommand(comm commander) *cobra.Command {
	// deleteCmd represents the get command
	var deleteCmd = &cobra.Command{
		Use:     "delete",
		Short:   "delete account",
		Example: "pk delete -n <name> -u <username>",
		Long:    `delete account details by specifying username and name`,
		Run:     comm.Run(commands.Delete),
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
		Run:     comm.Run(commands.Update),
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
		Run:     comm.Run(commands.DB),
	}

	return dbCmd
}
