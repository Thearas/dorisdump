package src

import (
	"bufio"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDecodeReplaySqls(t *testing.T) {
	t.Parallel()
	chroot()
	disableLog()

	replayFile, err := os.Open("fixture/q0.sql")
	assert.NoError(t, err)
	defer replayFile.Close()

	minTs, err := time.Parse("2006-01-02 15:04:05,000", "2024-08-06 23:44:11,041")
	assert.NoError(t, err)

	type args struct {
		s        *bufio.Scanner
		dbs      map[string]struct{}
		users    map[string]struct{}
		from     int64
		to       int64
		maxCount int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		want1   int64
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				s: bufio.NewScanner(replayFile),
			},
			want: map[string]int{
				"192.168.48.119:51970": 4,
				"192.168.48.118:51970": 4,
			},
			want1: minTs.UnixMilli(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeReplaySqls(tt.args.s, tt.args.dbs, tt.args.users, tt.args.from, tt.args.to, tt.args.maxCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeReplaySqls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotCount := lo.MapValues(got, func(v []*ReplaySql, _ string) int {
				return len(v)
			})
			if !reflect.DeepEqual(gotCount, tt.want) {
				t.Errorf("DecodeReplaySqls() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeReplaySqls() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
