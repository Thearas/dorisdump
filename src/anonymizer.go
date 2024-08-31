package src

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
	"golang.org/x/sync/errgroup"

	"github.com/Thearas/dorisdump/src/parser"
)

const (
	AnonymizeMinLength = 3
	AnonymizeHashBytes = 4
	AnonymizeHashFmt   = "h%x" // add prefix 'h'
)

func AnonymizeSqlsInPlace(method string, sqls []string) {
	// check vaild
	anonymizeF := getAnonymizeFunc(method)
	if anonymizeF == nil {
		return
	}

	g := errgroup.Group{}
	g.SetLimit(10)

	for i := range sqls {
		i, sql := i, sqls[i]
		g.Go(func() error {
			sqls[i] = AnonymizeSql(method, sql)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logrus.Warnln("Anonymization sql error:", err)
	}
}

func AnonymizeSql(method string, sql string) string {
	anonymizeF := getAnonymizeFunc(method)
	if anonymizeF == nil {
		return sql
	}

	p := parser.NewParser(sql, parser.NewListener(true, anonymizeF))

	return p.ToSQL()
}

func Anonymize(method string, s string) string {
	return getAnonymizeFunc(method)(s, false)
}

// NOTE: not thread safe.
func getAnonymizeFunc(method string) func(string, bool) string {
	h := blake3.New()
	hashF := func(id string, ignoreBuiltin bool) string {
		// do not anoymize identifier that is less than MinAnonymizeLength characters.
		if len(id) < AnonymizeMinLength {
			return id
		}

		// FIXME: db name is case-insensitive
		lowerid := strings.ToLower(id)

		// only take the first AnonymizeHashBytes bytes of hash.
		b := hashstr(h, lowerid)
		hash := fmt.Sprintf(AnonymizeHashFmt, b[:AnonymizeHashBytes])

		// do not anoymize builtin functions.
		if _, ok := BuiltinFunctionHashs[hash]; ignoreBuiltin && ok {
			return id
		}

		return hash
	}

	switch method {
	case "hash":
		return hashF
	default:
		logrus.Warnf("Anonymization method %s is not supported, keep going with no anonymization\n", method)
		return nil
	}
}
