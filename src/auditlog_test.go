package src

import (
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func TestExtractQueriesFromAuditLogs(t *testing.T) {
	t.Parallel()
	chroot()
	disableLog()

	type args struct {
		dbs               []string
		auditlogPaths     []string
		encoding          string
		queryMinCpuTimeMs int
		queryStates       []string
		parallel          int
		unique            bool
		uniqueNormalize   bool
		unescape          bool
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				auditlogPaths:     []string{"fixture/fe.audit.log"},
				encoding:          "auto",
				queryMinCpuTimeMs: 8,
				unescape:          true,
			},
			want: []int{8},
		},
		{
			name: "unique",
			args: args{
				auditlogPaths:   []string{"fixture/fe.audit.log"},
				encoding:        "auto",
				unique:          true,
				uniqueNormalize: true,
				unescape:        true,
			},
			want: []int{7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractQueriesFromAuditLogs(tt.args.dbs, tt.args.auditlogPaths, tt.args.encoding, tt.args.queryMinCpuTimeMs, tt.args.queryStates, tt.args.parallel, tt.args.unique, tt.args.uniqueNormalize, tt.args.unescape)
			gotCount := lo.Map(got, func(s []string, _ int) int { return len(s) })
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractQueriesFromAuditLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
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
