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
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Thearas/dorisdump/src"
)

var DumpConfig = Dump{}

type Dump struct {
	AuditLogPaths         []string
	AuditLogUnescape      bool
	OutputDDLDir          string
	OutputQueryDir        string
	LocalAuditLogCacheDir string
	AuditLogEncoding      string

	SSHAddress    string
	SSHPassword   string
	SSHPrivateKey string

	DumpSchema           bool
	DumpStats            bool
	DumpQuery            bool
	QueryOutputMode      string
	QueryUniqueNormalize bool
	QueryMinDuration_    time.Duration
	QueryMinDurationMs   int
	QueryStates          []string
	OnlySelect           bool
	Strict               bool

	Clean bool
}

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump schema and query for Doris",
	Long: `
Dump schema from DB and query from audit-log.

You may want to pass config by '$HOME/.dorisdump.yaml',
or environment variables with prefix 'DORIS_', e.g.
    DORIS_HOST=xxx
    DORIS_PORT=9030
	`,
	Aliases:                    []string{"d"},
	Example:                    "dorisdump dump --dump-schema --dump-query -dbs db1 --audit-logs /path/to/audit.log",
	TraverseChildren:           true,
	SuggestionsMinimumDistance: 2,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initConfig(cmd)
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		if err := completeDumpConfig(); err != nil {
			return err
		}

		if DumpConfig.Clean {
			if err := cleanAllFiles(true); err != nil {
				return err
			}
		}

		if AnonymizeConfig.Enabled {
			src.SetupAnonymizer(AnonymizeConfig.IdMinLength, AnonymizeConfig.ReserveIds...)
		}

		// dump schemas
		if DumpConfig.DumpSchema {
			schemas, err := dumpSchemas(ctx)
			if err != nil {
				return err
			}

			logrus.Infof("Found %d schema(s)\n", lo.SumBy(schemas, func(s *src.DBSchema) int { return len(s.Schemas) }))

			if err := outputSchemas(schemas); err != nil {
				return err
			}
		}

		// dump queries
		if DumpConfig.DumpQuery {
			queries, err := dumpQueries(ctx)
			if err != nil {
				return err
			}

			logrus.Infof("Found %d query(s)\n", lo.SumBy(queries, func(s []string) int { return len(s) }))

			if err := outputQueries(queries); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.PersistentFlags().SortFlags = false
	dumpCmd.Flags().SortFlags = false

	pFlags := dumpCmd.PersistentFlags()
	pFlags.BoolVar(&DumpConfig.DumpSchema, "dump-schema", false, "Dump schema")
	pFlags.BoolVar(&DumpConfig.DumpStats, "dump-stats", true, "Dump schema stats, only take effect when '--dump-schema=true'")
	pFlags.BoolVar(&DumpConfig.DumpQuery, "dump-query", false, "Dump query from audit log")
	pFlags.StringVar(&DumpConfig.QueryOutputMode, "query-output-mode", "default", "Dump query output mode, one of [default, unique]")
	pFlags.BoolVar(&DumpConfig.QueryUniqueNormalize, "query-unique-normalize", false, "Regard 'select 1 from b where a = 1' as 'select ? from b where a = ?' for unique, only take effect when '--query-output-mode=unique'")
	pFlags.DurationVar(&DumpConfig.QueryMinDuration_, "query-min-duration", 0, "Dump queries which execution duration is greater than or equal to")
	pFlags.StringSliceVar(&DumpConfig.QueryStates, "query-states", []string{}, "Dump queries with states, like 'ok', 'eof' and 'err'")
	pFlags.BoolVar(&DumpConfig.OnlySelect, "only-select", true, "Only dump SELECT queries")
	pFlags.BoolVarP(&DumpConfig.Strict, "strict", "s", false, "Filter out sqls that can't be parsed")
	pFlags.StringSliceVar(&DumpConfig.AuditLogPaths, "audit-logs", nil, "Audit log paths, either local path or ssh://xxx")
	pFlags.BoolVar(&DumpConfig.AuditLogUnescape, "audit-log-unescape", true, "Unescape '\\n', '\\t' and '\\r' in audit log")
	pFlags.StringVar(&DumpConfig.AuditLogEncoding, "audit-log-encoding", "auto", "Audit log encoding, like utf8, gbk, ...")
	pFlags.StringVar(&DumpConfig.SSHAddress, "ssh-address", "", "SSH address for downloading audit log, default is 'root@{db_host}:22'")
	pFlags.StringVar(&DumpConfig.SSHPassword, "ssh-password", "", "SSH password for '--ssh-address'")
	pFlags.StringVar(&DumpConfig.SSHPrivateKey, "ssh-private-key", "~/.ssh/id_rsa", "File path of SSH private key for '--ssh-address'")
	addAnonymizeBaseFlags(pFlags, false)

	flags := dumpCmd.Flags()
	flags.BoolVar(&DumpConfig.Clean, "clean", false, "Clean previous data and output directory")
}

