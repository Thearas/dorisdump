package src

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var (
	InternalSqlComment      = "/*dorisdump*/"
	InternalSqlCommentBytes = []byte(InternalSqlComment)

	sqlLikeReplacer = strings.NewReplacer(
		`"`, `\"`,
		`_`, `\_`,
		`%`, `\%`,
	)
)

type SchemaType string

var (
	SchemaTypeTable            SchemaType = "TABLE"
	SchemaTypeView             SchemaType = "VIEW"
	SchemaTypeMaterializedView SchemaType = "MATERIALIZED_VIEW"
)

func (s SchemaType) Lower() string {
	return strings.ToLower(string(s))
}

func (s SchemaType) sanitize() SchemaType {
	switch s {
	case "BASE TABLE":
		return SchemaTypeTable
	case "VIEW":
		return SchemaTypeView
	default:
		logrus.Warnf("unknown schema type: %s", s)
		return SchemaType(strings.ReplaceAll(string(s), " ", "_"))
	}
}

type DBSchema struct {
	Name    string        `yaml:"db"`
	Schemas []*Schema     `yaml:"-"`
	Stats   []*TableStats `yaml:"tables,omitempty"`
}

type Schema struct {
	Name       string     `db:"TABLE_NAME"`
	Type       SchemaType `db:"TABLE_TYPE"`
	DB         string     `db:"TABLE_SCHEMA"`
	CreateStmt string     `db:"-"`
}

func (s *Schema) String() string {
	return fmt.Sprintf("%s.%s", s.DB, s.Name)
}

type TableStats struct {
	Name     string         `yaml:"name"`
	RowCount int64          `yaml:"row_count"`
	Columns  []*ColumnStats `yaml:"columns,omitempty"`
}

type ColumnStats struct {
	Name        string `yaml:"name"`
	Count       int64  `yaml:"-"`
	Ndv         int64  `yaml:"ndv"`
	NullCount   int64  `yaml:"null_count"`
	DataSize    int64  `yaml:"data_size"`
	AvgSizeByte int64  `yaml:"avg_size_byte"`
	Min         string `yaml:"min"`
	Max         string `yaml:"max"`
}

func NewDB(host string, port int16, user, password, db string) (*sqlx.DB, error) {
	cfg := &mysql.Config{
		User:                 user,
		Passwd:               password,
		Addr:                 net.JoinHostPort(host, strconv.Itoa(int(port))),
		Net:                  "tcp",
		DBName:               db,
		AllowNativePasswords: true,
		Timeout:              3 * time.Second,
		InterpolateParams:    true, // some doris does not enable prepare stmt
		ParseTime:            false,
		ReadTimeout:          600 * time.Second,
		WriteTimeout:         600 * time.Second,
	}
	dsn := cfg.FormatDSN()
	logrus.Traceln("Connecting:", logrus.Fields{
		"Host": host,
		"Port": port,
		"User": user,
		"DB":   db,
	})
	return sqlx.Connect("mysql", dsn)
}

func ShowCreateTables(ctx context.Context, conn *sqlx.DB, db string, dbTables ...string) (schemas []*Schema, err error) {
	schemas_, err := ShowTables(ctx, conn, db)
	if err != nil {
		return nil, err
	}
	tables_ := lo.Map(schemas_, func(s *Schema, _ int) string { return s.Name })
	logrus.Debugln("found tables:", tables_)

	schemas = schemas_

	// filter tables
	if len(dbTables) > 0 {
		dbTables = lo.Filter(dbTables, func(t string, _ int) bool { return strings.HasPrefix(t, db+".") })

		schemas = make([]*Schema, 0, len(dbTables))
		for _, t := range dbTables {
			schema, find := lo.Find(schemas_, func(s *Schema) bool { return s.String() == t })
			if !find {
				return nil, fmt.Errorf("table %s not found in %s", t, db)
			}
			schemas = append(schemas, schema)
		}
	}

	for _, s := range schemas {
		createStmt, isMaterializedView, err := showCreateTable(ctx, conn, db, s.Name)
		if err != nil {
			return nil, err
		}

		s.CreateStmt = createStmt
		if isMaterializedView {
			s.Type = SchemaTypeMaterializedView
		}
	}

	return
}

func showCreateTable(ctx context.Context, conn *sqlx.DB, db, table string) (schema string, isMaterializedView bool, err error) {
	r, err := conn.QueryxContext(ctx, fmt.Sprintf(InternalSqlComment+"SHOW CREATE TABLE `%s`.`%s`", db, table))
	if err != nil {
		// may be a materialized view
		var err_ error
		r, err_ = conn.QueryxContext(ctx, fmt.Sprintf(InternalSqlComment+"SHOW CREATE MATERIALIZED VIEW `%s`.`%s`", db, table))
		if err_ != nil {
			return "", false, err
		}
		isMaterializedView = true
	}
	defer r.Close()

	logrus.Debugln("show create table:", table)

	schema, err = getStmtfromShowCreate(r)
	if err != nil {
		return "", false, err
	}

	// logrus.Traceln("create table:", schema)

	return
}

