package src

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

var (
	InternalSqlComment      = "/*dorisdump*/"
	InternalSqlCommentBytes = []byte(InternalSqlComment)
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

type Schema struct {
	Name       string     `db:"TABLE_NAME"`
	Type       SchemaType `db:"TABLE_TYPE"`
	DB         string     `db:"TABLE_SCHEMA"`
	CreateStmt string     `db:"-"`
}

func (s *Schema) String() string {
	return fmt.Sprintf("%s.%s", s.DB, s.Name)
}

func NewDB(host string, port int16, user, password, db string) (*sqlx.DB, error) {
	cfg := &mysql.Config{
		User:                 user,
		Passwd:               password,
		Addr:                 fmt.Sprintf("%s:%d", host, port),
		Net:                  "tcp",
		DBName:               db,
		AllowNativePasswords: true,
		Timeout:              10 * time.Second,
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
	schemas_, err := showTables(ctx, conn, db)
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
	r, err := conn.QueryxContext(ctx, fmt.Sprintf(InternalSqlComment+`SHOW CREATE TABLE %s.%s`, db, table))
	if err != nil {
		// may be a materialized view
		var err_ error
		r, err_ = conn.QueryxContext(ctx, fmt.Sprintf(InternalSqlComment+`SHOW CREATE MATERIALIZED VIEW %s.%s`, db, table))
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

	return
}

func showTables(ctx context.Context, conn *sqlx.DB, dbname string) (tables []*Schema, err error) {
	tables = []*Schema{}
	err = conn.SelectContext(ctx, &tables, InternalSqlComment+`SELECT lower(TABLE_NAME) as TABLE_NAME, TABLE_TYPE, TABLE_SCHEMA FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?`, dbname)
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

	return
}
