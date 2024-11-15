package src

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/edsrzf/mmap-go"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/Thearas/dorisdump/src/parser"
)

var (
	// stmtMatchFmt is the regex pattern to extract the query statement from the audit log.
	//
	// NOTE: A bit hacky, but it works for now.
	//
	// Tested on v2.0.x and v2.1.x. Not sure if it also works on others Doris version.
	stmtMatchFmt = `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d*) \[[^\]]+?\] \|Client=([^|]+?)\|User=([^|]+?)\|(?:.+?)\|Db=(%s?)\|State=%s\|(?:.+?)\|Time(?:\(ms\))?=(\d*)\|(?:.+?)\|QueryId=([a-z0-9-]+)\|IsQuery=%s\|(?:.+?)\|Stmt=(.+?)\|CpuTimeMS=`

	unescapeReplacer = strings.NewReplacer(
		"\\n", "\n",
		"\\t", "\t",
		"\\r", "\r",
	)

	// filterStmtRe filters out some statements from the audit log.
	filterStmtRe = regexp.MustCompile("(?i)^(EXPLAIN|SHOW|USE)")
)

type AuditLogScanner interface {
	Init()
	ScanOne(oneLine []byte) error
	Result() (sqls []string)
	Close()
}

var _ AuditLogScanner = (*SimpleAuditLogScanner)(nil)

type AuditLogScanOpts struct {
	// filter
	DBs                []string
	QueryMinDurationMs int
	QueryStates        []string
	OnlySelect         bool
	From, To           string

	Unescape bool
	Strict   bool
}

func (opts *AuditLogScanOpts) sqlConditions() string {
	// push down filter to db
	conditions := "1=1"
	if len(opts.DBs) > 0 {
		conditions += fmt.Sprintf(` AND db IN ("%s")`, strings.Join(opts.DBs, `", "`))
	}
	if opts.QueryMinDurationMs > 0 {
		conditions += fmt.Sprintf(` AND query_time >= %d`, opts.QueryMinDurationMs)
	}
	if len(opts.QueryStates) > 0 {
		conditions += fmt.Sprintf(` AND state IN ("%s")`, strings.Join(opts.QueryStates, `", "`))
	}
	if opts.OnlySelect {
		conditions += ` AND is_query = 1`
	}

	if opts.From != "" {
		conditions += fmt.Sprintf(` AND time >= "%s"`, opts.From)
	}
	if opts.To != "" {
		conditions += fmt.Sprintf(` AND time <= "%s"`, opts.To)
	}
	return conditions
}

// ExtractQueriesFromAuditLog extracts the query from an audit log.
func ExtractQueriesFromAuditLogs(
	auditlogPaths []string,
	encoding string,
	opts AuditLogScanOpts,
	parallel int,
) ([][]string, error) {
	logrus.Infof("Extracting queries of database %v, audit logs: %v\n", opts.DBs, auditlogPaths)

	g := ParallelGroup(parallel)
	useMmap := os.Getenv("MMAP") != ""

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

			var reader io.Reader = f
			if useMmap {
				m, err := mmap.Map(f, mmap.RDONLY, 0)
				if err != nil {
					logrus.Errorln("Unable to mmap audit log file:", auditlogPath, "err:", err)
					return err
				}
				defer m.Unmap()
				reader = bytes.NewReader(m)
			}

			buf := bufio.NewScanner(transform.NewReader(reader, enc.NewDecoder()))
			buf.Buffer(make([]byte, 10*1024*1024), 1024*1024)

			logrus.Debugln("Extracting queries from audit log:", auditlogPath, "with encoding:", enc)

			// read log file line by line
			s := NewAuditLogScanner(opts)
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

	return s.Result(), nil
}

type SimpleAuditLogScanner struct {
	AuditLogScanOpts

	hash *blake3.Hasher

	uniqsqls map[[32]byte]*uniqSql
	sqls     []string

	re *regexp2.Regexp
}

type uniqSql struct {
	count, sqlIdx int
}

func NewSimpleAuditLogScanner(opts AuditLogScanOpts) *SimpleAuditLogScanner {
	return &SimpleAuditLogScanner{
		AuditLogScanOpts: opts,
		hash:             blake3.New(),
		uniqsqls:         make(map[[32]byte]*uniqSql, 1024),
		sqls:             make([]string, 0, 1024),
	}
}

