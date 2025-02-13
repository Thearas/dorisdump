package src

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
)

const (
	ReplaySqlPrefix     = `/*dorisdump{`
	ReplaySqlSuffix     = `*/`
	replayTsFormat      = "2006-01-02 15:04:05.000"
	ReplayResultFileExt = ".result"
)

type ReplayResult struct {
	Ts      string `json:"ts,omitempty"`
	QueryId string `json:"queryId"`

	ReturnRows     int    `json:"returnRows"`
	ReturnRowsHash string `json:"returnRowsHash,omitempty"`
	DurationMs     int64  `json:"durationMs"`
	Err            string `json:"err,omitempty"`
}

func (re *ReplayResult) String() string {
	b, _ := json.Marshal(re)
	return string(b)
}

type ReplayClient struct {
	resultDir string
	dbcfg     *mysql.Config
	cluster   string

	client          string
	sqls            []*ReplaySql
	speed           float32
	maxHashRows     int
	maxConnIdleTime time.Duration
	minTs           int64

	db         *sqlx.DB
	connect    *sqlx.Conn
	resultFile *os.File

	hash *blake3.Hasher
}

func (c *ReplayClient) conn(ctx context.Context, currdb string, reconnect ...bool) (*sqlx.Conn, error) {
	if c.connect == nil || (len(reconnect) > 0 && reconnect[0]) {
		c.Close(false)

		dbcfg := c.dbcfg
		dbcfg.DBName = currdb

		if c.cluster != "" {
			dbcfg = dbcfg.Clone()
			dbcfg.DBName = fmt.Sprintf("%s@%s", dbcfg.DBName, c.cluster)
		}

		var err error
		c.db, err = sqlx.Open("mysql", dbcfg.FormatDSN())
		if err != nil {
			return nil, err
		}
		c.db.SetConnMaxIdleTime(0)
		c.db.SetConnMaxLifetime(0)
		c.db.SetMaxIdleConns(1)
		c.connect, err = c.db.Connx(ctx)
		if err != nil {
			return nil, err
		}
	}

	// switch db
	if currdb != "" && currdb != c.dbcfg.DBName {
		var clusterId string
		if c.cluster != "" {
			clusterId = fmt.Sprintf("@`%s`", c.cluster)
		}
		if _, err := c.connect.ExecContext(ctx, fmt.Sprintf("use `%s`%s", currdb, clusterId)); err != nil {
			logrus.Errorf("client %s switching to db %s failed, err: %v\n", c.client, currdb, err)
			return nil, err
		}
		logrus.Traceln("switching to db", currdb)
		c.dbcfg.DBName = currdb
	}

	return c.connect, nil
}

func (c *ReplayClient) query(ctx context.Context, currdb, stmt string, args ...any) (*sql.Rows, int64, error) {
	conn, err := c.conn(ctx, currdb)
	if err != nil {
		return nil, 0, err
	}

	startedAt := time.Now()
	r, err := conn.QueryContext(ctx, stmt, args...)
	duration := time.Since(startedAt).Milliseconds()

	if err != nil {
		return nil, duration, err
	}
	return r, duration, nil
}

func (c *ReplayClient) queryWithReconnect(ctx context.Context, currdb, stmt string, args ...any) (*sql.Rows, int64, error) {
	r, duration, err := c.query(ctx, currdb, stmt, args...)
	if errors.Is(err, net.ErrClosed) || errors.Is(err, mysql.ErrInvalidConn) {
		// reconnect
		_, err = c.conn(ctx, currdb, true)
		if err != nil {
			return nil, 0, err
		}
		logrus.Debugln("client", c.client, "reconnect")
		r, duration, err = c.query(ctx, currdb, stmt, args...)
	}
	return r, duration, err
}

func (c *ReplayClient) writeResult(b []byte) (err error) {
	if c.resultFile == nil {
		// result file
		resultFilePath := filepath.Join(c.resultDir, fmt.Sprintf("%s%s", c.client, ReplayResultFileExt))
		c.resultFile, err = os.OpenFile(resultFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			logrus.Errorf("open replay result file %s failed, err: %v\n", resultFilePath, err)
			return err
		}
	}
	if _, err := c.resultFile.Write(append(b, '\n')); err != nil {
		logrus.Errorln("client", c.client, "failed to write result:", err)
		return err
	}
	return nil
}

func (c *ReplayClient) Close(closefile bool) {
	if c.db != nil {
		c.db.Close()
		c.db = nil
		c.connect = nil
	}
	if closefile && c.resultFile != nil {
		c.resultFile.Sync()
		c.resultFile.Close()
		c.resultFile = nil
	}
}

func (c *ReplayClient) appendHash(r *sql.Rows) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
	if err != nil {
		return err
	}
	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(sql.RawBytes)
	}

	if err := r.Scan(values...); err != nil {
		return err
	}
	for _, v := range values {
		c.hash.Write(*v.(*sql.RawBytes))
		c.hash.Write([]byte{'\t'}) // append a tab between columns
	}
	c.hash.Write([]byte{'\n'}) // append a newline between rows
	return nil
}

