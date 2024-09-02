package src

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"

	"github.com/Thearas/dorisdump/src/parser"
)

const (
	AnonymizeHashBytes = 4
	AnonymizeHashFmt   = "h%x" // add prefix 'h'
)

var (
	AnonymizeMinLength       int
	AnonymizerreserveIdHashs map[string]string

	// Identifiers that should not be anonymized.
	ReserveIdentifiers = lo.Map([]string{
		"olap",
		"internal",
		"__internal_schema",
		"information_schema",
	}, func(s string, _ int) string { return strings.ToLower(s) })
)

func SetupAnonymizer(idMinLength int, reserveIds ...string) {
	AnonymizeMinLength = idMinLength
	SetupBuiltinHashs()

	reserveIds = append(ReserveIdentifiers, reserveIds...)
	AnonymizerreserveIdHashs = anonymizeHashSliceToMap(reserveIds)
}

func AnonymizeSqlsInPlace(method string, sqls []string, parallel int) {
	// check vaild
	anonymizeF := getAnonymizeFunc(method)
	if anonymizeF == nil {
		return
	}

	g := ParallelGroup(parallel)

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
		hash := anonymizeHashStr(h, lowerid)

		// do not anoymize reserve ids.
		if _, ok := AnonymizerreserveIdHashs[hash]; ok {
			return id
		}
		// do not anoymize builtin functions.
		if _, ok := BuiltinHashs[hash]; ignoreBuiltin && ok {
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

func anonymizeHashSliceToMap(xs []string) map[string]string {
	h := blake3.New()
	return lo.SliceToMap(xs, func(s string) (string, string) {
		return anonymizeHashStr(h, s), s
	})
}

func anonymizeHashStr(h *blake3.Hasher, s string) string {
	b := hashstr(h, s)

	return fmt.Sprintf(AnonymizeHashFmt, b[:AnonymizeHashBytes])
}