func (s *SimpleAuditLogScanner) Init() {
	s.re = regexp2.MustCompile(auditlogQueryRe(s.DBs, s.QueryStates, s.OnlySelect), regexp2.Multiline|regexp2.Singleline|regexp2.Unicode|regexp2.Compiled)
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
	s.onMatch(caps[1:], false)

	return nil
}

func (s *SimpleAuditLogScanner) Result() []string {
	return s.sqls
}

func (s *SimpleAuditLogScanner) Close() {}

func (s *SimpleAuditLogScanner) validateSQL(queryId, stmt string) error {
	p := parser.NewParser(queryId, stmt)
	_, err := p.Parse()
	return err
}

func (s *SimpleAuditLogScanner) onMatch(caps []string, skipOptsFilter bool) {
	time, client, user, db, durationMs, queryId, stmt := caps[0], caps[1], caps[2], caps[3], caps[4], caps[5], caps[6]

	stmt = strings.TrimSpace(stmt)
	ok := s.filterStmtFromMatch(time, durationMs, queryId, stmt, skipOptsFilter)
	if !ok {
		return
	}

	if s.Unescape {
		stmt = unescapeReplacer.Replace(stmt)
	}
	if s.Strict && s.validateSQL(queryId, stmt) != nil {
		return
	}

	// add leading meta comment
	outputStmt := EncodeReplaySql(time, client, user, db, queryId, stmt)

	s.sqls = append(s.sqls, outputStmt)
}

func (s *SimpleAuditLogScanner) filterStmtFromMatch(
	time, durationMs, queryId, stmt string,
	skipOptsFilter bool,
) bool {
	// remove empty stmt
	if len(stmt) == 0 {
		return false
	}

	// remove truncated queries (which length is larger than audit_plugin_max_sql_length)
	if logStmtTruncated(queryId, stmt) {
		return false
	}

	// remove dorisdump self queries
	if strings.HasPrefix(stmt, InternalSqlComment) {
		return false
	}

	// remove explain, show and use statements
	if !s.OnlySelect && filterStmtRe.MatchString(stmt) {
		return false
	}

	if skipOptsFilter {
		return true
	}

	// filter by opts below

	if s.From != "" && strings.SplitN(time, ",", 2)[0] < s.From {
		return false
	}
	if s.To != "" && strings.SplitN(time, ",", 2)[0] > s.To {
		return false
	}

	if s.QueryMinDurationMs > 0 {
		if len(durationMs) == 0 {
			return false
		}
		ms, err := strconv.Atoi(durationMs)
		if err != nil || ms < s.QueryMinDurationMs {
			return false
		}
	}

	return true
}

// Which length is larger than audit_plugin_max_sql_length.
func logStmtTruncated(queryId, stmt string) bool {
	var truncated bool
	if strings.HasSuffix(stmt, "...") {
		truncated = true
	} else if strings.HasSuffix(stmt, "*/") && strings.LastIndex(stmt, "... /*") > -1 {
		truncated = true
	}
	if truncated {
		logrus.Warningln("query has been truncated, query_id:", queryId)
	}
	return truncated
}

func auditlogQueryRe(dbs, states []string, onlySelect bool) string {
	var dbFilter string
	if len(dbs) > 0 {
		allowDBs := lo.Map(dbs, func(s string, _ int) string { return regexp.QuoteMeta(s) })
		dbFilter = strings.Join(allowDBs, "|")
	} else {
		dbFilter = "[^|]*"
	}

	var stateFilter string
	if len(states) > 0 {
		allowStates := lo.Map(states, func(s string, _ int) string { return regexp.QuoteMeta(s) })
		stateFilter = strings.Join(allowStates, "|")
	} else {
		stateFilter = "[^|]*"
	}

	isQuery := "(?:true|false)"
	if onlySelect {
		isQuery = "true"
	}

	return fmt.Sprintf(stmtMatchFmt, dbFilter, stateFilter, isQuery)
}

func detectAuditLogEncoding(encoding string, f *os.File) (encoding.Encoding, error) {
	// detect encoding
	var err error
	if encoding == "auto" || encoding == "" {
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
