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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Thearas/dorisdump/src"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean local data",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initConfig(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cleanFile(GlobalConfig.DataDir); err != nil {
			return err
		}
		if err := cleanFile(GlobalConfig.OutputDir); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func cleanFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	yes := src.Confirm(fmt.Sprintf("Delete %s", absPath))
	if yes && !GlobalConfig.DryRun {
		err = os.RemoveAll(absPath)
	}
	return err
}
