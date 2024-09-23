//go:build !chimera

package src

func NewAuditLogScanner(dbs []string, queryMinCpuTimeMs int, unique, uniqueNormalize, unescape bool) AuditLogScanner {
	return NewSimpleAuditLogScanner(dbs, queryMinCpuTimeMs, unique, uniqueNormalize, unescape)
}
