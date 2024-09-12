package src

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/edsrzf/mmap-go"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

var (
	// stmtMatchFmt is the regex pattern to extract the query statement from the audit log.
	//
	// NOTE: A bit hacky, but it works for now.
	//
	// Tested on v2.0.12 and v2.1.x. Not sure if it also works on others Doris version.
	stmtMatchFmt = `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d*) \[query\] \|Client=(.+)\|.*\|Db=%s\|.*\|Time(?:\(ms\))?=(\d*)\|.*\|IsQuery=true\|.*\|Stmt=(.*)\|CpuTimeMS=`

	IgnoreQueries = lo.Map([]string{
		`SELECT CONCAT("'", user, "'@'",host,"'") FROM mysql.user`,
		`SELECT @@max_allowed_packet`,
		`SELECT DATABASE()`,
		`SELECT name from mysql.help_topic WHERE name like "SHOW %"`,
		`select @@version_comment limit 1`,
		`select connection_id()`,
	}, func(s string, _ int) [32]byte { return hash(hasher, []byte(s)) })
)

func auditlogQueryRe(dbs []string) string {
	allowDBs := lo.Map(dbs, func(s string, _ int) string { return regexp.QuoteMeta(s) })

	var dbFilter string
	if len(dbs) > 0 {
		dbFilter = "(?:" + strings.Join(allowDBs, "|") + ")"
	} else {
		dbFilter = ".*"
	}

	return fmt.Sprintf(stmtMatchFmt, dbFilter)
}

func filterStmtFromMatch(queryMinDurationMs int, stmtTimeMs, stmt []byte) bool {
	if queryMinDurationMs > 0 {
		if len(stmtTimeMs) == 0 {
			return false
		}
		ms, err := strconv.Atoi(string(stmtTimeMs))
		if err != nil || ms < queryMinDurationMs {
			return false
		}
	}

	// remove dorisdump self queries
	if bytes.HasPrefix(stmt, InternalSqlCommentBytes) {
		return false
	}

	return true
}

// ExtractQueriesFromAuditLog extracts the query from an audit log.
// In unique mode: we will aggregate all queries in the first slot, the remaining slots will be empty.
// In non-unique mode: we will return all queries in the order of auditlogPath slots.
func ExtractQueriesFromAuditLogs(dbs []string, auditlogPaths []string, queryMinCpuTimeMs, parallel int, unique bool) ([][]string, error) {
	logrus.Infof("Extracting queries of database %v, audit logs: %v\n", dbs, auditlogPaths)

	g := ParallelGroup(parallel)

	hash2sqls := make([]map[[32]byte]string, len(auditlogPaths))
	sqlss := make([][]string, len(auditlogPaths))
	for i, auditlogPath := range auditlogPaths {
		i, auditlogPath := i, auditlogPath
		g.Go(func() error {
			f, err := os.Open(auditlogPath)
			if err != nil {
				logrus.Errorln("Unable to open audit log file:", auditlogPath)
				return err
			}
			defer f.Close()

			content, err := mmap.Map(f, mmap.RDONLY, 0)
			if err != nil {
				logrus.Errorln("Unable to mmap audit log file:", auditlogPath)
				return err
			}
			defer content.Unmap()

			logrus.Debugln("Extracting queries from audit log:", auditlogPath)

			hash2sql, sqls, err := ExtractQueriesFromAuditLog(dbs, []byte(content), queryMinCpuTimeMs, unique)
			if err != nil {
				return err
			}

			// remove ignored queries
			// FIXME: do we need to remove ignored queries when not in unique output mode?
			if len(hash2sql) > 0 {
				for _, q := range IgnoreQueries {
					delete(hash2sql, q)
				}
			}

			hash2sqls[i] = hash2sql
			sqlss[i] = sqls

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if unique {
		return [][]string{lo.Values(lo.Assign(hash2sqls...))}, nil
	}
	return sqlss, nil
}
