//go:build !chimera

package src

func NewAuditLogScanner(dbs []string, queryMinCpuTimeMs int, queryStates []string, unique, uniqueNormalize, unescape, onlySelect, strict bool) AuditLogScanner {
	return NewSimpleAuditLogScanner(dbs, queryMinCpuTimeMs, queryStates, unique, uniqueNormalize, unescape, onlySelect, strict)
}
