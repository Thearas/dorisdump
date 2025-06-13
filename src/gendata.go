package src

import (
	"bufio"
	"bytes"
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/Thearas/dorisdump/src/parser"
)

const (
	ColumnSeparator    = '☆'
	DefaultGenRowCount = 1000
	MaxGenRowCount     = 100_000
)

var (
	TypeAlias = map[string]string{
		"INTEGER":    "INT",
		"TEXT":       "STRING",
		"BOOL":       "BOOLEAN",
		"DECIMALV2":  "DECIMAL",
		"DECIMALV3":  "DECIMAL",
		"DATEV1":     "DATE",
		"DATEV2":     "DATE",
		"DATETIMEV1": "DATETIME",
		"DATETIMEV2": "DATETIME",
		"TIMESTAMP":  "DATETIME",
	}
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
	rows, customColumnRule := getCustomTableGenRule(table)

	// construct every columns
	columns := make([]string, 0, len(c.ColumnDefs().GetCols()))
	colGens := make([]Gen, 0, len(c.ColumnDefs().GetCols()))
	hasBitmap := false
	for _, col := range c.ColumnDefs().GetCols() {
		colName := strings.Trim(col.GetColName().GetText(), "`")
		visitor := &TypeVisitor{Colpath: fmt.Sprintf("%s.%s", table, colName)}
		colType_ := col.GetType_()
		colBaseType := visitor.GetBastType(colType_)

		// build column gen rule
		genRule := newColGenRule(col, colName, colBaseType, colStats, customColumnRule)
		visitor.GenRule = genRule

		if colBaseType == "BITMAP" {
			hasBitmap = true
			columns = append(columns, fmt.Sprintf("raw_%s,`%s`=bitmap_from_array(cast(raw_%s as ARRAY<BIGINT(20)>))", colName, colName, colName))
		} else {
			columns = append(columns, "`"+colName+"`")
		}
		colGens = append(colGens, visitor.GetTypeGen(colType_))
	}

	tg := &TableGen{rows: rows, colGens: colGens}
	if hasBitmap {
		tg.streamloadColumns = "columns:" + strings.Join(columns, ",")
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
	colGens           []Gen
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
		} else {
			w.WriteString(fmt.Sprint(val))
		}
		if i != len(tg.colGens)-1 {
			w.WriteRune(ColumnSeparator)
		}
	}
}

type GenRule = map[string]any

type TypeVisitor struct {
	Colpath string  // the path of the column, e.g. "db.table.col"
	GenRule GenRule // rules of generator
}

func NewTypeVisitor(colpath string, genRule GenRule) *TypeVisitor {
	if genRule == nil {
		genRule = GenRule{}
	}
	return &TypeVisitor{
		Colpath: colpath,
		GenRule: genRule,
	}
}

func (v *TypeVisitor) MergeDefaultRule(baseType string) *TypeVisitor {
	defaultGenRule, ok := DefaultTypeGenRules[baseType]
	if !ok {
		if ty_, ok := TypeAlias[baseType]; ok {
			baseType = ty_
		}
		defaultGenRule, ok = DefaultTypeGenRules[baseType]
		if !ok {
			return v
		}
	}
	if len(defaultGenRule) == 0 {
		return v
	}

	if err := mergo.Merge(&v.GenRule, defaultGenRule); err != nil {
		logrus.Fatalf("Unable to merge default gen rule for type '%s' in column '%s', err: %v\n", baseType, v.Colpath, err)
	}

	return v
}

func (v *TypeVisitor) HasGenRule() bool {
	return len(v.GenRule) > 0
}

