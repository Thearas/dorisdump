package generator

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// TODO: use disk to store refVals, now is in-memory impl.

var (
	_ Gen = &RefGen{}

	// table -> column -> refgen
	refGenMap     = map[string]map[string]*RefGen{}
	refGenMapLock sync.Mutex
)

func GetTableRefGen(table string) map[string]*RefGen {
	return refGenMap[table]
}

func getColumnRefGen(table, column string) *RefGen {
	tc := GetTableRefGen(table)
	if len(tc) == 0 {
		return nil
	}
	return tc[column]
}

type RefGen struct {
	Table  string
	Column string
	Limit  int

	refValsPtr *[]any

	possibility float32
}

func (g *RefGen) Clone() *RefGen {
	return &RefGen{
		Table:      g.Table,
		Column:     g.Column,
		Limit:      g.Limit,
		refValsPtr: g.refValsPtr,
	}
}

func (g *RefGen) TableColumn() string {
	return fmt.Sprintf("%s.%s", g.Table, g.Column)
}

// rowCount is the source table row count that ref points to.
func (g *RefGen) WithSourceTableRows(rowCount int) {
	g.possibility = float32(g.Limit) / float32(rowCount)
}

func (g *RefGen) AddRefVals(vals ...any) (full bool) {
	if g.possibility <= 0 {
		panic("unreachable: ref possibility must > 0")
	}
	refVals := *g.refValsPtr
	for _, v := range vals {
		if len(refVals) >= g.Limit {
			return true
		}

		// need at least one value in refVals
		pickup := len(refVals) == 0 || g.possibility >= 1 || rand.Float32() <= g.possibility
		if !pickup {
			continue
		}

		*g.refValsPtr = append(*g.refValsPtr, v)
	}

	return false
}

func (g *RefGen) Gen() any {
	refVals := *g.refValsPtr
	if len(refVals) == 0 {
		logrus.Fatalln("empty ref value point to", g.TableColumn())
	}

	limit := int(math.Min(float64(g.Limit), float64(len(refVals))))

	return refVals[gofakeit.IntN(limit)]
}

func NewRefGenerator(_ string, r GenRule) (Gen, error) {
	refGenMapLock.Lock()
	defer refGenMapLock.Unlock()

	tableColumn_ := cast.ToString(r["ref"])
	tableColumn := strings.SplitN(tableColumn_, ".", 2)
	if len(tableColumn) != 2 {
		return nil, fmt.Errorf("wrong ref, expect '<table>.<column>', got '%s'", tableColumn_)
	}

	limit := 1000
	if l := cast.ToInt(r["limit"]); l > 0 {
		limit = l
	}

	g := &RefGen{
		Table:      tableColumn[0],
		Column:     tableColumn[1],
		Limit:      limit,
		refValsPtr: &[]any{},
	}

	var sharedRefGen *RefGen

	c2g, ok := refGenMap[g.Table]
	if ok {
		sharedRefGen, ok = c2g[g.Column]
		if ok {
			// use the biggest limit
			sharedRefGen.Limit = int(math.Max(float64(g.Limit), float64(sharedRefGen.Limit)))

			// share the same refVals for all ref which ref to the same table.column
			g.refValsPtr = sharedRefGen.refValsPtr

			return g, nil
		}
	} else {
		refGenMap[g.Table] = map[string]*RefGen{}
	}

	sharedRefGen = g.Clone()
	refGenMap[g.Table][g.Column] = sharedRefGen

	return g, nil
}
