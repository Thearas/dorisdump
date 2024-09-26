package src

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/pingcap/tidb/pkg/parser"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

var (
	// stmtMatchFmt is the regex pattern to extract the query statement from the audit log.
	//
	// NOTE: A bit hacky, but it works for now.
	//
	// Tested on v2.0.x and v2.1.x. Not sure if it also works on others Doris version.
	stmtMatchFmt = `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d*) \[query\] \|Client=([^|]+)\|User=([^|]+)\|.*\|Db=(%s)\|.*\|Time(?:\(ms\))?=(\d*)\|.*\|QueryId=([a-z0-9-]+)\|IsQuery=true\|.*\|Stmt=(.*)\|CpuTimeMS=`

	unescapeReplacer = strings.NewReplacer(
		"\\n", "\n",
		"\\t", "\t",
		"\\r", "\r",
	)
)

type AuditLogScanner interface {
	Init()
	ScanOne(oneLine []byte) error
	Result() (sqls []string, err error)
	Close()
}

var _ AuditLogScanner = (*SimpleAuditLogScanner)(nil)

// ExtractQueriesFromAuditLog extracts the query from an audit log.
func ExtractQueriesFromAuditLogs(
	dbs []string,
	auditlogPaths []string,
	encoding string,
	queryMinCpuTimeMs int,
	parallel int,
	unique, uniqueNormalize bool,
	unescape bool,
) ([][]string, error) {
	logrus.Infof("Extracting queries of database %v, audit logs: %v\n", dbs, auditlogPaths)

	g := ParallelGroup(parallel)

	sqlss := make([][]string, len(auditlogPaths))
	for i, auditlogPath := range auditlogPaths {
		i, auditlogPath := i, auditlogPath
		g.Go(func() error {
			f, err := os.Open(auditlogPath)
			if err != nil {
				logrus.Errorln("Unable to open audit log file:", auditlogPath)
				return err
			}
			defer f.Close()

			// detect encoding
			enc, err := detectAuditLogEncoding(encoding, f)
			if err != nil {
				return err
			}
			buf := bufio.NewScanner(transform.NewReader(f, enc.NewDecoder()))

			logrus.Debugln("Extracting queries from audit log:", auditlogPath, "with encoding:", enc)

			// read log file line by line
			s := NewAuditLogScanner(dbs, queryMinCpuTimeMs, unique, uniqueNormalize, unescape)
			sqls, err := extractQueriesFromAuditLog(s, buf)
			if err != nil {
				return err
			}

			sqlss[i] = sqls

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return sqlss, nil
}

func extractQueriesFromAuditLog(
	s AuditLogScanner,
	auditlog *bufio.Scanner,
) ([]string, error) {
	s.Init()
	defer s.Close()

	// read log file line by line
	if !auditlog.Scan() {
		logrus.Warningln("Failed to scan audit log file, maybe empty?")
		return []string{}, nil
	}
	var (
		line   = auditlog.Bytes()
		lineRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
		eof    = false
	)

	for !eof {
		oneLog := bytes.Clone(line)

		// one log may have multiple lines
		// a line not starts with 'yyyy-mm-dd' is considered belonging to the previous line
		for {
			if !auditlog.Scan() {
				eof = true
				break
			}
			line = auditlog.Bytes()

			const minLenToMatch = len("yyyy-mm-dd")
			if len(line) >= minLenToMatch && lineRe.Match(line[:minLenToMatch+1]) {
				break
			}

			// append to previous line
			oneLog = append(oneLog, '\n')
			oneLog = append(oneLog, line...)
		}

		// parse log
		if err := s.ScanOne(oneLog); err != nil {
			logrus.Errorln("Failed to scan audit log file")
			return nil, err
		}
	}

	return s.Result()
}

type SimpleAuditLogScanner struct {
	hash *blake3.Hasher

	dbs                     []string
	encoding                encoding.Encoding
	queryMinCpuTimeMs       int
	unique, uniqueNormalize bool
	unescape                bool

	hash2sql map[[32]byte]string
	sqls     []string

	re *regexp2.Regexp
}

func NewSimpleAuditLogScanner(dbs []string, queryMinCpuTimeMs int, unique, uniqueNormalize, unescape bool) *SimpleAuditLogScanner {
	return &SimpleAuditLogScanner{
		hash:              blake3.New(),
		dbs:               dbs,
		queryMinCpuTimeMs: queryMinCpuTimeMs,
		unique:            unique,
		uniqueNormalize:   uniqueNormalize,
		unescape:          unescape,
		hash2sql:          make(map[[32]byte]string, 1024),
		sqls:              make([]string, 0, 1024),
	}
}

func (s *SimpleAuditLogScanner) Init() {
	s.re = regexp2.MustCompile(auditlogQueryRe(s.dbs), regexp2.Multiline|regexp2.Singleline|regexp2.Unicode|regexp2.Compiled)
}

func (s *SimpleAuditLogScanner) ScanOne(oneLog []byte) error {
	matches, err := s.re.FindStringMatch(string(oneLog))
	if err != nil {
		logrus.Errorln("Failed to scan audit log file")
		return err
	}
	if matches == nil || len(matches.Groups()) < 2 {
		return nil
	}

	caps := lo.Map(matches.Groups(), func(g regexp2.Group, _ int) string { return g.String() })
	s.onMatch(caps)

	return nil
}

func (s *SimpleAuditLogScanner) Result() (sqls []string, err error) {
	return s.sqls, nil
}

func (s *SimpleAuditLogScanner) Close() {}

func (s *SimpleAuditLogScanner) onMatch(caps []string) {
	time, client, user, db, durationMs, queryId, stmt := caps[1], caps[2], caps[3], caps[4], caps[5], caps[6], caps[7]

	stmt = strings.TrimSpace(stmt)
	ok := s.filterStmtFromMatch(s.queryMinCpuTimeMs, durationMs, queryId, stmt)
	if !ok {
		return
	}

	if s.unescape {
		stmt = unescapeReplacer.Replace(stmt)
	}

	// add leading meta comment
	outputStmt := EncodeReplaySql(time, client, user, db, queryId, stmt)

	// unique sqls
	if s.unique {
		var h [32]byte
		if s.uniqueNormalize {
			h = hashstr(s.hash, parser.NormalizeKeepHint(stmt))
		} else {
			h = hashstr(s.hash, stmt)
		}
		if _, ok := s.hash2sql[h]; ok {
			return
		}
		s.hash2sql[h] = outputStmt
	}

	// not unique sqls
	s.sqls = append(s.sqls, outputStmt)
}

func (s *SimpleAuditLogScanner) filterStmtFromMatch(queryMinDurationMs int, durationMs, queryId, stmt string) bool {
	if queryMinDurationMs > 0 {
		if len(durationMs) == 0 {
			return false
		}
		ms, err := strconv.Atoi(durationMs)
		if err != nil || ms < queryMinDurationMs {
			return false
		}
	}

	// remove empty stmt
	if len(stmt) == 0 {
		return false
	}

	// remove doris self queries

	// remove dorisdump self queries
	if strings.HasPrefix(stmt, InternalSqlComment) {
		return false
	}

	// remove truncated queries (which length is larger than audit_plugin_max_sql_length)
	truncated := false
	if strings.HasSuffix(stmt, "...") {
		truncated = true
	} else if strings.HasSuffix(stmt, "*/") && strings.LastIndex(stmt, "... /*") > -1 {
		truncated = true
	}
	if truncated {
		logrus.Warningln("query has been truncated, query_id:", queryId)
		return false
	}

	return true
}

func auditlogQueryRe(dbs []string) string {
	allowDBs := lo.Map(dbs, func(s string, _ int) string { return regexp.QuoteMeta(s) })

	var dbFilter string
	if len(dbs) > 0 {
		dbFilter = strings.Join(allowDBs, "|")
	} else {
		dbFilter = "[^|]*"
	}

	return fmt.Sprintf(stmtMatchFmt, dbFilter)
}

func detectAuditLogEncoding(encoding string, f *os.File) (encoding.Encoding, error) {
	// detect encoding
	var err error
	if encoding == "auto" {
		// detect encoding
		encoding, err = DetectCharset(bufio.NewReader(f))
		if err != nil {
			return nil, fmt.Errorf("cannot detect charset of audit log %s: %v", f.Name(), err)
		}
		if _, err = f.Seek(0, 0); err != nil {
			return nil, fmt.Errorf("cannot seek audit log %s: %v", f.Name(), err)
		}
	}
	enc, err := GetEncoding(encoding)
	if err != nil {
		return nil, err
	}

	return enc, nil
}
