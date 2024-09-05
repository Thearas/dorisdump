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
	stmtMatchFmt   = `\|Db=%s\|(?:(?!^\d{4}-\d{2}-\d{2})(?:.|\n))*\|Time(?:\(ms\))?=(\d*)\|(?:(?!^\d{4}-\d{2}-\d{2})(?:.|\n))*\|IsQuery=true\|(?:(?!^\d{4}-\d{2}-\d{2})(?:.|\n))*\|Stmt=((?:(?!^\d{4}-\d{2}-\d{2})(?:.|\n))*)\|CpuTimeMS=`
	stmtMatchStart = []byte("|Stmt=")
	stmtMatchEnd   = []byte("|CpuTimeMS=")
	timeMatchStart = []byte("|Time") // some is '|Time(ms)'

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

func ExtractQueriesFromAuditLogs(dbs []string, auditlogPaths []string, queryMinCpuTimeMs, parallel int) ([]string, error) {
	logrus.Infof("Extracting queries of database %v, audit logs: %v\n", dbs, auditlogPaths)

	g := ParallelGroup(parallel)

	hash2sqls := make([]map[[32]byte]string, len(auditlogPaths))
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

			hash2sql, err := ExtractQueriesFromAuditLog(dbs, []byte(content), queryMinCpuTimeMs)
			if err != nil {
				return err
			}

			// remove ignored queries
			for _, q := range IgnoreQueries {
				delete(hash2sql, q)
			}

			hash2sqls[i] = hash2sql

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return lo.Values(lo.Assign(hash2sqls...)), nil
}