func (c *ReplayClient) consumeHash() string {
	h := fmt.Sprintf("%x", c.hash.Sum(nil))
	c.hash.Reset()
	return h
}

func (c *ReplayClient) replay(ctx context.Context) error {
	logrus.Debugf("replay %d sqls for client %s\n", len(c.sqls), c.client)

	var (
		prevTs         = c.minTs
		prevDurationMs int64
	)

	for _, s := range c.sqls {
		// 1. Wait
		sleepMs := float32(s.Ts-prevTs-prevDurationMs) / c.speed
		if sleepMs > 2 /*ms*/ {
			sleep := time.Duration(sleepMs) * time.Millisecond

			if c.maxConnIdleTime > 0 && sleep > c.maxConnIdleTime {
				// close conn if idle time is too long
				logrus.Traceln("client", c.client, "close idle conn")
				c.Close(false)
			}
			time.Sleep(sleep)
		}
		prevTs = s.Ts
		prevDurationMs = s.DurationMs

		logrus.Traceln("client", c.client, "executing query_id:", s.QueryId, "sql:", s.Stmt)

		// 2. Execute query
		var (
			rowCount  int
			startedAt = time.Now()
		)
		r, durationMs, err := c.queryWithReconnect(ctx, s.Db, s.Stmt)
		if err != nil {
			logrus.Debugf("client %s executed sql failed at query_id: %s, err: %v\n", c.client, s.QueryId, err)
		} else {
			for r.Next() {
				rowCount++
				if rowCount < c.maxHashRows {
					if err = c.appendHash(r); err != nil {
						logrus.Errorf("scan sql return rows failed, query_id: %s, err: %v\n", s.QueryId, err)
						break
					}
				}
			}
			_ = r.Close()
		}

		logrus.Traceln("query_id:", s.QueryId, ", row count:", rowCount, ", duration:", durationMs, "ms")

		result := ReplayResult{
			Ts:         startedAt.Format(replayTsFormat),
			QueryId:    s.QueryId,
			ReturnRows: rowCount,
			DurationMs: durationMs,
		}
		if err != nil {
			result.Err = err.Error()
		}
		if c.maxHashRows > 0 && rowCount > 0 {
			result.ReturnRowsHash = c.consumeHash()
		}

		b, err := json.Marshal(result)
		if err != nil {
			logrus.Errorln("failed to marshal result:", err)
			continue
		}

		if err := c.writeResult(b); err != nil {
			return err
		}
	}

	logrus.Debugf("client %s replay done\n", c.client)

	return nil
}

func ReplaySqls(
	ctx context.Context,
	host string, port uint16, user, password, cluster string,
	resultDir string, client2sqls map[string][]*ReplaySql, speed float32, maxHashRows int, maxConnIdleTime time.Duration,
	minTs int64, parallel int,
) error {
	if len(client2sqls) == 0 {
		return fmt.Errorf("no sqls to replay")
	}
	if parallel < len(client2sqls) {
		logrus.Warnf("Parallel %d is less than client count %d", parallel, len(client2sqls))
		if Confirm("Set parallel to client count") {
			parallel = len(client2sqls)
		}
	}

	logrus.Infof("Replay with %d client, started at %v, speed %f\n",
		len(client2sqls),
		time.UnixMilli(minTs).UTC().Format("2006-01-02 15:04:05"),
		speed,
	)

	dbcfg := &mysql.Config{
		User:                 user,
		Passwd:               password,
		Addr:                 net.JoinHostPort(host, strconv.Itoa(int(port))),
		Net:                  "tcp",
		DBName:               "",
		AllowNativePasswords: true,
		Timeout:              5 * time.Second,
		InterpolateParams:    true,
		ReadTimeout:          600 * time.Second,
		WriteTimeout:         600 * time.Second,
	}

	// test connection
	db, err := sqlx.ConnectContext(ctx, "mysql", dbcfg.FormatDSN())
	if err != nil {
		return err
	}
	db.Close()

	g := ParallelGroup(parallel)
	for client, sqls := range client2sqls {
		client, sqls := client, sqls
		g.Go(func() error {
			cli := ReplayClient{
				resultDir:       resultDir,
				dbcfg:           dbcfg.Clone(),
				cluster:         cluster,
				client:          client,
				sqls:            sqls,
				speed:           speed,
				maxHashRows:     maxHashRows,
				maxConnIdleTime: maxConnIdleTime,
				minTs:           minTs,

				hash: blake3.New(),
			}
			defer cli.Close(true)

			return cli.replay(ctx)
		})
	}

	return g.Wait()
}

