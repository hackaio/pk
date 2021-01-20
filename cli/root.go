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
	"github.com/hackaio/pk"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)


var cfgFile string
var appHome string
var appDBDir string
var appCredDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pk",
	Short: "A simple commandline tool to store password on your laptop",
	Long: `pk is a really dump tool it does no magic, It just store passwords
and helps you retrieve them later easily`,
}

// Execute adds all child CLI to the root command and sets flags apkropriately.
// This is called by main.main(). It only needs to hapken once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	//initialize all commands
	configureCommands()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("name", "n", "", "name of account")
	rootCmd.PersistentFlags().StringP("username", "u", "", "account username")
	rootCmd.PersistentFlags().StringP("email", "e", "", "email used")
	rootCmd.PersistentFlags().StringP("password", "p", "", "password")
	rootCmd.PersistentFlags().StringP("token", "t", "", "login token")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pk.yaml)")

	rootCmd.AddCommand(initCmd, addCmd, getCmd, listCmd, updateCmd, deleteCmd, loginCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()

		//create application home directory at $HOME/pk
		appHome = filepath.Join(home, pk.AppDir)
		err = os.MkdirAll(appHome, 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//create databases directory at $HOME/pk/db
		appDBDir = filepath.Join(appHome, pk.DBDir)
		err = os.MkdirAll(appDBDir, 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}


		//create credentials directory at $HOME/pk/creds
		appCredDir = filepath.Join(appHome, pk.CredDir)
		err = os.MkdirAll(appDBDir, 0777)
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