func getStmtfromShowCreate(r *sqlx.Rows) (schema string, err error) {
	cols, err := r.Columns()
	if err != nil {
		return "", err
	}
	vals := lo.ToAnySlice(lo.ToSlicePtr(make([]string, len(cols))))

	for r.Next() {
		err := r.Scan(vals...)
		if err != nil {
			return "", err
		}
		// the second column is the create statement
		schema = *vals[1].(*string)
	}
	if err := r.Err(); err != nil {
		return schema, err
	}

	return
}

func ShowDatabases(ctx context.Context, conn *sqlx.DB, dbnamePrefix string) ([]string, error) {
	dbs := []string{}
	err := conn.SelectContext(ctx, &dbs, InternalSqlComment+`SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME not in ('__internal_schema', 'information_schema', 'mysql') AND SCHEMA_NAME like ? ORDER BY SCHEMA_NAME`, SanitizeLike(dbnamePrefix)+"%")
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func ShowTables(ctx context.Context, conn *sqlx.DB, dbname string, tablenamePrefix ...string) (tables []*Schema, err error) {
	tables = []*Schema{}
	if len(tablenamePrefix) > 0 {
		err = conn.SelectContext(ctx, &tables, InternalSqlComment+`SELECT TABLE_NAME, TABLE_TYPE, TABLE_SCHEMA FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME like ? ORDER BY TABLE_NAME`, dbname, SanitizeLike(tablenamePrefix[0])+"%")
	} else {
		err = conn.SelectContext(ctx, &tables, InternalSqlComment+`SELECT TABLE_NAME, TABLE_TYPE, TABLE_SCHEMA FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME`, dbname)
	}
	if err != nil {
		return nil, err
	}
	for _, t := range tables {
		t.Type = t.Type.sanitize()
	}
	return
}

func ShowFronendsDisksDir(ctx context.Context, conn *sqlx.DB, diskType string) (dir string, err error) {
	r, err := conn.QueryxContext(ctx, InternalSqlComment+"show frontends DISKS")
	if err != nil {
		return "", err
	}
	defer r.Close()

	cols, err := r.Columns()
	if err != nil {
		return "", err
	}
	colDirTypeIdx := lo.IndexOf(cols, "DirType")
	colDirIdx := lo.IndexOf(cols, "Dir")
	vals := lo.ToAnySlice(lo.ToSlicePtr(make([]string, len(cols))))

	for r.Next() {
		err := r.Scan(vals...)
		if err != nil {
			return "", err
		}

		if *vals[colDirTypeIdx].(*string) == diskType {
			dir = *vals[colDirIdx].(*string)
			break
		}
	}
	if err := r.Err(); err != nil {
		return dir, err
	}

	return
}

func GetTablesStats(ctx context.Context, conn *sqlx.DB, dbname string, tables ...string) ([]*TableStats, error) {
	if len(tables) == 0 {
		return []*TableStats{}, nil
	}

	stats := make([]*TableStats, 0, len(tables))
	for _, table := range tables {
		s, err := getTableStats(ctx, conn, dbname, table)
		if err != nil {
			logrus.Errorf("get table stats failed: db: %s, table: %s, err: %v\n", dbname, table, err)
			return nil, err
		}
		if s == nil {
			continue
		}
		stats = append(stats, s)
	}

	return stats, nil
}

func getTableStats(ctx context.Context, conn *sqlx.DB, dbname, table string) (*TableStats, error) {
	logrus.Debugln("get table stats:", table)

	// show column stats of all table.
	r, err := conn.QueryxContext(ctx, InternalSqlComment+fmt.Sprintf("SHOW COLUMN STATS `%s`.`%s`", dbname, table))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	cols := []*ColumnStats{}
	for r.Next() {
		vals := map[string]any{}
		if err := r.MapScan(vals); err != nil {
			return nil, err
		}

		min, max := vals["min"].([]byte), vals["max"].([]byte)
		if bytes.HasPrefix(min, []byte(`'`)) {
			min = bytes.ReplaceAll(min[1:len(min)-1], []byte(`''`), []byte(`'`))
		}
		if bytes.HasPrefix(max, []byte(`'`)) {
			max = bytes.ReplaceAll(max[1:len(max)-1], []byte(`''`), []byte(`'`))
		}
		cols = append(cols, &ColumnStats{
			Name:        cast.ToString(vals["column_name"]),
			Count:       int64(cast.ToFloat64((string(vals["count"].([]byte))))),
			Ndv:         int64(cast.ToFloat64((string(vals["ndv"].([]byte))))),
			NullCount:   int64(cast.ToFloat64((string(vals["num_null"].([]byte))))),
			AvgSizeByte: int64(cast.ToFloat64((string(vals["avg_size_byte"].([]byte))))),
			DataSize:    int64(cast.ToFloat64((string(vals["data_size"].([]byte))))),
			Min:         string(min),
			Max:         string(max),
		})
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	if len(cols) == 0 {
		logrus.Warnf("no column stats found for %s.%s\n", dbname, table)
		return nil, nil
	}

	tbl := &TableStats{
		Name:     table,
		RowCount: cols[0].Count,
		Columns:  cols,
	}
	return tbl, nil
}

func CountAuditlogs(
	ctx context.Context,
	db *sqlx.DB,
	dbname, table string,
	opts AuditLogScanOpts,
) (int, error) {
	query := fmt.Sprintf("SELECT count(*) FROM `%s`.`%s` WHERE %s", dbname, table, opts.sqlConditions())
	logrus.Traceln("query from audit log table:", query)

	var total int
	err := db.GetContext(ctx, &total, InternalSqlComment+query)
	if err != nil {
		logrus.Errorln("query audit log count failed, err:", err)
	}
	return total, err
}

func GetDBAuditLogs(
	ctx context.Context,
	db *sqlx.DB,
	dbname, table string,
	opts AuditLogScanOpts,
	parallel int,
) ([]string, error) {
	total, err := CountAuditlogs(ctx, db, dbname, table, opts)
	if err != nil {
		return nil, err
	}
	if total <= 0 {
		logrus.Warnln("no audit log found")
		return []string{}, nil
	}
	if total > 1_000_000 {
		if !Confirm(fmt.Sprintf("Audit log count(%d) may bigger than 1 million, continue", total)) {
			return []string{}, nil
		}
	}

	logrus.Debugf("need to scan %d audit log row(s)\n", total)

	logScans := make([]*SimpleAuditLogScanner, parallel)
	for i := range logScans {
		s := NewSimpleAuditLogScanner(opts)
		s.Init()
		defer s.Close()
		logScans[i] = s
	}

	const LimitPerSelect = 100
	chunksNum := total / LimitPerSelect
	if total%LimitPerSelect != 0 {
		chunksNum += 1
	}
	if chunksNum < parallel {
		parallel = chunksNum
	}
	chunkPerThread := chunksNum / parallel
	chunkRemain := chunksNum % parallel

	g := ParallelGroup(parallel)
	conditions := opts.sqlConditions()
	var scanned int
	for i := 0; i < parallel; i++ {
		logScan := logScans[i]

		chunks := chunkPerThread
		if i < chunkRemain {
			chunks += 1
		}
		limitPerThread := chunks * LimitPerSelect

		start := scanned
		end := start + limitPerThread
		if end > total {
			end = total
		}
		scanned = end

		g.Go(func() error {
			pageConds := ""
			for offset := start; offset < end; offset += LimitPerSelect {
				limit := LimitPerSelect

				overflow := offset + limit - end
				if overflow > 0 {
					limit -= overflow
				}

				offset_ := offset
				if pageConds != "" {
					offset_ = 0
				}

				time, queryId, err := getDBAuditLogs(ctx, logScan, db, dbname, table, conditions+pageConds, limit, offset_)
				if err != nil {
					return err
				}
				pageConds = ""
				if time != "" && queryId != "" {
					pageConds = fmt.Sprintf(` AND time > "%s" AND query_id > "%s"`, time, queryId)
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return lo.Flatten(lo.Map(logScans, func(s *SimpleAuditLogScanner, _ int) []string { return s.Result() })), nil
}

func getDBAuditLogs(
	ctx context.Context,
	logScan *SimpleAuditLogScanner,
	db *sqlx.DB,
	dbname, table string,
	conditions string,
	limit, offset int,
) (string, string, error) {
	stmt := fmt.Sprintf("SELECT %s FROM `%s`.`%s` WHERE %s LIMIT %d OFFSET %d ORDER BY time asc, query_id asc",
		strings.Join(captureFieldCols, ", "),
		dbname,
		table,
		conditions,
		limit,
		offset,
	)
	logrus.Traceln("query audit log:", stmt)

	var (
		r     *sqlx.Rows
		err   error
		retry int
	)
	const MaxRetry = 3
	for ; retry < MaxRetry; retry++ {
		r, err = db.QueryxContext(ctx, InternalSqlComment+stmt)
		if err == nil {
			break
		}
	}
	if err != nil {
		logrus.Errorf("query audit log table failed: retry: %d, db: %s, table: %s, err: %v\n", retry, dbname, table, err)
		return "", "", err
	}
	defer r.Close()

	var lastTime, lastQueryId string
	for i := 0; r.Next(); i++ {
		vals_, err := r.SliceScan()
		if err != nil {
			return "", "", err
		}

		vals, err := cast.ToStringSliceE(vals_)
		if err != nil {
			logrus.Errorf("read audit log table failed: db: %s, table: %s, err: %v\n", dbname, table, err)
			return "", "", err
		}
		lastTime, lastQueryId = vals[0], vals[5]
		logScan.onMatch(vals, true)
	}

	return lastTime, lastQueryId, r.Err()
}

func SanitizeLike(s string) string {
	return sqlLikeReplacer.Replace(s)
}
