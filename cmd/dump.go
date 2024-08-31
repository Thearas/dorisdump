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
	"slices"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/Thearas/dorisdump/src"
)

var DumpConfig = Dump{}

type Dump struct {
	AuditLogPaths         []string
	OutputDDLDir          string
	OutputQueryDir        string
	LocalAuditLogCacheDir string

	SSHAddress    string
	SSHPassword   string
	SSHPrivateKey string

	AnonymizerEnabled    bool
	AnonymizerMethod     string
	StripComment         bool
	AnonymizerReserveIds []string

	DumpSchema bool
	DumpQuery  bool

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
	Aliases:          []string{"d"},
	Example:          "dorisdump dump --audit-log /path/to/audit.log -D db1,db2 -T table1,db2.table2",
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initConfig(cmd)
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		completeDumpConfig()

		if DumpConfig.Clean {
			if err := cleanCmd.RunE(nil, nil); err != nil {
				return err
			}
		}

		if DumpConfig.AnonymizerEnabled {
			src.SetupAnonymizer(DumpConfig.AnonymizerReserveIds...)
		}

		// dump schemas
		if DumpConfig.DumpSchema {
			schemas, err := dumpSchemas(ctx)
			if err != nil {
				return err
			}

			logrus.Infof("Found %d schema(s)\n", len(schemas))

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

			logrus.Infof("Found %d query(s)\n", len(queries))

			if err := outputQueries(queries); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	pFlags := dumpCmd.PersistentFlags()
	pFlags.BoolVar(&DumpConfig.AnonymizerEnabled, "anonymize", false, "Anonymize sqls")
	pFlags.StringSliceVar(&DumpConfig.AnonymizerReserveIds, "anonymize-reserve-ids", nil, "Skip anonymization for these ids, usually database names")
	pFlags.StringVar(&DumpConfig.AnonymizerMethod, "anonymize-method", "hash", "Anonymize method, hash only for now")
	pFlags.MarkHidden("anonymize-method")
	pFlags.BoolVar(&DumpConfig.StripComment, "strip-comment", false, "Strip comments")
	pFlags.BoolVar(&DumpConfig.DumpSchema, "dump-schema", true, "Dump schema")
	pFlags.BoolVar(&DumpConfig.DumpQuery, "dump-query", false, "Dump query from audit log")
	pFlags.StringSliceVar(&DumpConfig.AuditLogPaths, "audit-logs", nil, "Audit log paths, either local path or ssh://xxx, default is ssh://root@{db_host}:22/{fe_dir}/log/fe.audit.log")
	pFlags.StringVar(&DumpConfig.SSHAddress, "ssh-address", "", "SSH address for downloading audit log, default is root@{db_host}:22")
	pFlags.StringVar(&DumpConfig.SSHPassword, "ssh-password", "", "SSH password for --ssh-address")
	pFlags.StringVar(&DumpConfig.SSHPrivateKey, "ssh-private-key", "~/.ssh/id_rsa", "File path of SSH private key for --ssh-address")

	flags := dumpCmd.Flags()
	flags.BoolVar(&DumpConfig.Clean, "clean", false, "Clean previous data and output directory")
}

func completeDumpConfig() error {
	DumpConfig.OutputDDLDir = filepath.Join(GlobalConfig.OutputDir, "ddl")
	DumpConfig.OutputQueryDir = filepath.Join(GlobalConfig.OutputDir, "sql")
	DumpConfig.LocalAuditLogCacheDir = filepath.Join(GlobalConfig.DataDir, "auditlog")

	GlobalConfig.DBs, GlobalConfig.Tables = lo.Uniq(GlobalConfig.DBs), lo.Uniq(GlobalConfig.Tables)
	dbs, tables := GlobalConfig.DBs, GlobalConfig.Tables
	if len(dbs) == 0 {
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

	DumpConfig.SSHPrivateKey = src.ExpandHome(DumpConfig.SSHPrivateKey)
	if DumpConfig.SSHAddress == "" {
		DumpConfig.SSHAddress = fmt.Sprintf("ssh://root@%s:22", GlobalConfig.DBHost)
	}
	if !strings.HasPrefix(DumpConfig.SSHAddress, "ssh://") {
		DumpConfig.SSHAddress = "ssh://" + DumpConfig.SSHAddress
	}

	return nil
}

func dumpSchemas(ctx context.Context) ([]*src.Schema, error) {
	dbs, tables := GlobalConfig.DBs, GlobalConfig.Tables
	g := errgroup.Group{}
	g.SetLimit(10)

	schemas := make([][]*src.Schema, len(dbs))
	for i, db := range dbs {
		db := db
		g.Go(func() error {
			logrus.Infof("Dumping schemas from %s...\n", db)
			conn, err := connectDB(db)
			if err != nil {
				return err
			}
			defer conn.Close()

			schemas[i], err = src.ShowCreateTables(ctx, conn, db, tables...)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return slices.Concat(schemas...), nil
}

func outputSchemas(schemas []*src.Schema) error {
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

	g := errgroup.Group{}
	g.SetLimit(10)
	for _, s := range schemas {
		s := s
		g.Go(func() error {
			if DumpConfig.AnonymizerEnabled {
				s.Name = src.Anonymize(DumpConfig.AnonymizerMethod, s.Name)
				s.CreateStmt = src.AnonymizeSql(DumpConfig.AnonymizerMethod, s.CreateStmt)
			}

			if printSql {
				logrus.Tracef("schema: %+v\n", *s)
			}

			filename := fmt.Sprintf("%s.%s.%s.sql", s.DB, s.Name, s.Type.Lower())
			path := filepath.Join(DumpConfig.OutputDDLDir, filename)
			if GlobalConfig.DryRun {
				return nil
			}
			return src.WriteFile(path, s.CreateStmt)
		})
	}

	return g.Wait()
}

func dumpQueries(ctx context.Context) ([]string, error) {
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

	queries, err := src.ExtractQueriesFromAuditLogs(GlobalConfig.DBs, auditLogFiles)
	if err != nil {
		logrus.Errorf("Extract queries from audit logs failed, %v\n", err)
		return nil, err
	}

	return queries, nil
}

func outputQueries(queries []string) error {
	if len(queries) == 0 {
		return nil
	}

	if !GlobalConfig.DryRun {
		if err := os.MkdirAll(DumpConfig.OutputQueryDir, 0755); err != nil {
			logrus.Errorln("Create output query directory failed, ", err)
			return err
		}
	}

	printSql := logrus.GetLevel() == logrus.TraceLevel

	count := 0
	size := len(queries)
	for size != 0 {
		size /= 10
		count++
	}
	format := fmt.Sprintf("q%%0%dd.sql", count)

	g := errgroup.Group{}
	g.SetLimit(10)
	for i, query := range queries {
		i, query := i, query
		g.Go(func() error {
			if DumpConfig.AnonymizerEnabled {
				query = src.AnonymizeSql(DumpConfig.AnonymizerMethod, query)
			}

			name := fmt.Sprintf(format, i)

			if printSql {
				logrus.Tracef("queries %s: %+v\n", name, query)
			}

			path := filepath.Join(DumpConfig.OutputQueryDir, name)
			if GlobalConfig.DryRun {
				return nil
			}
			return src.WriteFile(path, query)
		})
	}

	return g.Wait()
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
