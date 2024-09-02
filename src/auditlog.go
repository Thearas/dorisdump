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
	// Tested on v2.1.x. Not sure if it also works on others Doris version.
	stmtMatchFmt        = `\|Db=%s\|.*\|IsQuery=true\|.*\|Stmt=(.*)\|CpuTimeMS=\d*`
	stmtMatchStart      = []byte("|Stmt=")
	stmtMatchEnd        = []byte("|CpuTimeMS=")
	cpuTimeMsMatchStart = stmtMatchEnd
)

func auditlogQueryRe(dbs []string) string {
	allowDBs := lo.Map(dbs, func(s string, _ int) string { return regexp.QuoteMeta(s) })
	dbFilter := "(?:" + strings.Join(allowDBs, "|") + ")"
	return fmt.Sprintf(stmtMatchFmt, dbFilter)
}

func retrieveStmtFromMatch(match []byte, minCpuTimeMs int, filterDorisDumpSelfSql bool) []byte {
	// 1. Retrieve query
	s := bytes.Index(match, stmtMatchStart)
	e := bytes.LastIndex(match, stmtMatchEnd)
	stmt := match[s+len(stmtMatchStart) : e]

	// remove dorisdump self queries
	if filterDorisDumpSelfSql && bytes.HasPrefix(stmt, InternalSqlCommentBytes) {
		return nil
	}

	// 2. Retrieve cpu time
	if minCpuTimeMs == 0 {
		return stmt
	}
	cpuTime_ := string(match[e+len(cpuTimeMsMatchStart):])
	cpuTimeMs, err := strconv.Atoi(cpuTime_)
	if err != nil || cpuTimeMs < minCpuTimeMs {
		return nil
	}

	return stmt
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
