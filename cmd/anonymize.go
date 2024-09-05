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
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/Thearas/dorisdump/src"
)

var AnonymizeConfig = Anonymize{}

type Anonymize struct {
	// common
	Enabled     bool
	Method      string
	IdMinLength int
	ReserveIds  []string

	// only for anonymize cmd
	File string
}

// anonymizeCmd represents the anonymize command
var anonymizeCmd = &cobra.Command{
	Use:     "anonymize",
	Short:   "Anonymize sqls",
	Example: `echo "select * from table1" | dorisdump anonymize -f -`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initConfig(cmd)
	},
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		var input []byte

		switch AnonymizeConfig.File {
		case "-":
			// read from stdin
			input, err = io.ReadAll(os.Stdin)
		default:
			input, err = os.ReadFile(AnonymizeConfig.File)
		}
		if err != nil {
			return err
		}

		src.SetupAnonymizer(AnonymizeConfig.IdMinLength, AnonymizeConfig.ReserveIds...)

		sql := src.AnonymizeSql(AnonymizeConfig.Method, "", string(input))
		fmt.Println(sql)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(anonymizeCmd)

	pFlags := anonymizeCmd.PersistentFlags()
	addAnonymizeBaseFlags(pFlags, true)
	pFlags.MarkHidden("anonymize")

	flags := anonymizeCmd.Flags()
	flags.StringVarP(&AnonymizeConfig.File, "file", "f", "", "File path to anonymize sqls, '-' for reading from stdin")
	anonymizeCmd.MarkFlagRequired("file")
}

func addAnonymizeBaseFlags(pFlags *pflag.FlagSet, defaultEnabled bool) {
	pFlags.BoolVar(&AnonymizeConfig.Enabled, "anonymize", defaultEnabled, "Anonymize sqls")
	pFlags.IntVar(&AnonymizeConfig.IdMinLength, "anonymize-id-min-length", 3, "Skip anonymization for id which length is less than this value")
	pFlags.StringSliceVar(&AnonymizeConfig.ReserveIds, "anonymize-reserve-ids", nil, "Skip anonymization for these ids, usually database names")
	pFlags.StringVar(&AnonymizeConfig.Method, "anonymize-method", "hash", "Anonymize method, hash only for now")
	pFlags.MarkHidden("anonymize-method")
}
