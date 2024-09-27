//go:build !chimera

package src

func NewAuditLogScanner(dbs []string, queryMinCpuTimeMs int, queryStates []string, unique, uniqueNormalize, unescape bool) AuditLogScanner {
	return NewSimpleAuditLogScanner(dbs, queryMinCpuTimeMs, queryStates, unique, uniqueNormalize, unescape)
}
