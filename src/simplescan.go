//go:build !chimera

package src

func NewAuditLogScanner(opts AuditLogScanOpts) AuditLogScanner {
	return NewSimpleAuditLogScanner(opts)
}
