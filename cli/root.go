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
	"fmt"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/bcrypt"
	"github.com/hackaio/pk/credstore"
	"github.com/hackaio/pk/pg"
	"github.com/hackaio/pk/rsa"
	"github.com/spf13/cobra"
	"os"

	jwt "github.com/hackaio/pk/jwt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var verboseResp bool
var tokenStr string
var home string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pk",
	Short: "A simple command line tool to store and retrieve passwords",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cli *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pk.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verboseResp,"verbose", "v", false, "verbose command output")
	rootCmd.PersistentFlags().StringVarP(&tokenStr,"token","t","","auth token")
	rootCmd.PersistentFlags().StringP("name","n","","name of the account (e.g github)")
	rootCmd.PersistentFlags().StringP("username","u","","username of the account (e.g alicebob)")
	rootCmd.PersistentFlags().StringP("email","e","","email of the account")
	rootCmd.PersistentFlags().StringP("password","p","","the account password")

	pgDatabase,err := pg.Connect()
	if err != nil {
		panic(err)
	}

	store := pg.NewStore(pgDatabase)
	hasher := bcrypt.New()
	tokenizer := jwt.NewTokenizer("pk")
	/*logger := log.New(os.Stdout,"pk",0)
	logMiddleware := pk.LoggingMiddleware(logger)*/

	homeDir, err := homedir.Dir()
	es,err := rsa.NewEncoderSigner(homeDir)
	if err != nil {
		panic(err)
	}

	cs := credstore.New()

	keeper := pk.NewPasswordKeeper(hasher,store,tokenizer,es,cs)

	comm := commander{keeper: keeper}

	commands := MakeAllCommands(comm)

	rootCmd.AddCommand(
		commands.Init,
		commands.Update,
		commands.Delete,
		commands.Get,
		commands.Login,
		commands.List,
		commands.DB,
		commands.Add,
		)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pk" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pk")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
