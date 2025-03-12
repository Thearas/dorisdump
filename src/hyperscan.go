//go:build chimera
// +build chimera

package src

import (
	"sync"

	"github.com/flier/gohs/chimera"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func NewAuditLogScanner(opts AuditLogScanOpts) AuditLogScanner {
	return &HyperAuditLogScanner{
		SimpleAuditLogScanner: *NewSimpleAuditLogScanner(opts),
	}
}

var _ AuditLogScanner = (*HyperAuditLogScanner)(nil)

type HyperAuditLogScanner struct {
	SimpleAuditLogScanner

	database chimera.BlockDatabase
	scratch  *chimera.Scratch
}

func (s *HyperAuditLogScanner) Init() {
	s.database, s.scratch = hs_alloc(s.DBs, s.QueryStates, s.OnlySelect)
}

func (s *HyperAuditLogScanner) ScanOne(oneLog []byte) error {
	if err := s.database.Scan(oneLog, s.scratch, &hs_handler{}, s); err != nil {
		logrus.Errorln("[hyperscan] Failed to scan audit log file")
		return err
	}

	return nil
}

func (s *HyperAuditLogScanner) Close() {
	s.scratch.Free()
}

type hs_handler struct{}

// OnMatch will be invoked whenever a match is located in the target data during the execution of a scan.
func (h *hs_handler) OnMatch(_ uint, _, _ uint64, _ uint, captured []*chimera.Capture, ctx any) chimera.Callback {
	c, _ := ctx.(*HyperAuditLogScanner)

	caps := lo.Map(captured[1:], func(cap *chimera.Capture, _ int) string { return string(cap.Bytes) })
	c.onMatch(caps, false)

	return chimera.Continue
}

// OnError will be invoked when an error event occurs during matching;
// this indicates that some matches for a given expression may not be reported.
func (h *hs_handler) OnError(event chimera.ErrorEvent, _ uint, _, _ any) chimera.Callback {
	logrus.Errorln("[hyperscan] OnError:", event.Error())
	return chimera.Continue
}

type hyperscanAlloc = func() (chimera.BlockDatabase, *chimera.Scratch)

var (
	hsAlloc hyperscanAlloc
	hslock  sync.Mutex
)

func hs_alloc(dbs, states []string, onlySelect bool) (chimera.BlockDatabase, *chimera.Scratch) {
	hslock.Lock()
	defer hslock.Unlock()

	if hsAlloc == nil {
		hsAlloc = hs_makeAuditLogQueryRegex(dbs, states, onlySelect)
	}
	return hsAlloc()
}

func hs_makeAuditLogQueryRegex(dbs, states []string, onlySelect bool) hyperscanAlloc {
	re := auditlogQueryRe(dbs, states, onlySelect)

	pattern := chimera.NewPattern(re, chimera.MultiLine|chimera.DotAll|chimera.SingleMatch|chimera.Utf8Mode|chimera.UnicodeProperty)
	database, err := chimera.NewManagedBlockDatabase(pattern)
	if err != nil {
		logrus.Fatalf(`[hyperscan] Unable to compile pattern "%s": %v\n`, pattern.String(), err)
	}
	scratch, err := chimera.NewScratch(database)
	if err != nil {
		logrus.Fatalln("[hyperscan] Unable to allocate scratch space")
	}

	scratchAlloc := func() (chimera.BlockDatabase, *chimera.Scratch) {
		scratchClone, err := scratch.Clone()
		if err != nil {
			logrus.Fatalln("[hyperscan] Unable to clone scratch")
		}
		return database, scratchClone
	}
	return scratchAlloc
}