func completeDumpConfig() error {
	if !DumpConfig.DumpSchema && !DumpConfig.DumpQuery {
		return errors.New("Expected at least one of --dump-schema or --dump-query")
	}

	DumpConfig.OutputDDLDir = filepath.Join(GlobalConfig.OutputDir, "ddl")
	DumpConfig.OutputQueryDir = filepath.Join(GlobalConfig.OutputDir, "sql")
	DumpConfig.LocalAuditLogCacheDir = filepath.Join(GlobalConfig.DataDir, "auditlog")

	if DumpConfig.QueryMinDuration_ > 0 {
		DumpConfig.QueryMinDurationMs = int(DumpConfig.QueryMinDuration_.Milliseconds())
	}

	GlobalConfig.DBs, GlobalConfig.Tables = lo.Uniq(GlobalConfig.DBs), lo.Uniq(GlobalConfig.Tables)
	dbs, tables := GlobalConfig.DBs, GlobalConfig.Tables
	if DumpConfig.DumpSchema && len(dbs) == 0 {
		return errors.New("Expected at least one database, please use --dbs flag")
	} else if len(dbs) == 1 {
		// prepend default database if only one database specified
		prefix := dbs[0] + "."
		for i, t := range tables {
			if !strings.Contains(t, ".") {
				tables[i] = prefix + t
			}
		}
	} else {
		for _, t := range tables {
			if !strings.Contains(t, ".") {
				return errors.New("Expected database in table name when multiple databases specified, e.g. --tables db1.table1,db2.table2")
			}
		}
	}

	DumpConfig.QueryStates = lo.Map(DumpConfig.QueryStates, func(s string, _ int) string {
		return strings.ToUpper(s)
	})

	DumpConfig.SSHPrivateKey = src.ExpandHome(DumpConfig.SSHPrivateKey)
	if DumpConfig.SSHAddress == "" {
		DumpConfig.SSHAddress = fmt.Sprintf("ssh://root@%s:22", GlobalConfig.DBHost)
	}
	if !strings.HasPrefix(DumpConfig.SSHAddress, "ssh://") {
		DumpConfig.SSHAddress = "ssh://" + DumpConfig.SSHAddress
	}

	return nil
}

