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
	oldpk "github.com/hackaio/pk"
	pk "github.com/hackaio/pk/v0.1.0/pk"
	"github.com/hackaio/pk/v0.1.0/pk/cli/commands"
	"github.com/spf13/cobra"
)

var _ commands.Runner = (*commander)(nil)


type commander struct {
	keeper      pk.PasswordKeeper
	credentials oldpk.CredStore
	reader      BulkReader
	writer      BulkWriter
}

func NewCommander(keeper pk.PasswordKeeper, store oldpk.CredStore,
	reader BulkReader, writer BulkWriter)commands.Runner  {
	return &commander{
		keeper:      keeper,
		credentials: store,
		reader:      reader,
		writer:      writer,
	}
}

func (c *commander) Run(command commands.Command) commands.RunFunc {

	switch command {
	
	case commands.Register:
		
		return c.runRegisterCommand()

	case commands.Init:
		return c.runInitCommand()

	case commands.Login:
		return c.runLoginCommand()

	case commands.Add:
		return c.runAddCommand()

	case commands.Get:
		return c.runGetCommand()

	case commands.Delete:
		return c.runDeleteCommand()

	case commands.List:
		return c.runListCommand()

	case commands.Update:
		return c.runUpdateCommand()

	case commands.DB:
		return c.runDBCommand()

	default:
		return func(cmd *cobra.Command, args []string) {
			
		}
		
	}
}

func (c *commander) runInitCommand() commands.RunFunc {
	
}

func (c *commander) runRegisterCommand() commands.RunFunc {

}

func (c *commander) runDBCommand() commands.RunFunc {
	
}

func (c *commander) runUpdateCommand() commands.RunFunc {
	
}

func (c *commander) runListCommand() commands.RunFunc {
	
}

func (c *commander) runDeleteCommand() commands.RunFunc {
	
}

func (c *commander) runLoginCommand() commands.RunFunc {
	
}

func (c *commander) runGetCommand() commands.RunFunc {
	
}

func (c *commander) runAddCommand() commands.RunFunc {
	
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

func makeListCommand(comm commander) *cobra.Command {

}

func makeDBCommand(comm commander) *cobra.Command {

}

func makeUpdateCommand(comm commander) *cobra.Command {

}

func makeDeleteCommand(comm commander) *cobra.Command {

}

func makeGetCommand(comm commander) *cobra.Command {

}

func makeAddCommand(comm commander) *cobra.Command {

}

func makeRegisterCommand(comm commander) *cobra.Command {

}

func makeLoginCommand(comm commander) *cobra.Command {

}

func makeInitCommand(comm commander) *cobra.Command {

}
