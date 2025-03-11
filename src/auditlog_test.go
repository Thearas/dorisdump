package src

import (
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type sqlWriter struct {
	sqls []string
}

func (w *sqlWriter) WriteSql(s string) error {
	w.sqls = append(w.sqls, s)
	return nil
}

func (w *sqlWriter) Close() error {
	return nil
}

func TestExtractQueriesFromAuditLogs(t *testing.T) {
	t.Parallel()
	chroot()
	disableLog()

	type args struct {
		dbs               []string
		auditlogPaths     []string
		encoding          string
		queryMinCpuTimeMs int64
		queryStates       []string
		parallel          int
		unescape          bool
		onlySelect        bool
		strict            bool
		from, to          string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				auditlogPaths:     []string{"fixture/fe.audit.log"},
				encoding:          "auto",
				queryMinCpuTimeMs: 8,
				unescape:          true,
				onlySelect:        true,
				strict:            true,
			},
			want: 8,
		},
		{
			name: "not_only_select",
			args: args{
				auditlogPaths: []string{"fixture/fe.audit.log"},
				encoding:      "auto",
				unescape:      true,
				onlySelect:    false,
				strict:        true,
			},
			want: 9,
		},
		{
			name: "from_to",
			args: args{
				auditlogPaths: []string{"fixture/fe.audit.log"},
				encoding:      "auto",
				unescape:      true,
				onlySelect:    false,
				strict:        true,
				from:          "2024-08-06 23:44:11",
				to:            "2024-08-06 23:44:12",
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := AuditLogScanOpts{
				DBs:                tt.args.dbs,
				QueryMinDurationMs: tt.args.queryMinCpuTimeMs,
				QueryStates:        tt.args.queryStates,
				OnlySelect:         tt.args.onlySelect,
				From:               tt.args.from,
				To:                 tt.args.to,
				Unescape:           tt.args.unescape,
				Strict:             tt.args.strict,
			}
			writers := lo.RepeatBy(len(tt.args.auditlogPaths), func(index int) SqlWriter { return &sqlWriter{} })
			gotCount, err := ExtractQueriesFromAuditLogs(writers, tt.args.auditlogPaths, tt.args.encoding, opts, tt.args.parallel)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractQueriesFromAuditLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, sql := range writers[0].(*sqlWriter).sqls {
				assert.Contains(t, sql, `"user":"root"`)
				assert.True(t, strings.Contains(sql, `"db":"mydb"`) || strings.Contains(sql, `"db":"__internal_schema"`))
			}
			if !reflect.DeepEqual(gotCount, tt.want) {
				t.Errorf("ExtractQueriesFromAuditLogs() = %v, want %v", gotCount, tt.want)
			}
		})
	}
}

func chroot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
}

func disableLog() {
	logrus.SetLevel(logrus.ErrorLevel)
}
