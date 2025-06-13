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
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Thearas/dorisdump/src"
)

const (
	DefaultGenDataDDLs = "output/ddl"
	DefaultGenDataDir  = "output/gendata"
)

// GendataConfig holds the configuration values
var GendataConfig = Gendata{}

// Gendata holds the configuration for the gendata command
type Gendata struct {
	DDL           string
	OutputDataDir string
	GenConf       string
	NumRows       int

	genFromDDLs []string
}

// gendataCmd represents the gendata command
var gendataCmd = &cobra.Command{
	Use:   "gendata",
	Short: "Generates CSV data based on DDL and stats files.",
	Long: `Gendata command reads table structures from DDL (.table.sql) files and table statistics files (.stats.yaml) to generate fake CSV data.

Example:
  dorisdump gendata --dbs db1,db2
  dorisdump gendata --dbs db1 --tables t1,t2 --rows 500 --ddl output/ddl/
  dorisdump gendata --ddl create.table.sql`,
	Aliases: []string{"g"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(cmd)
	},
	SilenceUsage:               true,
	SuggestionsMinimumDistance: 6,
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := cmd.Context()

		if err := completeGendataConfig(); err != nil {
			return err
		}
		if err := src.SetupDefaultGenRules(GendataConfig.GenConf); err != nil {
			return err
		}
		GlobalConfig.Parallel = lo.Min([]int{GlobalConfig.Parallel, len(GendataConfig.genFromDDLs)})

		logrus.Infof("Generate data for %d tables, parallel: %d\n", len(GendataConfig.genFromDDLs), GlobalConfig.Parallel)

		g := src.ParallelGroup(GlobalConfig.Parallel)
		for _, ddlFile := range GendataConfig.genFromDDLs {
			logrus.Debugf("generating data to %s ...\n", strings.TrimSuffix(ddlFile, ".table.sql"))

			ddl, err := src.ReadFileOrStdin(ddlFile)
			if err != nil {
				return err
			}
			stats, err := findTableStats(ddlFile)
			if err != nil {
				return err
			}
			tg, err := src.NewTableGen(ddl, stats)
			if err != nil {
				return err
			}

			if GlobalConfig.DryRun {
				continue
			}

			o, err := createOutputGenDataWriter(ddlFile)
			if err != nil {
				return err
			}

			g.Go(func() error {
				defer o.Close()

				w := bufio.NewWriterSize(o, 256*1024)
				if err := tg.GenCSV(w, GendataConfig.NumRows); err != nil {
					return err
				}

				return w.Flush()
			})
		}

		return g.Wait()
	},
}

func init() {
	rootCmd.AddCommand(gendataCmd)
	gendataCmd.PersistentFlags().SortFlags = false
	gendataCmd.Flags().SortFlags = false

	pFlags := gendataCmd.PersistentFlags()
	pFlags.StringVarP(&GendataConfig.DDL, "ddl", "d", DefaultGenDataDDLs, "Directory or file containing DDL (.table.sql) and stats (.stats.yaml) files")
	pFlags.StringVarP(&GendataConfig.OutputDataDir, "output-data-dir", "o", DefaultGenDataDir, "Directory where CSV files will be generated")
	pFlags.IntVarP(&GendataConfig.NumRows, "rows", "r", 0, fmt.Sprintf("Number of rows to generate per table (default %d)", src.DefaultGenRowCount))
	pFlags.StringVarP(&GendataConfig.GenConf, "genconf", "c", "", "Generator config file")
}

