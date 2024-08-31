package src

import (
	"os"
	"sync"

	"github.com/flier/gohs/hyperscan"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
)

// ExtractQueryOne extracts the query from an audit log.
func ExtractQueriesFromAuditLog(dbs []string, auditlog []byte) (map[[32]byte]string, error) {
	database, scratch, close := hs_alloc(dbs)
	defer close()

	c := hs_newContext(auditlog)
	if err := database.Scan(auditlog, scratch, hs_Callback, c); err != nil {
		logrus.Errorln("[hyperscan] Failed to scan audit log file")
		return nil, err
	}

	return c.hash2sql, nil
}

func hs_makeAuditLogQueryRegex(dbs []string) hyperscanAlloc {
	re := auditlogQueryRe(dbs)

	pattern := hyperscan.NewPattern(re, hyperscan.SomLeftMost)
	database, err := hyperscan.NewBlockDatabase(pattern)
	if err != nil {
		logrus.Fatalf(`[hyperscan] Unable to compile pattern "%s": %v\n`, pattern.String(), err)
	}
	scratchPool := sync.Pool{
		New: func() interface{} {
			scratch, err := hyperscan.NewManagedScratch(database)
			if err != nil {
				logrus.Fatalln(os.Stderr, "[hyperscan] Unable to allocate scratch space")
			}
			return scratch
		},
	}

	scratchAlloc := func() (hyperscan.BlockDatabase, *hyperscan.Scratch, func()) {
		scratch, _ := scratchPool.Get().(*hyperscan.Scratch)
		return database, scratch, func() { scratchPool.Put(scratch) }
	}
	return scratchAlloc
}

type hyperscanAlloc = func() (hyperscan.BlockDatabase, *hyperscan.Scratch, func())

var (
	hsAlloc hyperscanAlloc
	hslock  sync.Mutex
)

func hs_alloc(dbs []string) (hyperscan.BlockDatabase, *hyperscan.Scratch, func()) {
	hslock.Lock()
	defer hslock.Unlock()

	if hsAlloc == nil {
		hsAlloc = hs_makeAuditLogQueryRegex(dbs)
	}
	return hsAlloc()
}

type hyperscanContext struct {
	content  []byte
	hash     *blake3.Hasher
	hash2sql map[[32]byte]string
}

func hs_newContext(content []byte) hyperscanContext {
	return hyperscanContext{
		content:  content,
		hash:     blake3.New(),
		hash2sql: make(map[[32]byte]string, 1024),
	}
}

func hs_Callback(id uint, from, to uint64, _ uint, ctx any) error {
	c, _ := ctx.(hyperscanContext)

	match := c.content[from:to]
	sql := retrieveStmtFromMatch(match, true)
	if sql == nil {
		return nil
	}

	hash := hash(c.hash, sql)
	if _, ok := c.hash2sql[hash]; ok {
		return nil
	}

	c.hash2sql[hash] = string(sql)

	return nil
}
