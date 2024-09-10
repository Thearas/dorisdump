//go:build chimera
// +build chimera

package src

import (
	"fmt"
	"strings"
	"sync"

	"github.com/flier/gohs/chimera"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
)

// ExtractQueryOne extracts the query from an audit log.
func ExtractQueriesFromAuditLog(dbs []string, auditlog []byte, queryMinCpuTimeMs int, unique bool) (map[[32]byte]string, []string, error) {
	database, scratch, close := hs_alloc(dbs)
	defer close()

	c := hs_newContext(auditlog, queryMinCpuTimeMs, unique)
	if err := database.Scan(auditlog, scratch, &hs_handler{}, &c); err != nil {
		logrus.Errorln("[hyperscan] Failed to scan audit log file")
		return nil, nil, err
	}

	return c.hash2sql, c.sqls, nil
}

type hs_handler struct{}

// OnMatch will be invoked whenever a match is located in the target data during the execution of a scan.
func (h *hs_handler) OnMatch(_ uint, _, _ uint64, flags uint, captured []*chimera.Capture, ctx any) chimera.Callback {
	c, _ := ctx.(*hyperscanContext)

	time, client, durationMs, stmt := captured[1].Bytes, captured[2].Bytes, captured[3].Bytes, captured[4].Bytes

	ok := filterStmtFromMatch(c.queryMinCpuTimeMs, durationMs, stmt)
	if !ok {
		return chimera.Continue
	}

	// unique sqls
	if c.unique {
		// not unique sqls
		hash := hash(c.hash, stmt)
		if _, ok := c.hash2sql[hash]; !ok {
			c.hash2sql[hash] = string(stmt)
		}
		return chimera.Continue
	}

	// not unique sqls
	// replace all newline into space
	stmt_ := strings.ReplaceAll(string(stmt), "\n", " ")

	// add leading comment in JSON with time and client
	time_ := string(time)
	client_ := anonymizeHashStr(c.hash, string(client))
	stmt_ = fmt.Sprintf(`/*{"time": "%s", "client": "%s"}*/ %s`, time_, client_, stmt_)

	c.sqls = append(c.sqls, stmt_)
	return chimera.Continue
}

// OnError will be invoked when an error event occurs during matching;
// this indicates that some matches for a given expression may not be reported.
func (h *hs_handler) OnError(_ chimera.ErrorEvent, _ uint, _, _ any) chimera.Callback {
	return chimera.Continue
}

func hs_makeAuditLogQueryRegex(dbs []string) hyperscanAlloc {
	re := auditlogQueryRe(dbs)

	pattern := chimera.NewPattern(re, chimera.MultiLine)
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

func hs_alloc(dbs []string) (chimera.BlockDatabase, *chimera.Scratch, func()) {
	hslock.Lock()
	defer hslock.Unlock()

	if hsAlloc == nil {
		hsAlloc = hs_makeAuditLogQueryRegex(dbs)
	}
	return hsAlloc()
}

type hyperscanContext struct {
	content           []byte
	queryMinCpuTimeMs int
	unique            bool

	// when unique
	hash     *blake3.Hasher
	hash2sql map[[32]byte]string
	// when not unique
	sqls []string
}

func hs_newContext(content []byte, queryMinCpuTimeMs int, unique bool) hyperscanContext {
	return hyperscanContext{
		content:           content,
		queryMinCpuTimeMs: queryMinCpuTimeMs,
		unique:            unique,
		hash:              blake3.New(),
		hash2sql:          make(map[[32]byte]string, 1024),
		sqls:              make([]string, 0, 1024),
	}
}
