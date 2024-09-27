//go:build chimera
// +build chimera

package src

import (
	"sync"

	"github.com/flier/gohs/chimera"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func NewAuditLogScanner(dbs []string, queryMinCpuTimeMs int, queryStates []string, unique, uniqueNormalize, unescape, strict bool) AuditLogScanner {
	return &HyperAuditLogScanner{
		SimpleAuditLogScanner: *NewSimpleAuditLogScanner(dbs, queryMinCpuTimeMs, queryStates, unique, uniqueNormalize, unescape, strict),
	}
}

var _ AuditLogScanner = (*HyperAuditLogScanner)(nil)

type HyperAuditLogScanner struct {
	SimpleAuditLogScanner

	database chimera.BlockDatabase
	scratch  *chimera.Scratch
	close    func()
}

func (s *HyperAuditLogScanner) Init() {
	s.database, s.scratch, s.close = hs_alloc(s.dbs, s.queryStates)
}

func (s *HyperAuditLogScanner) ScanOne(oneLog []byte) error {
	if err := s.database.Scan(oneLog, s.scratch, &hs_handler{}, s); err != nil {
		logrus.Errorln("[hyperscan] Failed to scan audit log file")
		return err
	}

	return nil
}

func (s *HyperAuditLogScanner) Close() {
	s.close()
}

type hs_handler struct{}

// OnMatch will be invoked whenever a match is located in the target data during the execution of a scan.
func (h *hs_handler) OnMatch(_ uint, _, _ uint64, _ uint, captured []*chimera.Capture, ctx any) chimera.Callback {
	c, _ := ctx.(*HyperAuditLogScanner)

	caps := lo.Map(captured, func(cap *chimera.Capture, _ int) string { return string(cap.Bytes) })
	c.onMatch(caps)

	return chimera.Continue
}

// OnError will be invoked when an error event occurs during matching;
// this indicates that some matches for a given expression may not be reported.
func (h *hs_handler) OnError(event chimera.ErrorEvent, _ uint, _, _ any) chimera.Callback {
	logrus.Errorln("[hyperscan] OnError:", event.Error())
	return chimera.Continue
}

func hs_makeAuditLogQueryRegex(dbs, states []string) hyperscanAlloc {
	re := auditlogQueryRe(dbs, states)

	pattern := chimera.NewPattern(re, chimera.MultiLine|chimera.DotAll|chimera.SingleMatch|chimera.Utf8Mode|chimera.UnicodeProperty)
	database, err := chimera.NewBlockDatabase(pattern)
	if err != nil {
		logrus.Fatalf(`[hyperscan] Unable to compile pattern "%s": %v\n`, pattern.String(), err)
	}
	scratchPool := sync.Pool{
		New: func() any {
			scratch, err := chimera.NewManagedScratch(database)
			if err != nil {
				logrus.Fatalln("[hyperscan] Unable to allocate scratch space")
			}
			return scratch
		},
	}

	scratchAlloc := func() (chimera.BlockDatabase, *chimera.Scratch, func()) {
		scratch, _ := scratchPool.Get().(*chimera.Scratch)
		return database, scratch, func() { scratchPool.Put(scratch) }
	}
	return scratchAlloc
}

type hyperscanAlloc = func() (chimera.BlockDatabase, *chimera.Scratch, func())

var (
	hsAlloc hyperscanAlloc
	hslock  sync.Mutex
)

func hs_alloc(dbs, states []string) (chimera.BlockDatabase, *chimera.Scratch, func()) {
	hslock.Lock()
	defer hslock.Unlock()

	if hsAlloc == nil {
		hsAlloc = hs_makeAuditLogQueryRegex(dbs, states)
	}
	return hsAlloc()
}
