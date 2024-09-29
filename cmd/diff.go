/*
Copyright © 2024 Thearas thearas850@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Thearas/dorisdump/src"
)

var (
	noColor         bool
	minDurationDiff time.Duration
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:     "diff",
	Short:   "Diff replay result",
	Example: "dorisdump diff /path/to/replay1 /path/to/replay2",
	RunE: func(cmd *cobra.Command, args []string) error {
		if noColor {
			if err := os.Setenv("NO_COLOR", "true"); err != nil {
				return err
			}
		}

		if len(args) != 2 {
			return errors.New("diff requires two file path")
		}

		replay1, replay2 := args[0], args[1]
		lstats, err := os.Stat(replay1)
		if err != nil {
			return err
		}
		rstats, err := os.Stat(replay2)
		if err != nil {
			return err
		}
		if lstats.IsDir() != rstats.IsDir() {
			return errors.New("file path should be both file or both directory")
		}

		return diff(replay1, replay2, lstats.IsDir())
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	flags := diffCmd.Flags()
	flags.BoolVar(&noColor, "no-color", false, "Disable color output")
	flags.DurationVar(&minDurationDiff, "min-duration-diff", 100*time.Millisecond, "Print diff if duration difference is greater than this value")
}

func diff(replay1, replay2 string, isDir bool) error {
	if !isDir {
		return diffFile(replay1, replay2)
	}

	return filepath.WalkDir(replay1, func(path1 string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path1, src.ReplayResultFileExt) {
			return nil
		}

		relativePath := strings.TrimPrefix(path1, replay1)
		path2 := filepath.Join(replay2, relativePath)

		if err := diffFile(path1, path2); err != nil {
			logrus.Errorf("diff %s and %s failed, err: %v\n", path1, path2, err)
		}
		return nil
	})
}

func diffFile(file1, file2 string) error {
	f1, err := os.Open(file1)
	if err != nil {
		return err
	}
	defer f1.Close()
	scan1 := bufio.NewScanner(f1)
	f2, err := os.Open(file2)
	if err != nil {
		return err
	}
	defer f2.Close()
	scan2 := bufio.NewScanner(f2)

	logrus.Debugf("diffing %s and %s\n", file1, file2)

	id2diff := make(map[string]string)
	for scan1.Scan() {
		b1 := scan1.Bytes()

		var b2 []byte
		if scan2.Scan() {
			b2 = scan2.Bytes()
		}

		if len(b1) == 0 && len(b2) == 0 {
			continue
		}

		var d diff2
		err = json.Unmarshal(b1, &d.r1)
		if err != nil {
			logrus.Errorf("unmarshal %s failed, err: %v\n", string(b2), err)
		}
		err = json.Unmarshal(b2, &d.r2)
		if err != nil {
			logrus.Errorf("unmarshal %s failed, err: %v\n", string(b2), err)
		}

		if d.r1.QueryId != d.r2.QueryId {
			id2diff[d.r1.QueryId] = fmt.Sprintf("query id not match, %s != %s", d.r1.QueryId, d.r2.QueryId)
			break
		}
		if diffmsg := d.result(); diffmsg != "" {
			id2diff[d.r1.QueryId] = diffmsg
		}
	}

	// print diff result
	for id, diffmsg := range id2diff {
		fmt.Printf("QueryId: %s, %s\n", id, diffmsg)
	}

	return nil
}

type diff2 struct {
	r1, r2 src.ReplayResult
	diff   string
}

func (d *diff2) result() string {
	if d.r1.Err != d.r2.Err {
		return fmt.Sprintf(`err not match:
%s
------
%s`, color.GreenString(d.r1.Err), color.RedString(d.r2.Err))
	}
	if d.r1.ReturnRows != d.r2.ReturnRows {
		return fmt.Sprintf("rows count not match: %s != %s", color.GreenString(strconv.Itoa(d.r1.ReturnRows)), color.RedString(strconv.Itoa(d.r2.ReturnRows)))
	}
	if d.r1.ReturnRowsHash != d.r2.ReturnRowsHash {
		return color.RedString("rows hash not match (count: %d)", d.r1.ReturnRows)
	}
	if math.Abs(float64(d.r1.DurationMs-d.r2.DurationMs)) > float64(minDurationDiff.Milliseconds()) {
		return fmt.Sprintf("duration not match: %s vs %s", color.GreenString("%dms", d.r1.DurationMs), color.RedString("%dms", d.r2.DurationMs))
	}
	return ""
}
