package src

import (
	"bufio"
	"fmt"
	"strings"

	"dario.cat/mergo"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	gen "github.com/Thearas/dorisdump/src/generator"
	"github.com/Thearas/dorisdump/src/parser"
)

const (
	ColumnSeparator    = '☆' // make me happy
	DefaultGenRowCount = 1000
	MaxGenRowCount     = 1_000_000
)

func NewTableGen(createTableStmt string, stats *TableStats) (*TableGen, error) {
	// parse create-table statement
	sqlId := "create-table"
	if stats != nil {
		sqlId = stats.Name
	}
	p := parser.NewParser(sqlId, createTableStmt)
	c, ok := p.SupportedCreateStatement().(*parser.CreateTableContext)
	if !ok {
		logrus.Fatalln("SQL parser error")
	} else if p.ErrListener.LastErr != nil {
		return nil, p.ErrListener.LastErr
	}

	// get table stats
	table := strings.ReplaceAll(strings.ReplaceAll(c.GetName().GetText(), "`", ""), " ", "")
	colStats := make(map[string]*ColumnStats)
	if stats != nil {
		colStats = lo.SliceToMap(stats.Columns, func(s *ColumnStats) (string, *ColumnStats) {
			s.Count = stats.RowCount
			return s.Name, s
		})
		logrus.Debugf("using stats for table '%s'\n", table)
	} else {
		logrus.Debugf("stats not found for table '%s'\n", table)
	}

	// get custom table gen rule
	rows, customColumnRule := gen.GetCustomTableGenRule(table)

	// construct every streamLoadCols
	streamLoadCols := make([]string, 0, len(c.ColumnDefs().GetCols()))
	colGens := make([]gen.Gen, 0, len(c.ColumnDefs().GetCols()))
	hasBitmap := false
	for _, col := range c.ColumnDefs().GetCols() {
		colName := strings.Trim(col.GetColName().GetText(), "`")
		colType_ := col.GetType_()
		visitor := &gen.TypeVisitor{Colpath: fmt.Sprintf("%s.%s", table, colName)}
		colBaseType := visitor.GetBaseType(colType_)

		if colBaseType == "BITMAP" {
			hasBitmap = true
			streamLoadCols = append(streamLoadCols, fmt.Sprintf("raw_%s,`%s`=bitmap_from_array(cast(raw_%s as ARRAY<BIGINT(20)>))", colName, colName, colName))
		} else {
			streamLoadCols = append(streamLoadCols, "`"+colName+"`")
		}

		// get column gen rule
		visitor.GenRule = newColGenRule(col, colName, colBaseType, colStats, customColumnRule)

		// build column generator
		colGens = append(colGens, visitor.GetTypeGen(colType_))
	}

	tg := &TableGen{rows: rows, colGens: colGens}
	if hasBitmap {
		tg.streamloadColumns = "columns:" + strings.Join(streamLoadCols, ",")
	}

	return tg, nil
}

func newColGenRule(col parser.IColumnDefContext, colName, colType string, colStats map[string]*ColumnStats, customColumnRule map[string]GenRule) GenRule {
	genRule := GenRule{}

	// 1. Merge rules in stats
	if colstats, ok := colStats[colName]; ok {
		var nullFreq float32
		if colstats.Count > 0 {
			nullFreq = float32(colstats.NullCount) / float32(colstats.Count)
		}
		if nullFreq >= 0 && nullFreq < 1 {
			genRule["null_frequency"] = nullFreq
		}

		if IsStringType(colType) {
			avgLen := colstats.AvgSizeByte
			genRule["length"] = avgLen

			// HACK: +-5 on string avg size as length
			if avgLen > 5 && colType != "CHAR" {
				genRule["length"] = GenRule{
					"min": avgLen - 5,
					"max": avgLen + 5,
				}
			}
		} else {
			if colstats.Min != "" {
				genRule["min"] = colstats.Min
			}
			if colstats.Max != "" {
				genRule["max"] = colstats.Max
			}
		}
	}

	// 2. Merge rules in global custom rules
	customRule, ok := customColumnRule[colName]
	if !ok || len(customRule) == 0 {
		return genRule
	}
	if err := mergo.Merge(&genRule, customRule, mergo.WithOverride); err != nil {
		logrus.Fatalln(err)
	}

	notnull := col.NOT() != nil && col.GetNullable() != nil
	if notnull {
		genRule["null_frequency"] = 0
	}

	return genRule
}

type TableGen struct {
	streamloadColumns string
	rows              int
	colGens           []gen.Gen
}

// Gen generates multiple CSV line into writer.
func (tg *TableGen) GenCSV(w *bufio.Writer, rows int) error {
	if tg.streamloadColumns != "" {
		if _, err := w.WriteString(tg.streamloadColumns); err != nil {
			return err
		}
		w.WriteByte('\n')
	}
	if rows <= 0 {
		rows = tg.rows
	}
	if rows == 0 {
		rows = DefaultGenRowCount
	}
	if err := CheckGenRowCount(rows); err != nil {
		return err
	}

	for l := range rows {
		tg.genOne(w)
		if l != rows-1 {
			if err := w.WriteByte('\n'); err != nil {
				return err
			}
		}
	}
	return nil
}

// GenOne generates one CSV line into writer.
func (tg *TableGen) genOne(w *bufio.Writer) {
	for i, g := range tg.colGens {
		val := g.Gen()
		if val == nil {
			w.WriteString(`\N`)
		} else if v, ok := val.(json.RawMessage); ok {
			w.Write(v)
		} else if s, ok := val.(string); ok {
			w.WriteString(s)
		} else {
			w.WriteString(fmt.Sprint(val))
		}
		if i != len(tg.colGens)-1 {
			w.WriteRune(ColumnSeparator)
		}
	}
}

type GenRule = gen.GenRule

func CheckGenRowCount(rows int) error {
	if rows < 0 {
		return fmt.Errorf("--rows/row_count must be a positive integer, got %d", rows)
	} else if rows > MaxGenRowCount {
		return fmt.Errorf("--rows/row_count must be smaller than 100_000, got %d", rows)
	}
	return nil
}
