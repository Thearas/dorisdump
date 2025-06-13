package src

import (
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	GLOBAL_NULL_FREQUENCY = 0.0 // Default null frequency is 0%
	MAX_DECIMAL_INT_LEN   = len(strconv.FormatInt(math.MaxInt64, 10))
	GlobalGenRule         = GenRule{
		"null_frequency": GLOBAL_NULL_FREQUENCY,
	}
	DefaultTypeGenRules = map[string]GenRule{
		"ARRAY": {
			"length": GenRule{
				"min": 1,
				"max": 3,
			},
		},
		"MAP": {
			"length": GenRule{
				"min": 1,
				"max": 3,
			},
		},
		"JSON": {
			"structure": `STRUCT<col1:SMALLINT, col2:SMALLINT>`, // Default structure
		},
		"JSONB": {
			"structure": `STRUCT<col1:SMALLINT, col2:SMALLINT>`, // Default structure
		},
		"VARIANT": {
			"structure": `STRUCT<col1:SMALLINT, col2:SMALLINT>`, // Default structure
		},
		"BITMAP": {
			"length": 5,
			"min":    0,
			"max":    math.MaxInt32,
		},
		"TEXT": {
			"length": GenRule{
				"min": 1,
				"max": 10,
			},
		},
		"STRING": {
			"length": GenRule{
				"min": 1,
				"max": 10,
			},
		},
		"VARCHAR": {
			"length": GenRule{
				"min": 1,
				"max": 10,
			},
		},
		"TINYINT": {
			"min": math.MinInt8,
			"max": math.MaxInt8,
		},
		"SMALLINT": {
			"min": math.MinInt16,
			"max": math.MaxInt16,
		},
		"INT": {
			"min": math.MinInt32,
			"max": math.MaxInt32,
		},
		"BIGINT": {
			"min": math.MinInt32,
			"max": math.MaxInt32,
		},
		"LARGEINT": {
			"min": math.MinInt32,
			"max": math.MaxInt32,
		},
		"FLOAT": {
			"min": math.MinInt16,
			"max": math.MaxInt16,
		},
		"DOUBLE": {
			"min": math.MinInt32,
			"max": math.MaxInt32,
		},
		"DECIMAL": {
			"min": math.MinInt32,
			"max": math.MaxInt32,
		},
		"DATE": {
			"min": time.Now().AddDate(-10, 0, 0),
			"max": time.Now(),
		},
		"DATETIME": {
			"min": time.Now().AddDate(-10, 0, 0),
			"max": time.Now(),
		},
	}
)

func SetupDefaultGenRules(configFile string) error {
	if configFile != "" {
		b, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(b, &GlobalGenRule); err != nil {
			return err
		}
	}
	if g, ok := GlobalGenRule["type"]; !ok || g == nil {
		GlobalGenRule["type"] = GenRule{}
	}
	GlobalGenRule["type"] = lo.MapEntries(GlobalGenRule["type"].(GenRule), func(ty string, g any) (string, GenRule) {
		if g == nil {
			g = GenRule{}
		}
		genRule, ok := g.(GenRule)
		if !ok {
			logrus.Fatalf("Type gen rule for '%s' should be a map, but got '%T'\n", ty, g)
		}
		return strings.ToUpper(ty), genRule
	})
	if err := mergo.Merge(&DefaultTypeGenRules, GlobalGenRule["type"], mergo.WithOverride); err != nil {
		logrus.Fatalln("Merge global and default gen rules failed, err:", err)
	}

	// copy null_frequency to every types' gen rule
	for _, genRule := range DefaultTypeGenRules {
		if r, ok := genRule["null_frequency"]; !ok || r == nil {
			genRule["null_frequency"] = GlobalGenRule["null_frequency"]
		}
	}

	return nil
}

func getCustomTableGenRule(table string) (rows int, colrules map[string]GenRule) {
	tableParts := strings.Split(table, ".")
	tablePart := tableParts[len(tableParts)-1]

	g, ok := GlobalGenRule["tables"].([]any)
	if !ok || len(g) == 0 {
		logrus.Debugf("no custom gen rule for table '%s'\n", table)
		return 0, map[string]GenRule{}
	}

	tg_, found := lo.Find(g, func(tg_ any) bool {
		tg, ok := tg_.(GenRule)
		if !ok {
			logrus.Fatalf("custom table gen rule for '%s' should be a map", table)
		}
		return tg["name"] == tablePart
	})
	if !found {
		logrus.Debugf("no custom gen rule for table '%s'\n", table)
		return 0, map[string]GenRule{}
	}
	tg := tg_.(GenRule)

	// get table row_count
	rowCount, ok := tg["row_count"].(int)
	if ok {
		rows = rowCount
	}

	// get table columns gen rule
	cgs, ok := tg["columns"].([]any)
	if !ok || len(cgs) == 0 {
		logrus.Debugf("no custom gen rule for table columns '%s'\n", table)
		return 0, map[string]GenRule{}
	}

	i := 0
	colrules = lo.SliceToMap(cgs, func(cg_ any) (string, GenRule) {
		cg, ok := cg_.(GenRule)
		if !ok {
			logrus.Fatalf("custom column gen rule for '%s.#%d' should be a map", table, i)
		}

		name, ok := cg["name"].(string)
		if !ok {
			logrus.Fatalf("Column field #%d has no name in table '%s'\n", i, table)
		}
		i++
		return name, cg
	})
	return
}