func DecodeReplaySqls(
	s *bufio.Scanner,
	dbs, users map[string]struct{},
	from, to int64, // ms
	maxCount int,
) (map[string][]*ReplaySql, int64, error) {
	if !s.Scan() {
		logrus.Warningln("Failed to scan reply sql file, maybe empty?")
		return nil, 0, nil
	}

	var (
		line  = s.Bytes()
		eof   = false
		count int
	)

	// check the replay file is valid by first line prefix
	if !bytes.HasPrefix(line, []byte(ReplaySqlPrefix)) {
		return nil, 0, errors.New("invalid sql replay file")
	}

	client2sqls := make(map[string][]*ReplaySql, 1024)
	minTs := int64(math.MaxInt64)
	for !eof {
		oneSql := bytes.Clone(line)

		// one log may have multiple lines
		// a line not starts with `replaySqlPrefix` is considered belonging to the previous line
		for {
			if !s.Scan() {
				eof = true
				break
			}
			line = s.Bytes()

			if bytes.HasPrefix(line, []byte(ReplaySqlPrefix)) {
				break
			}

			// append to previous line
			oneSql = append(oneSql, '\n')
			oneSql = append(oneSql, line...)
		}

		// parse one sql

		// decode meta
		// include '{' and '}'
		metaStart := len(ReplaySqlPrefix) - 1
		metaEnd := bytes.Index(oneSql, []byte(ReplaySqlSuffix))
		if metaEnd < 0 || oneSql[metaEnd-1] != '}' {
			logrus.Warningln("Failed to extract replay sql meta at:", string(oneSql))
			continue
		}
		meta, err := decodeReplaySqlMeta(oneSql[metaStart:metaEnd])
		if err != nil {
			logrus.Warningln("Failed to parse replay sql meta, err: ", err, ", query:", meta.QueryId)
			continue
		}

		// decode stmt
		stmt := string(bytes.TrimSpace(oneSql[metaEnd+len(ReplaySqlSuffix):]))
		if stmt == "" {
			logrus.Warningln("empty replay sql stmt, query_id:", meta.QueryId)
			continue
		}

		// filter
		if _, ok := dbs[meta.Db]; len(dbs) > 0 && !ok {
			continue
		}
		if _, ok := users[meta.User]; len(users) > 0 && !ok {
			continue
		}
		if !meta.matchTime(from, to) {
			continue
		}
		if maxCount > 0 && count >= maxCount {
			break
		}

		// log may out of order
		ts, err := meta.Timestamp()
		if err != nil {
			continue
		}
		if ts < minTs {
			minTs = meta.Ts
		}

		// add to result
		client2sqls[meta.Client] = append(client2sqls[meta.Client], &ReplaySql{
			ReplaySqlMeta: meta,
			Stmt:          stmt,
		})
		count++
	}

	logrus.Infoln("Found", count, "replay sql(s)")

	return client2sqls, minTs, nil
}

type ReplaySql struct {
	ReplaySqlMeta

	Stmt string
}

func EncodeReplaySql(ts, client, user, db, queryId, stmt string, durationMs int64) string {
	b, err := json.Marshal(ReplaySqlMeta{
		Ts_:        ts,
		Client:     client,
		User:       user,
		Db:         db,
		QueryId:    queryId,
		DurationMs: durationMs,
	})
	if err != nil {
		panic(err)
	}

	outputStmt := fmt.Sprintf(`/*dorisdump%s*/ %s`, b, stmt)
	if !strings.HasSuffix(outputStmt, ";") {
		outputStmt += ";"
	}

	return outputStmt
}

// ReplaySqlMeta will be prepend to every sql as a comment.
//
// e.g.	"/*dorisdump{"ts": "2024-09-20 00:00:00", "client": "127.0.0.1:32345", "user": "root", "db": "test", "queryId": "1"}*/ <the sql>"
type ReplaySqlMeta struct {
	Ts_        string `json:"ts"`
	Ts         int64  `json:"-"`
	Client     string `json:"client"`
	User       string `json:"user"`
	Db         string `json:"db"`
	QueryId    string `json:"queryId"`
	DurationMs int64  `json:"durationMs,omitempty"`
}

func (m *ReplaySqlMeta) matchTime(fromMs, toMs int64) bool {
	if fromMs <= 0 && toMs <= 0 {
		return true
	}

	ts, err := m.Timestamp()
	if err != nil {
		return false
	}

	if fromMs > 0 && ts < fromMs {
		return false
	}
	if toMs > 0 && ts > toMs {
		return false
	}
	return true
}

func (m *ReplaySqlMeta) Timestamp() (ms int64, err error) {
	if m.Ts != 0 {
		return
	}

	ts, err := time.Parse(replayTsFormat, m.Ts_)
	if err != nil {
		logrus.Warningln("wrong timestamp:", m.Ts_, "at query_id:", m.QueryId)
		return 0, err
	}

	m.Ts = ts.UnixMilli()

	return m.Ts, nil
}

func decodeReplaySqlMeta(b []byte) (ReplaySqlMeta, error) {
	meta := ReplaySqlMeta{}
	err := json.Unmarshal(b, &meta)
	return meta, err
}