func dumpSchemas(ctx context.Context) ([]*src.DBSchema, error) {
	dbs, tables := GlobalConfig.DBs, GlobalConfig.Tables
	g := src.ParallelGroup(GlobalConfig.Parallel)

	schemas := make([]*src.DBSchema, len(dbs))
	for i, db := range dbs {
		i, db := i, db
		g.Go(func() error {
			logrus.Infof("Dumping schemas from %s...\n", db)
			conn, err := connectDB(db)
			if err != nil {
				return err
			}
			defer conn.Close()

			// dump schema
			createTables, err := src.ShowCreateTables(ctx, conn, db, tables...)
			if err != nil {
				return err
			}

			// dump stats
			if !DumpConfig.DumpStats {
				return nil
			}
			tbls := lo.Map(createTables, func(s *src.Schema, _ int) string { return s.Name })
			stats, err := src.GetTablesStats(ctx, conn, db, tbls...)
			if err != nil {
				return err
			}

			schemas[i] = &src.DBSchema{
				Name:    db,
				Schemas: createTables,
				Stats:   stats,
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return schemas, nil
}

func outputSchemas(schemas []*src.DBSchema) error {
	if len(schemas) == 0 {
		return nil
	}

	if !GlobalConfig.DryRun {
		if err := os.MkdirAll(DumpConfig.OutputDDLDir, 0755); err != nil {
			logrus.Errorln("Create output ddl directory failed, ", err)
			return err
		}
	}

	printSql := logrus.GetLevel() == logrus.TraceLevel

	g := src.ParallelGroup(GlobalConfig.Parallel)
	for _, s := range schemas {
		s := s
		g.Go(func() error {
			// 1. write each schema into split file
			for _, s := range s.Schemas {
				var filename string
				if AnonymizeConfig.Enabled {
					s.DB = src.Anonymize(AnonymizeConfig.Method, s.DB)
					s.Name = src.Anonymize(AnonymizeConfig.Method, s.Name)
				}

				filename = fmt.Sprintf("%s.%s.%s.sql", s.DB, s.Name, s.Type.Lower())
				if AnonymizeConfig.Enabled {
					s.CreateStmt = src.AnonymizeSql(AnonymizeConfig.Method, filename, s.CreateStmt)
				}

				if printSql {
					logrus.Tracef("schema: %+v\n", *s)
				}

				path := filepath.Join(DumpConfig.OutputDDLDir, filename)
				if GlobalConfig.DryRun {
					return nil
				}
				if err := src.WriteFile(path, s.CreateStmt); err != nil {
					return err
				}

			}

			// 2. write all stats into one file
			if len(s.Stats) == 0 {
				return nil
			}
			if AnonymizeConfig.Enabled {
				s.Name = src.Anonymize(AnonymizeConfig.Method, s.Name)
				for _, s := range s.Stats {
					s.Name = src.Anonymize(AnonymizeConfig.Method, s.Name)
					for _, c := range s.Columns {
						c.Name = src.Anonymize(AnonymizeConfig.Method, c.Name)
					}
				}
			}
			yml_, err := yaml.Marshal(s)
			if err != nil {
				return err
			}
			yml := string(yml_)

			if printSql {
				logrus.Tracef("stats: \n%s\n", yml)
			}
			if GlobalConfig.DryRun {
				return nil
			}

			filename := fmt.Sprintf("%s.stats.yaml", s.Name)
			path := filepath.Join(DumpConfig.OutputDDLDir, filename)
			return src.WriteFile(path, yml)
		})
	}

	return g.Wait()
}

func dumpQueries(ctx context.Context) ([][]string, error) {
	auditLogs := DumpConfig.AuditLogPaths
	if len(auditLogs) == 0 {
		sshUrl, err := chooseRemoteAuditLog(ctx)
		if err != nil {
			return nil, fmt.Errorf("Please specific audit log path by --audit-logs, error: %v", err)
		}
		auditLogs = []string{sshUrl}
	}

	logrus.Debugf("audit log paths: %+v", auditLogs)

	auditLogFiles := []string{}
	for _, auditLog := range auditLogs {
		var localPath string

		// 1. Remote audit log. scp remote path to local path
		if strings.HasPrefix(auditLog, "ssh://") {
			localPath = filepath.Join(DumpConfig.LocalAuditLogCacheDir, path.Base(auditLog))
			if err := copyAuditLog(ctx, auditLog, localPath); err != nil {
				logrus.Errorln("Copy remote audit log failed:", err)
				return nil, err
			}
			auditLogFiles = append(auditLogFiles, localPath)
			continue
		}

		// 2. Local audit log.
		localPath = strings.TrimPrefix(auditLog, "file://")
		localPaths, err := filepath.Glob(localPath)
		if err != nil {
			return nil, fmt.Errorf("Invalid audit log path: %s, error: %v", localPath, err)
		}

		auditLogFiles = append(auditLogFiles, localPaths...)
	}

	logrus.Infoln("Dumping queries from audit logs...")

	queries, err := src.ExtractQueriesFromAuditLogs(
		GlobalConfig.DBs,
		auditLogFiles,
		DumpConfig.AuditLogEncoding,
		DumpConfig.QueryMinDurationMs,
		DumpConfig.QueryStates,
		GlobalConfig.Parallel,
		DumpConfig.QueryOutputMode == "unique",
		DumpConfig.QueryUniqueNormalize,
		DumpConfig.AuditLogUnescape,
		DumpConfig.OnlySelect,
		DumpConfig.Strict,
	)
	if err != nil {
		logrus.Errorf("Extract queries from audit logs failed, %v\n", err)
		return nil, err
	}

	return queries, nil
}

func outputQueries(queries [][]string) error {
	if !GlobalConfig.DryRun {
		if err := os.MkdirAll(DumpConfig.OutputQueryDir, 0755); err != nil {
			logrus.Errorln("Create output query directory failed, ", err)
			return err
		}
	}

	printSql := logrus.GetLevel() == logrus.TraceLevel

	return outputDefaultQueries(queries, printSql)
}

func outputDefaultQueries(queriess [][]string, printSql bool) error {
	if len(queriess) == 0 {
		return nil
	}

	format := outputQueryFileNameFormat(len(queriess))

	g := src.ParallelGroup(GlobalConfig.Parallel)
	for i, queries := range queriess {
		i, queries := i, queries
		g.Go(func() (err error) {
			name := fmt.Sprintf(format, i)
			path := filepath.Join(DumpConfig.OutputQueryDir, name)

			if AnonymizeConfig.Enabled {
				for i, query := range queries {
					queries[i] = src.AnonymizeSql(AnonymizeConfig.Method, name+"#"+strconv.Itoa(i), query)
				}
			}

			var f *os.File
			if !GlobalConfig.DryRun {
				f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
				if err != nil {
					return err
				}
				defer f.Close()
			}
			for i, query := range queries {
				if printSql {
					logrus.Tracef("queries %s: %+v\n", name+"#"+strconv.Itoa(i), query)
				}
				if GlobalConfig.DryRun {
					continue
				}
				_, err = f.WriteString(query + "\n")
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	return g.Wait()
}

func outputQueryFileNameFormat(total int) string {
	count := 0
	for total != 0 {
		total /= 10
		count++
	}

	return fmt.Sprintf("q%%0%dd.sql", count)
}

func chooseRemoteAuditLog(ctx context.Context) (string, error) {
	conn, err := connectDB(GlobalConfig.DBs[0])
	if err != nil {
		return "", err
	}
	defer conn.Close()

	dir, err := src.ShowFronendsDisksDir(ctx, conn, "audit-log")
	if err != nil {
		return "", err
	}

	sshUrl, err := expandSSHPath(fmt.Sprintf("%s%s", DumpConfig.SSHAddress, dir))
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(sshUrl, "/") {
		sshUrl += "/"
	}
	sshUrl += "fe.audit.log*"

	auditLogs, err := src.SshLs(ctx, DumpConfig.SSHPrivateKey, sshUrl)
	if err != nil {
		return "", fmt.Errorf("SSH list remote audit log failed: %v", err)
	}
	if len(auditLogs) == 0 {
		return "", errors.New("No audit log found on remote server")
	}

	choosed, err := src.Choose("Choose audit log on remote server to dump:", auditLogs)
	if err != nil {
		return "", err
	}

	return expandSSHPath(fmt.Sprintf("%s%s", DumpConfig.SSHAddress, choosed))
}

func copyAuditLog(ctx context.Context, remotePath, localPath string) error {
	remotePath, err := expandSSHPath(remotePath)
	if err != nil {
		return err
	}

	err = src.ScpFromRemote(ctx, DumpConfig.SSHPrivateKey, remotePath, localPath)
	if err != nil {
		err = fmt.Errorf("scp failed, please check --ssh-password or --ssh-private-key: %v", err)
	}
	return err
}

func expandSSHPath(remotePath string) (string, error) {
	// default
	u, err := url.Parse(DumpConfig.SSHAddress)
	if err != nil {
		return "", err
	}
	defaultUser := u.User.Username()
	defaultHost := u.Host
	defaultPass, passInAddr := u.User.Password()
	if !passInAddr {
		defaultPass = DumpConfig.SSHPassword
	}

	// remotePath custom
	u, err = u.Parse(remotePath)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		u.Host = defaultHost
	}
	user := u.User.Username()
	pass, passOk := u.User.Password()
	if user == "" {
		user = defaultUser
	}
	if !passOk {
		pass = defaultPass
	}
	u.User = url.UserPassword(user, pass)

	return u.String(), nil
}

func connectDB(db string) (*sqlx.DB, error) {
	if db == "" {
		return nil, fmt.Errorf("database name is required, please use --db flag")
	}
	return src.NewDB(GlobalConfig.DBHost, GlobalConfig.DBPort, GlobalConfig.DBUser, GlobalConfig.DBPassword, db)
}
