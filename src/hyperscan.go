//go:build chimera
// +build chimera

package src

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/flier/gohs/chimera"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
)

func ExtractQueriesFromAuditLog(dbs []string, auditlog *os.File, queryMinCpuTimeMs int, unique bool) (map[[32]byte]string, []string, error) {
	c := hs_newContext(queryMinCpuTimeMs, unique)
	database, scratch, close := hs_alloc(dbs)
	defer close()

	// read log file line by line
	buf := bufio.NewScanner(auditlog)
	if !buf.Scan() {
		logrus.Errorf("Failed to scan audit log file %s, maybe empty?\n", auditlog.Name())
		return c.hash2sql, c.sqls, nil
	}
	var (
		line   = buf.Bytes()
		lineRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
	)

outer:
	for {
		oneLog := line

		// one log may have multiple lines
		// a line not starts with 'yyyy-mm-dd' is considered belonging to the previous line
		for {
			if !buf.Scan() {
				break outer
			}
			line = buf.Bytes()

			const minLenToMatch = len("yyyy-mm-dd")
			if len(line) >= minLenToMatch && lineRe.Match(line[:minLenToMatch+1]) {
				break
			}

			// append to previous line
			oneLog = append(oneLog, '\n')
			oneLog = append(oneLog, line...)
		}

		// parse log
		if err := database.Scan(oneLog, scratch, &hs_handler{}, &c); err != nil {
			logrus.Errorln("[hyperscan] Failed to scan audit log file")
			return nil, nil, err
		}
		oneLog = nil
	}

	return c.hash2sql, c.sqls, nil
}

type hs_handler struct{}

// OnMatch will be invoked whenever a match is located in the target data during the execution of a scan.
func (h *hs_handler) OnMatch(_ uint, _, _ uint64, flags uint, captured []*chimera.Capture, ctx any) chimera.Callback {
	c, _ := ctx.(*hyperscanContext)

	time, client, durationMs, stmt := captured[1].Bytes[:], captured[2].Bytes[:], captured[3].Bytes[:], captured[4].Bytes[:]

	ok := filterStmtFromMatch(c.queryMinCpuTimeMs, durationMs, stmt)
	if !ok {
		return chimera.Continue
	}

	stmt = ShortenTabSpaces(stmt)

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
	// add leading comment in JSON with time and client
	time_ := string(time)
	client_ := anonymizeHashStr(c.hash, string(client))
	stmt_ := strings.TrimSpace(string(stmt))
	if !strings.HasSuffix(stmt_, ";") {
		stmt_ = stmt_ + ";"
	}
	stmt_ = fmt.Sprintf(`/*{"time": "%s", "client": "%s"}*/ %s`, time_, client_, stmt_)

	c.sqls = append(c.sqls, stmt_)
	return chimera.Continue
}

// OnError will be invoked when an error event occurs during matching;
// this indicates that some matches for a given expression may not be reported.
func (h *hs_handler) OnError(event chimera.ErrorEvent, _ uint, _, _ any) chimera.Callback {
	logrus.Errorln("[hyperscan] OnError:", event.Error())
	return chimera.Continue
}

func hs_makeAuditLogQueryRegex(dbs []string) hyperscanAlloc {
	re := auditlogQueryRe(dbs)

	pattern := chimera.NewPattern(re, chimera.MultiLine|chimera.DotAll|chimera.SingleMatch)
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
	queryMinCpuTimeMs int
	unique            bool

	// when unique
	hash     *blake3.Hasher
	hash2sql map[[32]byte]string
	// when not unique
	sqls []string
}

func hs_newContext(queryMinCpuTimeMs int, unique bool) hyperscanContext {
	return hyperscanContext{
		queryMinCpuTimeMs: queryMinCpuTimeMs,
		unique:            unique,
		hash:              blake3.New(),
		hash2sql:          make(map[[32]byte]string, 1024),
		sqls:              make([]string, 0, 1024),
	}
}