func (v *TypeVisitor) GetRule(name string, defaultValue ...any) any {
	if !v.HasGenRule() {
		return nil
	}
	if r, ok := v.GenRule[name]; ok {
		return r
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (v *TypeVisitor) GetMinMax() (min, max any) {
	return v.GetRule("min"), v.GetRule("max")
}

func (v *TypeVisitor) GetLength() (min, max int) {
	l := v.GetRule("length")
	if l == nil {
		logrus.Fatalf("length not found for column '%s'\n", v.Colpath)
	}

	switch l := l.(type) {
	case int, float32, float64:
		length := cast.ToInt(l)
		min, max = length, length
	case GenRule:
		min, max = cast.ToInt(l["min"]), cast.ToInt(l["max"])
	}
	if max < min {
		logrus.Debugf("length max(%d) < min(%d), set max to min for column '%s'\n", max, min, v.Colpath)
		min = max
	}
	return
}

func (v *TypeVisitor) ChildGenRule(name string) GenRule {
	r := v.GetRule(name)
	if r == nil {
		return nil
	}
	return r.(GenRule)
}

func (v *TypeVisitor) GetChildGenRule(name string, childType parser.IDataTypeContext) Gen {
	return NewTypeVisitor(v.Colpath+"."+name, v.ChildGenRule(name)).GetTypeGen(childType)
}

func (v *TypeVisitor) GetNullFrequency() float32 {
	nullFrequency, err := cast.ToFloat32E(v.GetRule("null_frequency", GLOBAL_NULL_FREQUENCY))
	if err != nil || nullFrequency < 0 || nullFrequency > 1 {
		logrus.Fatalf("Invalid null frequency '%v' for column '%s': %v\n", v.GetRule("null_frequency"), v.Colpath, err)
	}
	return nullFrequency
}

func (v *TypeVisitor) GetBastType(type_ parser.IDataTypeContext) (t string) {
	switch ty := type_.(type) {
	case *parser.ComplexDataTypeContext:
		t = ty.GetComplex_().GetText()
	case *parser.PrimitiveDataTypeContext:
		t = ty.PrimitiveColType().GetType_().GetText()
	default:
		logrus.Fatalf("Unsupported column type '%s' for column '%s'\n", type_.GetText(), v.Colpath)
	}
	return strings.ToUpper(t)
}

func (v *TypeVisitor) GetTypeGen(type_ parser.IDataTypeContext) Gen {
	baseType := v.GetBastType(type_)
	v.MergeDefaultRule(baseType) // Merge global (aka. default) generate rules first.
	if logrus.GetLevel() > logrus.DebugLevel {
		logrus.Tracef("gen rule of '%s': %s\n", v.Colpath, string(MustJsonMarshal(v.GenRule)))
	}

	var (
		nullFrequency = v.GetNullFrequency()
		g             Gen
	)

	switch ty := type_.(type) {
	case *parser.ComplexDataTypeContext:
		switch baseType {
		case "ARRAY":
			// Handle array type
			g_ := &ArrayGen{}
			g_.LenMin, g_.LenMax = v.GetLength()
			g_.SetElementGen(v.GetChildGenRule("element", ty.DataType(0)))
			g = g_
		case "MAP":
			// Handle map type
			kv := ty.AllDataType()
			if len(kv) != 2 {
				logrus.Fatalf("Invalid map type: '%s' for column '%s', expected 2 types for key and value\n", ty.GetText(), v.Colpath)
			}

			// Handle key-value pair in map
			g_ := &MapGen{}
			g_.LenMin, g_.LenMax = v.GetLength()
			g_.SetKeyGen(v.GetChildGenRule("key", kv[0]))
			g_.SetValueGen(v.GetChildGenRule("value", kv[1]))
			g = g_
		case "STRUCT":
			// Handle struct type
			g_ := &StructGen{}

			// Handle each field in the struct
			fields_ := v.GetRule("fields")
			if fields_ == nil {
				fields_ = v.GetRule("field")
			}
			fieldRules, ok := fields_.([]GenRule) // Ensure fields is a slice of maps
			if !ok {
				if fields_ != nil {
					logrus.Fatalf("Invalid struct fields type '%T' for column '%s'\n", fields_, v.Colpath)
				}
				fieldRules = []GenRule{}
			}
			i := 0
			fields := lo.SliceToMap(fieldRules, func(field GenRule) (string, GenRule) {
				fieldName, ok := field["name"].(string)
				if !ok {
					logrus.Fatalf("Struct field #%d has no name in column '%s'\n", i, v.Colpath)
				}
				i++
				return fieldName, field
			})
			for _, field := range ty.ComplexColTypeList().AllComplexColType() {
				fieldName := strings.Trim(field.Identifier().GetText(), "`")
				fieldType := field.DataType()
				fieldGenRule, ok := fields[fieldName]
				if !ok {
					fieldGenRule = nil
				}
				fieldVisitor := NewTypeVisitor(v.Colpath+"."+fieldName, fieldGenRule)
				g_.AddChild(fieldName, fieldVisitor.GetTypeGen(fieldType))
			}
			g = g_
		default:
			logrus.Fatalf("Unsupported complex type: '%s' for column '%s'\n", ty.GetComplex_().GetText(), v.Colpath)
		}
	case *parser.PrimitiveDataTypeContext:
		min_, max_ := v.GetMinMax()
		switch baseType {
		case "BITMAP":
			// Generate a random bitmap array with a length between lenMin and lenMax
			lenMin, lenMax := v.GetLength()
			min, max := CastMinMax[int64](min_, max_, baseType, v.Colpath)
			g = NewFuncGen(func() any {
				return json.RawMessage(MustJsonMarshal(lo.RepeatBy(gofakeit.IntRange(lenMin, lenMax), func(_ int) int64 {
					return rand.Int64N(max-min+1) + min
				})))
			})
		case "JSON", "JSONB", "VARIANT":
			var genRule GenRule
			structure, ok := v.GetRule("structure").(string)
			structure = strings.TrimSpace(structure)
			if ok && structure != "" {
				genRule = maps.Clone(v.GenRule)
				delete(genRule, "structure")
			} else {
				logrus.Fatalf("JSON/JSONB/VARIANT must have gen rule 'structure at column '%s'\n", v.Colpath)
			}

			p := parser.NewParser(v.Colpath, structure)
			dataType := p.DataType()
			if err := p.ErrListener.LastErr; err != nil {
				logrus.Fatalf("Invalid JSON structure '%s' for column '%s': %v\n", structure, v.Colpath, err)
			}
			visitor := NewTypeVisitor(v.Colpath, genRule)
			g = visitor.GetTypeGen(dataType)
		case "BOOL", "BOOLEAN":
			enum := []int{0, 1}
			g = NewFuncGen(func() any { return gofakeit.RandomInt(enum) }) // BOOLEAN is typically 0 or 1
		case "TINYINT":
			min, max := CastMinMax[int8](min_, max_, baseType, v.Colpath)
			g = NewIntGen(min, max)
		case "SMALLINT":
			min, max := CastMinMax[int16](min_, max_, baseType, v.Colpath)
			g = NewIntGen(min, max)
		case "INT", "INTEGER":
			min, max := CastMinMax[int32](min_, max_, baseType, v.Colpath)
			g = NewIntGen(min, max)
		case "BIGINT", "LARGEINT": // TODO: Need larger INT?
			min, max := CastMinMax[int64](min_, max_, baseType, v.Colpath)
			range_ := max - min + 1
			g = NewFuncGen(func() int64 { return rand.Int64N(range_) + min })
		case "FLOAT":
			min, max := CastMinMax[float32](min_, max_, baseType, v.Colpath)
			g = NewFuncGen(func() any { return gofakeit.Float32Range(min, max) })
		case "DOUBLE":
			min, max := CastMinMax[float64](min_, max_, baseType, v.Colpath)
			g = NewFuncGen(func() any { return gofakeit.Float64Range(min, max) })
		case "DECIMAL", "DECIMALV2", "DECIMALV3": // TODO: Need larger DECIMAL?
			var precision, scale int = 999, 999
			if v.GetRule("precision") != nil {
				precision = cast.ToInt(v.GetRule("precision"))
			}
			if v.GetRule("scale") != nil {
				scale = cast.ToInt(v.GetRule("scale"))
			}

			intVals := ty.AllINTEGER_VALUE()
			p := cast.ToInt(intVals[0].GetText())
			if p > 38 {
				p = 38
			}
			if precision > p {
				precision = p
				// logrus.Debugf("Precision '%d' is larger than the defined precision '%d' for column '%s', using %d instead\n", precision, p, v.Colpath, p)
			}
			if len(intVals) > 1 {
				s := cast.ToInt(intVals[1].GetText())
				if s < 0 || s > precision {
					// logrus.Debugf("Scale '%d' is invalid for precision '%d' in column '%s', using 0 instead\n", s, precision, v.Colpath)
					s = 0
				} else if scale > s {
					// logrus.Debugf("Scale '%d' is larger than the defined scale '%d' for column '%s', using %d instead\n", scale, s, v.Colpath, s)
					scale = s
				}
			}

			var min, max int64
			if min_ == nil {
				min = -int64(math.Pow10(int(precision))) + 1 // Default min value
			} else {
				min = cast.ToInt64(min_)
			}
			if max_ == nil {
				max = int64(math.Pow10(int(precision))) - 1 // Default max value
			} else {
				max = cast.ToInt64(max_)
			}

			// TODO: Support larger precision
			intLen := precision - scale
			if intLen > MAX_DECIMAL_INT_LEN {
				intLen = MAX_DECIMAL_INT_LEN
			}

			g = NewFuncGen(func() any {
				var res [2]int64
				if intLen == 0 {
					res[0] = 0
				} else if min < 0 && rand.Float32() < 0.5 {
					delta := -float64(min)
					n := int64(math.Min(delta, math.Pow10(intLen)-1))
					res[0] = -rand.Int64N(n)
				} else {
					delta := float64(max) - math.Max(0, float64(min)) + 1
					lowerBound := int64(math.Max(0, float64(min)))
					n := int64(math.Min(delta, math.Pow10(intLen)-1))
					res[0] = lowerBound + rand.Int64N(n)
				}

				n := int64(math.Pow10(scale))
				if n <= 0 {
					res[1] = 0
				} else {
					res[1] = rand.Int64N(n)
				}

				return json.RawMessage(fmt.Sprintf("%d.%0*d", res[0], scale, res[1])) // Format as decimal string
			})
		case "DATE", "DATEV1", "DATEV2":
			min, max := CastMinMax[time.Time](min_, max_, baseType, v.Colpath)
			g = NewFuncGen(func() any { return gofakeit.DateRange(min, max).Format("2006-01-02") })
		case "DATETIME", "DATETIMEV1", "DATETIMEV2", "TIMESTAMP":
			min, max := CastMinMax[time.Time](min_, max_, baseType, v.Colpath)
			g = NewFuncGen(func() any { return gofakeit.DateRange(min, max).Format("2006-01-02 15:04:05") })
		case "TEXT", "STRING":
			lenMin, lenMax := v.GetLength()
			lenMin = lo.Max([]int{1, lenMin})
			lenMax = lo.Max([]int{1, lenMax})
			g = NewFuncGen(func() any { return RandomStr(lenMin, lenMax) })
		case "VARCHAR":
			var (
				length         int
				lenMin, lenMax = v.GetLength()
			)
			lenMin = lo.Max([]int{1, lenMin})
			lenMax = lo.Max([]int{1, lenMax})
			length_ := ty.INTEGER_VALUE(0)
			if length_ != nil {
				length = lo.Max([]int{1, cast.ToInt(length_.GetText())})
			} else {
				length = lenMax
			}
			if length < lenMax {
				lenMax = length
			}
			if lenMin > lenMax {
				lenMin = 1
			}
			g = NewFuncGen(func() any { return RandomStr(lenMin, lenMax) })
		case "CHAR":
			length_ := ty.INTEGER_VALUE(0)
			if length_ == nil {
				logrus.Fatalf("CHAR type must have a length in column '%s'\n", v.Colpath)
			}
			length := lo.Max([]int{1, cast.ToInt(length_.GetText())})
			if length > 255 {
				length = 255
			}
			g = NewFuncGen(func() any { return RandomStr(length, length) })
		case "IPV4":
			g = NewFuncGen(func() any { return gofakeit.IPv4Address() })
		case "IPV6":
			g = NewFuncGen(func() any { return gofakeit.IPv6Address() })
		default: // TODO: HLL, AGG_STATE, QUANTILE_STATE
			logrus.Fatalf("Unsupported column type '%s' for column '%s'\n", type_.GetText(), v.Colpath)
		}
	}
	if nullFrequency > 0 && nullFrequency <= 1 && baseType != "BITMAP" {
		return NewFuncGen(func() any {
			if rand.Float32() < nullFrequency {
				return nil
			}
			return g.Gen()
		})
	}
	return g
}

type Gen interface {
	Gen() any
}

type StructGen struct {
	Fields []*StructFieldGen
}

func (g *StructGen) AddChild(name string, child Gen) {
	g.Fields = append(g.Fields, &StructFieldGen{Name: name, Value: child})
}

func (g *StructGen) Gen() any {
	field2Data := lo.SliceToMap(g.Fields, func(field *StructFieldGen) (string, any) {
		return field.Name, field.Gen()
	})

	return json.RawMessage(MustJsonMarshal(field2Data))
}

type StructFieldGen struct {
	Name  string
	Value Gen
}

func (g *StructFieldGen) Gen() any {
	return g.Value.Gen()
}

type ArrayGen struct {
	Element        Gen
	LenMin, LenMax int
}

func (g *ArrayGen) Gen() any {
	len := rand.IntN(g.LenMax-g.LenMin+1) + g.LenMin

	elementData := lo.RepeatBy(len, func(_ int) any {
		return g.Element.Gen()
	})

	return json.RawMessage(MustJsonMarshal(elementData))
}

func (g *ArrayGen) SetElementGen(elem Gen) {
	g.Element = elem
}

type MapGen struct {
	Key, Value     Gen
	LenMin, LenMax int
}

func (g *MapGen) Gen() any {
	len := rand.IntN(g.LenMax-g.LenMin+1) + g.LenMin

	var b bytes.Buffer
	b.WriteByte('{')
	for i := 0; i < len; i++ {
		key := g.Key.Gen()
		value := g.Value.Gen()
		b.Write(MustJsonMarshal(key))
		b.WriteByte(':')
		b.Write(MustJsonMarshal(value))
		if i < len-1 {
			b.WriteByte(',')
		}
	}
	b.WriteByte('}')
	return json.RawMessage(b.Bytes())
}

func (g *MapGen) SetKeyGen(k Gen) {
	g.Key = k
}

func (g *MapGen) SetValueGen(v Gen) {
	g.Value = v
}

type fgen[T any] struct {
	f func() T
}

func NewIntGen[T int8 | int16 | int | int32](min, max T) Gen {
	return NewFuncGen(func() int { return gofakeit.IntRange(int(min), int(max)) })
}

func NewFuncGen[T any](f func() T) Gen {
	return &fgen[T]{f: f}
}

func (g *fgen[T]) Gen() any {
	return g.f()
}

func CastMinMax[R int8 | int16 | int | int32 | int64 | float32 | float64 | time.Time](min_, max_ any, baseType, colpath string, errmsg ...string) (R, R) {
	min, max, err := Cast2[R](min_, max_)
	if err != nil {
		msg := fmt.Sprintf("Invalid min/max %s '%v/%v' for column '%s': %v, expect %T", baseType, min_, max_, colpath, err, min)
		if len(errmsg) > 0 {
			msg += ", " + errmsg[0]
		}
		logrus.Fatalln(msg)
	}

	minBigger := false
	switch any(min).(type) {
	case int8:
		minBigger = any(max).(int8) < any(min).(int8)
	case int16:
		minBigger = any(max).(int16) < any(min).(int16)
	case int:
		minBigger = any(max).(int) < any(min).(int)
	case int32:
		minBigger = any(max).(int32) < any(min).(int32)
	case int64:
		minBigger = any(max).(int64) < any(min).(int64)
	case float32:
		minBigger = any(max).(float32) < any(min).(float32)
	case float64:
		minBigger = any(max).(float64) < any(min).(float64)
	case time.Time:
		minBigger = any(max).(time.Time).Before(any(min).(time.Time))
	}
	if minBigger {
		logrus.Warnf("Column '%s' max(%v) < min(%v), set max to min\n", colpath, max, min)
		max = min
	}
	return min, max
}

func CheckGenRowCount(rows int) error {
	if rows < 0 {
		return fmt.Errorf("--rows/row_count must be a positive integer, got %d", rows)
	} else if rows > MaxGenRowCount {
		return fmt.Errorf("--rows/row_count must be smaller than 100_000, got %d", rows)
	}
	return nil
}
