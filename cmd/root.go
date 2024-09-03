/*
Copyright © 2024 Thearas thearas850@gmail.com

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
package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/Thearas/dorisdump/src"
)

var (
	GlobalConfig = Global{}
	completionDB *sqlx.DB
)

type Global struct {
	ConfigFile string
	LogLevel   string
	DataDir    string
	OutputDir  string
	DryRun     bool
	Parallel   int

	DBHost     string
	DBPort     int16
	DBUser     string
	DBPassword string
	DBs        []string
	Tables     []string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dorisdump",
	Short: "Dump schema and query for Doris",
	Long: `
Dump schema and query for Doris.

You may want to pass config by '$HOME/.dorisdump.yaml',
or environment variables with prefix 'DORIS_', e.g.
    DORIS_HOST=xxx
    DORIS_PORT=9030
	`,
	Example:          "dorisdump dump --help",
	SuggestFor:       []string{"dump"},
	ValidArgs:        []string{"completion", "help", "clean", "dump", "anonymize"},
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initConfig(cmd)
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringVar(&GlobalConfig.ConfigFile, "config", "", "Config file (default is $HOME/.dorisdump.yaml)")
	pFlags.StringVarP(&GlobalConfig.LogLevel, "log-level", "L", "info", "Log level, one of: trace, debug, info, warn")
	pFlags.StringVar(&GlobalConfig.DataDir, "data-dir", "./.dorisdump/", "Directory for storing data")
	pFlags.StringVarP(&GlobalConfig.OutputDir, "output", "O", "./dorisdump_output/", "Directory for storing dump sql")
	pFlags.BoolVar(&GlobalConfig.DryRun, "dry-run", false, "Dry run")
	pFlags.IntVar(&GlobalConfig.Parallel, "parallel", 10, "Parallel dump worker")

	pFlags.StringVarP(&GlobalConfig.DBHost, "host", "H", "127.0.0.1", "DB Host")
	pFlags.Int16VarP(&GlobalConfig.DBPort, "port", "P", 9030, "DB Port")
	pFlags.StringVarP(&GlobalConfig.DBUser, "user", "U", "root", "DB User")
	pFlags.StringVar(&GlobalConfig.DBPassword, "password", "", "DB password")
	pFlags.StringSliceVarP(&GlobalConfig.DBs, "dbs", "D", []string{}, "DBs to work on")
	pFlags.StringSliceVarP(&GlobalConfig.Tables, "tables", "T", []string{}, "Tables to work on")

	// completion
	rootCmd.RegisterFlagCompletionFunc("dbs", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var (
			items      = []string{}
			splits     = strings.SplitAfterN(toComplete, ",", 2)
			tocomplete = splits[0]
			err        error
		)
		prefix := ""
		if len(splits) > 1 {
			prefix = splits[1] + ","
		}

		db := getCompletionDB()
		if db != nil {
			items, err = src.ShowDatabases(cmd.Context(), db, tocomplete)
		}
		if len(items) == 0 || err != nil {
			items = []string{}
		}
		return lo.Map(items, func(item string, _ int) string {
			return prefix + item
		}), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	cfgFile := GlobalConfig.ConfigFile
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dorisdump" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dorisdump")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("DORIS")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return err
	}

	bindFlags(cmd, viper.GetViper())

	if err := initLog(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "")
	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)

			vals := []any{val}
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				vals = val.([]any)
			}

			flags := cmd.Flags()
			for _, val := range vals {
				flags.Set(f.Name, fmt.Sprintf("%v", val))
			}
		}
	})
}

func initLog() error {
	logLevel, err := logrus.ParseLevel(GlobalConfig.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %s, err: %v", GlobalConfig.LogLevel, err)
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{})
	return nil
}

func getCompletionDB() *sqlx.DB {
	if completionDB != nil {
		return completionDB
	}

	db, err := connectDB("information_schema")
	if err != nil {
		return nil
	}
	return db
}
