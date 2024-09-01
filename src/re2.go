//go:build !(linux && amd64)
// +build !linux !amd64

package src

import (
	regexp "github.com/wasilibs/go-re2"
	"github.com/zeebo/blake3"
)

// ExtractQueryOne extracts the query from an audit log.
func ExtractQueriesFromAuditLog(dbs []string, auditlog []byte, queryMinCpuTimeMs int) (map[[32]byte]string, error) {
	regex, err := regexp.Compile(auditlogQueryRe(dbs))
	if err != nil {
		return nil, err
	}

	matches := regex.FindAll(auditlog, -1)

	hash2sql := make(map[[32]byte]string, 1024)
	h := blake3.New()
	for _, match := range matches {
		sql := retrieveStmtFromMatch(match, queryMinCpuTimeMs, true)
		if sql == nil {
			continue
		}
		hash2sql[hash(h, sql)] = string(sql)
	}

	return hash2sql, nil
}
