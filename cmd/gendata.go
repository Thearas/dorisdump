/*
Copyright © 2025 Thearas thearas850@gmail.com

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

	"github.com/Thearas/dorisdump/src" // Import the src package
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
)

// GendataConfig holds the configuration for the gendata command
type GendataConfig struct {
	InputDDLDir  string
	OutputCsvDir string
	NumRows      int
}

// GendataConfig holds the configuration values
var GendataCfg = GendataConfig{}

// completeGendataConfig validates and completes the gendata configuration
func completeGendataConfig() error {
	// Resolve absolute paths if OutputDir is available from a global config, otherwise use defaults as is
	// Assuming GlobalConfig.OutputDir is accessible similar to dumpCmd. If not, this part needs adjustment.
	// For now, let's assume GlobalConfig.OutputDir exists or use current directory if not.
	// This part might need to be adapted based on how GlobalConfig is structured and initialized.
	// If GlobalConfig.OutputDir is not meant to be used here, we can remove it.
	baseOutputDir := "." // Default to current directory if GlobalConfig.OutputDir is not set/relevant
	if GlobalConfig.OutputDir != "" {
		baseOutputDir = GlobalConfig.OutputDir
	}

	if !filepath.IsAbs(GendataCfg.InputDDLDir) {
		GendataCfg.InputDDLDir = filepath.Join(baseOutputDir, GendataCfg.InputDDLDir)
	}
	if !filepath.IsAbs(GendataCfg.OutputCsvDir) {
		GendataCfg.OutputCsvDir = filepath.Join(baseOutputDir, GendataCfg.OutputCsvDir)
	}

	if GendataCfg.NumRows <= 0 {
		return fmt.Errorf("NumRows must be a positive integer, got %d", GendataCfg.NumRows)
	}

	// Create OutputCsvDir if it doesn't exist
	if err := os.MkdirAll(GendataCfg.OutputCsvDir, 0755); err != nil {
		logrus.Errorf("Failed to create output CSV directory '%s': %v", GendataCfg.OutputCsvDir, err)
		return err
	}
	logrus.Infof("Output CSV directory set to: %s", GendataCfg.OutputCsvDir)
	logrus.Infof("Input DDL directory set to: %s", GendataCfg.InputDDLDir)


	return nil
}

// gendataCmd represents the gendata command
var gendataCmd = &cobra.Command{
	Use:   "gendata",
	Short: "Generates CSV data based on DDL and stats files.",
	Long: `Gendata command reads table structures from DDL (.sql) files and table statistics 
from .stats.yaml files to generate mock CSV data.

Example:
  dorisdump gendata --input-ddl-dir output/ddl --output-csv-dir output/csv_data --num-rows 500`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// It's good practice to initialize global configs here if not already done
		if err := initConfig(cmd); err != nil {
			return err
		}
		return completeGendataConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Info("Gendata command called.")
		logrus.Infof("Configuration:")
		logrus.Infof("  Input DDL Directory: %s", GendataCfg.InputDDLDir)
		logrus.Infof("  Output CSV Directory: %s", GendataCfg.OutputCsvDir)
		logrus.Infof("  Number of Rows per Table: %d", GendataCfg.NumRows)

		// Create src.Config and populate it from GendataCfg
		srcConfig := src.Config{
			InputDDLDir:  GendataCfg.InputDDLDir,
			OutputCsvDir: GendataCfg.OutputCsvDir,
			NumRows:      GendataCfg.NumRows,
		}

		// Call the GenerateData function from the src package
		if err := src.GenerateData(srcConfig); err != nil {
			logrus.Errorf("Error generating data: %v", err)
			return err
		}

		logrus.Info("Data generation completed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gendataCmd)

	gendataCmd.Flags().StringVar(&GendataCfg.InputDDLDir, "input-ddl-dir", "output/ddl", "Directory containing DDL (.sql) and stats (.stats.yaml) files.")
	gendataCmd.Flags().StringVar(&GendataCfg.OutputCsvDir, "output-csv-dir", "output/csv_data", "Directory where CSV files will be generated.")
	gendataCmd.Flags().IntVar(&GendataCfg.NumRows, "num-rows", 1000, "Number of rows to generate per table.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