// completeGendataConfig validates and completes the gendata configuration
func completeGendataConfig() (err error) {
	isDefaultDDLs := GendataConfig.DDL == DefaultGenDataDDLs
	if isDefaultDDLs {
		GendataConfig.DDL = filepath.Join(GlobalConfig.OutputDir, "ddl")
	}
	isDefaultGenDataOutputDir := GendataConfig.OutputDataDir == DefaultGenDataDir
	if isDefaultGenDataOutputDir {
		GendataConfig.OutputDataDir = filepath.Join(GlobalConfig.OutputDir, "gendata")
	}

	if err := src.CheckGenRowCount(GendataConfig.NumRows); err != nil {
		return err
	}

	// if --ddl is a sql file, not need --dbs or --tables
	if stat, err := os.Stat(GendataConfig.DDL); err == nil && !stat.IsDir() {
		if !strings.HasSuffix(stat.Name(), ".sql") {
			return errors.New("the ddl file must ends with '.sql'")
		}
		GendataConfig.genFromDDLs = []string{GendataConfig.DDL}
		return nil
	}

	GlobalConfig.DBs, GlobalConfig.Tables = lo.Uniq(GlobalConfig.DBs), lo.Uniq(GlobalConfig.Tables)
	dbs, tables := GlobalConfig.DBs, GlobalConfig.Tables
	if len(dbs) == 0 && len(tables) == 0 {
		return errors.New("expected at least one database or tables, please use --dbs/--tables flag or --ddl flag with a '.sql' file")
	} else if len(dbs) == 1 {
		// prepend default database if only one database specified
		prefix := dbs[0] + "."
		for i, t := range GlobalConfig.Tables {
			if !strings.Contains(t, ".") {
				GlobalConfig.Tables[i] = prefix + t
			}
		}
	} else {
		for _, t := range tables {
			if !strings.Contains(t, ".") {
				return errors.New("expected database in table name when zero/multiple databases specified, e.g. --tables db1.table1,db2.table2")
			}
		}
	}

	ddls := []string{}
	if len(GlobalConfig.Tables) == 0 {
		for _, db := range GlobalConfig.DBs {
			fmatch := filepath.Join(GendataConfig.DDL, fmt.Sprintf("%s.*.table.sql", db))
			tableddls, err := src.FileGlob([]string{fmatch})
			if err != nil {
				logrus.Errorf("Get db '%s' ddls in '%s' failed\n", db, fmatch)
				return err
			}
			ddls = append(ddls, tableddls...)
		}
	} else {
		for _, table := range GlobalConfig.Tables {
			tableddl := filepath.Join(GendataConfig.DDL, fmt.Sprintf("%s.table.sql", table))
			ddls = append(ddls, tableddl)
		}
	}
	GendataConfig.genFromDDLs = ddls

	return nil
}

func findTableStats(ddlFileName string) (*src.TableStats, error) {
	ddlFileDir := filepath.Dir(ddlFileName)
	ddlFileName = filepath.Base(ddlFileName)

	// table ddl file has 4 parts: {db}.{table}.table.sql
	// db stats file has 2 parts: {db}.stats.yaml
	isDumpTable := len(strings.Split(ddlFileName, ".")) == 4 && strings.HasSuffix(ddlFileName, ".table.sql")
	if !isDumpTable {
		return nil, nil
	}

	db := strings.SplitN(ddlFileName, ".", 2)[0]
	dbStatsFile := filepath.Join(ddlFileDir, db+".stats.yaml")
	b, err := os.ReadFile(dbStatsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logrus.Debugf("stats file '%s' not found for db '%s'\n", dbStatsFile, db)
			return nil, nil
		}
		return nil, err
	}

	dbstats := &src.DBSchema{}
	if err := yaml.Unmarshal(b, dbstats); err != nil {
		return nil, err
	}

	table := strings.SplitN(ddlFileName, ".", 3)[1]
	for _, tableStats := range dbstats.Stats {
		if tableStats.Name != table || len(tableStats.Columns) == 0 || tableStats.RowCount <= 0 {
			continue
		}
		if tableStats.Columns[0].Method != "FULL" {
			logrus.Warnf("Table stats '%s.%s' is '%s' in '%s', better to run 'ANALYZE DATABASE `%s` WITH SYNC' and dump again\n", db, table, tableStats.Columns[0].Method, dbStatsFile, db)
		}
		return tableStats, nil
	}

	logrus.Warnf("Table stats '%s.%s' not found in '%s', better to run 'ANALYZE DATABASE `%s` WITH SYNC' and dump again\n", db, table, dbStatsFile, db)
	return nil, nil
}

func createOutputGenDataWriter(ddlFileName string) (*os.File, error) {
	ddlFileName = filepath.Base(ddlFileName)
	dir := filepath.Join(GendataConfig.OutputDataDir, strings.TrimSuffix(strings.TrimSuffix(ddlFileName, ".table.sql"), ".sql"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	file := filepath.Join(dir, "1.csv")
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logrus.Fatalln("Can not open output data file:", file, ", err:", err)
	}
	return f, nil
}
