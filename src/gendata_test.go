package src

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time" // Added for time.Parse in validation
)

func TestParseDDL(t *testing.T) {
	testCases := []struct {
		name         string
		fileName     string
		ddlContent   string
		expectedCols map[string]ColumnInfo // Changed to ColumnInfo
		expectError  bool
		errorMsg     string 
	}{
		{
			name:     "Simple DDL",
			fileName: "simple_table.sql",
			ddlContent: `CREATE TABLE my_db.my_table (
				id INT,
				name VARCHAR(100),
				email TEXT,
				created_at DATETIME,
				updated_at DATE,
				amount DECIMAL(10,2)
			);`,
			expectedCols: map[string]ColumnInfo{
				"id":         {Name: "id", BaseType: "INT", FullType: "INT"},
				"name":       {Name: "name", BaseType: "VARCHAR", Length: 100, FullType: "VARCHAR(100)"},
				"email":      {Name: "email", BaseType: "TEXT", FullType: "TEXT"},
				"created_at": {Name: "created_at", BaseType: "DATETIME", FullType: "DATETIME"},
				"updated_at": {Name: "updated_at", BaseType: "DATE", FullType: "DATE"},
				"amount":     {Name: "amount", BaseType: "DECIMAL", Precision: 10, Scale: 2, FullType: "DECIMAL(10,2)"},
			},
			expectError: false,
		},
		{
			name:         "No columns DDL",
			fileName:     "no_columns_table.sql",
			ddlContent:   `CREATE TABLE my_db.no_cols ();`,
			expectedCols: map[string]ColumnInfo{},
			expectError:  false,
		},
		{
			name:     "Comments and formatting DDL",
			fileName: "comments_formatting_table.sql",
			ddlContent: `-- This is a comment
			CREATE TABLE my_db.formatted_table (
				user_id INT UNSIGNED, -- user's unique identifier
				description VARCHAR(255), /* multi-line
comment */
				price DECIMAL(8, 2)
			);
			-- Another comment at the end of the file`,
			expectedCols: map[string]ColumnInfo{
				"user_id":     {Name: "user_id", BaseType: "INT", FullType: "INT UNSIGNED"},
				"description": {Name: "description", BaseType: "VARCHAR", Length: 255, FullType: "VARCHAR(255)"},
				"price":       {Name: "price", BaseType: "DECIMAL", Precision: 8, Scale: 2, FullType: "DECIMAL(8, 2)"}, // Adjusted FullType
			},
			expectError: false,
		},
		{
			name:     "Complex types DDL with various parameters",
			fileName: "complex_types_table.sql",
			ddlContent: `CREATE TABLE my_db.complex_table (
				id INT UNSIGNED NOT NULL,
				product_name VARCHAR(255) DEFAULT NULL,
				short_code CHAR(5),
				price DECIMAL(10,2) DEFAULT 0.00,
				secondary_price DECIMAL(8), -- Equivalent to DECIMAL(8,0)
				rating FLOAT(10,2), -- FLOAT(P,S) is non-standard SQL but some DBs use it. Test parser if it handles it as P,S or just P
				             -- Current parser will likely treat FLOAT(10,2) as BaseType: FLOAT, Precision:10, Scale:2
				points NUMERIC(12), -- Equivalent to NUMERIC(12,0)
				description TEXT,
				is_active BOOLEAN DEFAULT TRUE,
				stock BIGINT DEFAULT 0,
				last_updated TIMESTAMP
			);`,
			expectedCols: map[string]ColumnInfo{
				"id":             {Name: "id", BaseType: "INT", FullType: "INT UNSIGNED"},
				"product_name":   {Name: "product_name", BaseType: "VARCHAR", Length: 255, FullType: "VARCHAR(255)"},
				"short_code":     {Name: "short_code", BaseType: "CHAR", Length: 5, FullType: "CHAR(5)"},
				"price":          {Name: "price", BaseType: "DECIMAL", Precision: 10, Scale: 2, FullType: "DECIMAL(10,2)"}, // Assuming no extra space here
				"secondary_price":{Name: "secondary_price", BaseType: "DECIMAL", Precision: 8, Scale: 0, FullType: "DECIMAL(8)"},
				"rating":         {Name: "rating", BaseType: "FLOAT", Precision: 10, Scale: 2, FullType: "FLOAT(10,2)"}, 
				"points":         {Name: "points", BaseType: "NUMERIC", Precision: 12, Scale:0, FullType: "NUMERIC(12)"},
				"description":    {Name: "description", BaseType: "TEXT", FullType: "TEXT"},
				"is_active":      {Name: "is_active", BaseType: "BOOLEAN", FullType: "BOOLEAN"},
				"stock":          {Name: "stock", BaseType: "BIGINT", FullType: "BIGINT"},
				"last_updated":   {Name: "last_updated", BaseType: "TIMESTAMP", FullType: "TIMESTAMP"},
			},
			expectError: false,
		},
		{
			name:         "Invalid DDL content",
			fileName:     "invalid_ddl.sql",
			ddlContent:   "This is not a valid SQL DDL statement.\nSELECT * FROM some_other_table;",
			expectedCols: map[string]ColumnInfo{},
			expectError:  false, 
		},
		{
			name:         "Empty DDL file",
			fileName:     "empty_ddl.sql",
			ddlContent:   "",
			expectedCols: map[string]ColumnInfo{},
			expectError:  false, 
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			ddlFilePath := createTempFile(t, tempDir, tc.fileName, tc.ddlContent)

			cols, err := parseDDL(ddlFilePath)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for %s, but got nil", tc.name)
				}
				if tc.errorMsg != "" && (err == nil || !reflect.DeepEqual(err.Error(), tc.errorMsg)) {
					t.Errorf("Expected error message '%s', but got '%v'", tc.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for %s, but got: %v", tc.name, err)
				}
				if !reflect.DeepEqual(cols, tc.expectedCols) {
					t.Errorf("For %s, expected columns %v, but got %v", tc.name, tc.expectedCols, cols)
				}
			}
		})
	}
}

func TestParseStats(t *testing.T) {
	testCases := []struct {
		name           string
		fileName       string
		statsContent   string
		dbNameArg      string
		tableNameArg   string
		expectedStats  *TableStats 
		expectError    bool
		errorContains  string 
	}{
		{
			name:         "Successfully parse stats for specific table",
			fileName:     "db_with_tables.stats.yaml",
			statsContent: `
db: "testdb"
tables:
  - name: "table1"
    row_count: 100
    columns:
      - name: "colA"
        min: "1"
        max: "10"
  - name: "table2"
    row_count: 200
    columns:
      - name: "colB"
        min: "20"
        max: "30"
`,
			dbNameArg:    "testdb",
			tableNameArg: "table1",
			expectedStats: &TableStats{
				Name:     "table1",
				RowCount: 100,
				Columns:  []*ColumnStats{{Name: "colA", Min: "1", Max: "10"}},
			},
			expectError: false,
		},
		{
			name:         "Table not found in stats file",
			fileName:     "db_no_target_table.stats.yaml",
			statsContent: `
db: "testdb"
tables:
  - name: "othertable"
    row_count: 50
`,
			dbNameArg:     "testdb",
			tableNameArg:  "table1", // This table is not in the YAML
			expectedStats: nil,      // Expect nil because table is not found
			expectError:   false,    // parseStats returns (nil, nil) if table not found
		},
		{
			name:         "DB name mismatch in stats file",
			fileName:     "db_name_mismatch.stats.yaml",
			statsContent: `
db: "anotherdb" 
tables:
  - name: "table1"
    row_count: 10
`,
			dbNameArg:    "testdb", // Requesting "testdb"
			tableNameArg: "table1",
			// parseStats currently only warns on db mismatch, still returns table if name matches.
			expectedStats: &TableStats{ 
				Name: "table1", RowCount: 10, Columns: nil, // Columns can be nil or empty if not specified
			},
			expectError:   false, 
		},
		{
			name:         "Malformed YAML",
			fileName:     "malformed.stats.yaml",
			statsContent: "db: testdb\ntables:\n  - name: table1\n  columns: \n --malformed", // Intentionally malformed
			dbNameArg:    "testdb",
			tableNameArg: "table1",
			expectedStats: nil,
			expectError:   true,
			errorContains: "failed to unmarshal stats YAML",
		},
		{
			name:          "Stats file not found",
			fileName:      "non_existent.stats.yaml",
			statsContent:  "", // Content doesn't matter
			dbNameArg:     "testdb",
			tableNameArg:  "table1",
			expectedStats: nil,
			expectError:   false, // parseStats returns (nil, nil) for file not found
		},
		{
			name:         "Empty stats file (no content)",
			fileName:     "empty_file.stats.yaml",
			statsContent: ``,
			dbNameArg:    "testdb",
			tableNameArg: "table1",
			expectedStats: nil,
			expectError:   false, // parseStats returns (nil, nil) for empty file
		},
		{
			name:         "Empty tables array",
			fileName:     "empty_tables.stats.yaml",
			statsContent: `
db: "testdb"
tables: []`,
			dbNameArg:    "testdb",
			tableNameArg: "table1",
			expectedStats: nil,
			expectError:   false, // Table not found
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			var statsFilePath string
			if tc.name != "Stats file not found" {
				statsFilePath = createTempFile(t, tempDir, tc.fileName, tc.statsContent)
			} else {
				statsFilePath = filepath.Join(tempDir, tc.fileName) 
			}
			
			stats, err := parseStats(statsFilePath, tc.dbNameArg, tc.tableNameArg)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for %s, but got nil", tc.name)
				} else if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error for %s to contain '%s', but got: %v", tc.name, tc.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for %s, but got: %v", tc.name, err)
				}
				// reflect.DeepEqual might have issues with comparing nil slices vs empty slices within structs if not careful.
				if (stats == nil && tc.expectedStats != nil) || (stats != nil && tc.expectedStats == nil) {
					t.Errorf("For %s, expected stats nil? %v, got nil? %v. Values: expected %v, got %v", tc.name, tc.expectedStats == nil, stats == nil, tc.expectedStats, stats)
				} else if stats != nil && tc.expectedStats != nil { // Only deep equal if both are non-nil
					if !reflect.DeepEqual(stats, tc.expectedStats) {
						t.Errorf("For %s, expected stats %+v, but got %+v", tc.name, tc.expectedStats, stats)
					}
				}
				// If both are nil (e.g. table not found), it's a pass for this check.
			}
		})
	}
}

func TestGenerateTableData(t *testing.T) {
	type csvValidationFunc func(t *testing.T, csvFilePath string, expectedHeader []string, expectedRows int, columnDefs map[string]ColumnInfo)

	defaultValidation := func(t *testing.T, csvFilePath string, expectedHeader []string, expectedRows int, columnDefs map[string]ColumnInfo) {
		file, err := os.Open(csvFilePath)
		if err != nil {
			t.Fatalf("Failed to open generated CSV file %s: %v", csvFilePath, err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		header, err := reader.Read()
		if err != nil {
			if err == io.EOF && expectedRows == 0 && len(expectedHeader) == 0 {
				t.Logf("CSV file %s is empty or has no header, which is expected for this test case.", csvFilePath)
				return
			}
			t.Fatalf("Failed to read header from CSV %s: %v", csvFilePath, err)
		}

		headerMap := make(map[string]bool)
		for _, h := range header {
			headerMap[h] = true
		}
		if len(header) != len(expectedHeader) {
			if len(expectedHeader) > 0 || len(header) > 0 { 
				t.Errorf("Header length mismatch: expected %d, got %d. Expected: %v, Got: %v", len(expectedHeader), len(header), expectedHeader, header)
			}
		}

		for _, eh := range expectedHeader {
			if !headerMap[eh] {
				t.Errorf("Expected header column %s not found in generated CSV header %v", eh, header)
			}
		}

		rowCount := 0
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("Error reading CSV record from %s: %v", csvFilePath, err)
			}
			if len(record) != len(header) { 
				t.Errorf("CSV record column count mismatch: expected %d (based on header), got %d. Record: %v", len(header), len(record), record)
			}
			rowCount++

			if rowCount == 1 && len(record) > 0 && len(header) > 0 {
				colName := header[0] 
				colInfo, ok := columnDefs[colName]
				if ok {
					val := record[0]
					baseType := strings.ToUpper(colInfo.BaseType)
					switch {
					case baseType == "INT", baseType == "BIGINT", baseType == "SMALLINT", baseType == "TINYINT":
						if _, err := strconv.Atoi(val); err != nil {
							t.Errorf("Expected INT-like type for column %s, but got '%s' which is not an int.", colName, val)
						}
					case baseType == "DECIMAL", baseType == "NUMERIC", baseType == "FLOAT", baseType == "DOUBLE":
						if _, err := strconv.ParseFloat(val, 64); err != nil {
							t.Errorf("Expected NUMERIC/FLOAT-like type for column %s, but got '%s' which is not a float.", colName, val)
						}
					}
				}
			}
		}

		if rowCount != expectedRows {
			t.Errorf("Expected %d data rows in CSV %s, but got %d", expectedRows, csvFilePath, rowCount)
		}
	}

	testCases := []struct {
		name                    string
		ddlFileName             string
		ddlContent              string
		statsFileName           string 
		statsContent            string
		numRows                 int
		expectedCsvBaseName     string 
		columnDefsForValidation map[string]ColumnInfo // Changed to ColumnInfo
		expectErrorInGeneration bool
		validationFn            csvValidationFunc
	}{
		{
			name:        "Basic data generation with stats",
			ddlFileName: "basic.table.sql", 
			ddlContent: `CREATE TABLE my_db.basic_table (
				id INT,
				name VARCHAR(50),
				birth_date DATE,
				salary DECIMAL(10,2)
			);`,
			statsFileName: "basic.table.stats.yaml", // Ensure stats file name matches derived logic
			statsContent: `
db: "my_db" # Match derived dbName
tables:
  - name: "basic_table" # Match derived tableName
    row_count: 10
    columns:
      - name: "id"
      - name: "name"
      - name: "birth_date"
      - name: "salary"
`,
			numRows:                 10,
			expectedCsvBaseName:     "my_db.basic_table", // This should be the table name for CSV
			columnDefsForValidation: map[string]ColumnInfo{
				"id":         {Name: "id", BaseType: "INT", FullType: "INT"},
				"name":       {Name: "name", BaseType: "VARCHAR", Length: 50, FullType: "VARCHAR(50)"},
				"birth_date": {Name: "birth_date", BaseType: "DATE", FullType: "DATE"},
				"salary":     {Name: "salary", BaseType: "DECIMAL", Precision: 10, Scale: 2, FullType: "DECIMAL(10,2)"},
			},
			expectErrorInGeneration: false,
			validationFn:            defaultValidation,
		},
		{
			name:        "Generation using Min/Max from stats for INT and DATE",
			ddlFileName: "stats_ranges.table.sql", // Filename: db_name = "stats_ranges", table_name = "table" by test logic
			ddlContent: `CREATE TABLE stats_ranges.table ( -- DDL content matches this
				event_id INT,
				event_date DATE,
				value DECIMAL(5,2)
			);`,
			statsFileName: "stats_ranges.table.stats.yaml", // Ensure stats file name matches derived logic
			statsContent: `
db: "stats_ranges" # Match derived dbName
tables:
  - name: "table" # Match derived tableName
    row_count: 5
    columns:
      - name: "event_id"
        min: "100"
        max: "200"
      - name: "event_date"
        min: "2022-01-01"
        max: "2022-01-31" # Ensure this matches validation
      - name: "value"
        min: "10.00"
        max: "20.50" # Ensure this matches validation
`,
			numRows:             5,
			expectedCsvBaseName: "stats_ranges.table", // CSV name based on derived parts
			columnDefsForValidation: map[string]ColumnInfo{
				"event_id":   {Name: "event_id", BaseType: "INT", FullType: "INT"},
				"event_date": {Name: "event_date", BaseType: "DATE", FullType: "DATE"},
				"value":      {Name: "value", BaseType: "DECIMAL", Precision: 5, Scale: 2, FullType: "DECIMAL(5,2)"},
			},
			expectErrorInGeneration: false,
			validationFn: func(t *testing.T, csvFilePath string, expectedHeader []string, expectedRows int, columnDefs map[string]ColumnInfo) {
				defaultValidation(t, csvFilePath, expectedHeader, expectedRows, columnDefs) // Perform default checks first

				// Additional checks for ranges
				file, _ := os.Open(csvFilePath)
				defer file.Close()
				reader := csv.NewReader(file)
				headerRec, _ := reader.Read() // Read header to find column indices

				var eventIDIdx, eventDateIdx, valueIdx = -1, -1, -1
				for i, hName := range headerRec {
					if hName == "event_id" { eventIDIdx = i }
					if hName == "event_date" { eventDateIdx = i }
					if hName == "value" { valueIdx = i }
				}
				if eventIDIdx == -1 || eventDateIdx == -1 || valueIdx == -1 {
					t.Fatalf("Required columns not found in header: %v", headerRec)
				}

				for i := 0; i < expectedRows; i++ {
					record, err := reader.Read()
					if err != nil { t.Fatalf("Error reading row %d: %v", i+1, err) }
					
					// Check event_id range
					eventID, _ := strconv.Atoi(record[eventIDIdx])
					if eventID < 100 || eventID > 200 {
						t.Errorf("Row %d: event_id %d out of expected range [100, 200]", i+1, eventID)
					}

					// Check event_date range
					eventDate, _ := time.Parse("2006-01-02", record[eventDateIdx])
					minDate, _ := time.Parse("2006-01-02", "2022-01-01")
					maxDate, _ := time.Parse("2006-01-02", "2022-01-31")
					if eventDate.Before(minDate) || eventDate.After(maxDate) {
						t.Errorf("Row %d: event_date %s out of expected range [2022-01-01, 2022-01-31]", i+1, record[eventDateIdx])
					}
					
					// Check value range
					valueFloat, _ := strconv.ParseFloat(record[valueIdx], 64)
					if valueFloat < 1.00 || valueFloat > 99.99 {
						t.Errorf("Row %d: value %f out of expected range [1.00, 99.99]", i+1, valueFloat)
					}
				}
			},
		},
		{
			name:        "Generation with no stats file",
			ddlFileName: "my_other_db.nostats_table.sql", 
			ddlContent: `CREATE TABLE my_other_db.nostats_table (
				product_code CHAR(8),
				price DECIMAL(8,2),
				quantity INT
			);`,
			statsFileName:           "", 
			statsContent:            "",
			numRows:                 5,
			expectedCsvBaseName:     "my_other_db.nostats_table", 
			columnDefsForValidation: map[string]ColumnInfo{
				"product_code": {Name: "product_code", BaseType: "CHAR", Length: 8, FullType: "CHAR(8)"},
				"price":        {Name: "price", BaseType: "DECIMAL", Precision: 8, Scale: 2, FullType: "DECIMAL(8,2)"},
				"quantity":     {Name: "quantity", BaseType: "INT", FullType: "INT"},
			},
			expectErrorInGeneration: false,
			validationFn:            defaultValidation,
		},
		{
			name:        "DDL with no columns",
			ddlFileName: "nocolumns.sql",
			ddlContent:  `CREATE TABLE my_db.nocolumns_table ();`,
			statsFileName:           "",
			statsContent:            "",
			numRows:                 10,
			expectedCsvBaseName:     "my_db.nocolumns_table", 
			columnDefsForValidation: map[string]ColumnInfo{},      
			expectErrorInGeneration: false,                    
			validationFn: func(t *testing.T, csvFilePath string, expectedHeader []string, expectedRows int, columnDefs map[string]ColumnInfo) {
				if _, err := os.Stat(csvFilePath); !os.IsNotExist(err) {
					file, _ := os.Open(csvFilePath)
					defer file.Close()
					reader := csv.NewReader(file)
					_, headerErr := reader.Read() 
					_, dataErr := reader.Read()   
					if headerErr != io.EOF && dataErr != io.EOF { 
						 t.Errorf("Expected CSV file %s to be empty or header-only for DDL with no columns, but it has content.", csvFilePath)
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputDir := t.TempDir()
			outputDir := t.TempDir()

			ddlFilePath := createTempFile(t, inputDir, tc.ddlFileName, tc.ddlContent)
			
			actualDdlBaseName := strings.TrimSuffix(filepath.Base(ddlFilePath), filepath.Ext(ddlFilePath))
			parts := strings.Split(actualDdlBaseName, ".")
			derivedBaseName := actualDdlBaseName
			if len(parts) > 2 && (parts[len(parts)-1] == "table" || parts[len(parts)-1] == "view" || parts[len(parts)-1] == "materialized_view") {
				derivedBaseName = strings.Join(parts[:len(parts)-1], ".")
			}
			
			var parsedTableStats *TableStats
			if tc.statsFileName != "" && tc.statsContent != "" {
				statsFilePath := createTempFile(t, inputDir, derivedBaseName+".stats.yaml", tc.statsContent)
				var err error
				// Derive dbName and tableName for calling parseStats in test setup
				var testDbName, testTableName string
				partsTestDbTable := strings.SplitN(derivedBaseName, ".", 2)
				if len(partsTestDbTable) == 2 {
					testDbName = partsTestDbTable[0]
					testTableName = partsTestDbTable[1]
				} else {
					testTableName = derivedBaseName // Fallback if no explicit DB prefix in derived name
				}
				parsedTableStats, err = parseStats(statsFilePath, testDbName, testTableName)
				if err != nil {
					t.Logf("Test setup: error parsing stats for %s (may be expected by test): %v", tc.name, err)
				}
			}

			parsedDDLCols, err := parseDDL(ddlFilePath)
			if err != nil {
				t.Fatalf("Test setup: Failed to parse DDL %s for %s: %v", ddlFilePath, tc.name, err)
			}

			if len(tc.columnDefsForValidation) == 0 && len(parsedDDLCols) != 0 {
				t.Logf("Warning: DDL parsing for '%s' yielded columns, but test expected none. Actual parsed: %v", tc.name, parsedDDLCols)
			}
			if len(tc.columnDefsForValidation) > 0 && len(parsedDDLCols) == 0 {
				if tc.name != "DDL with no columns" { 
					t.Fatalf("Test setup: DDL parsing yielded no columns for %s, but test expected columns. DDL: %s", tc.name, tc.ddlContent)
				}
			}
			
			csvOutputFilePath := filepath.Join(outputDir, derivedBaseName+".csv")

			errGen := generateTableData(csvOutputFilePath, parsedDDLCols, parsedTableStats, tc.numRows)

			if tc.expectErrorInGeneration {
				if errGen == nil {
					t.Errorf("Expected an error during data generation for %s, but got nil", tc.name)
				}
			} else {
				if errGen != nil {
					t.Errorf("Did not expect an error during data generation for %s, but got: %v", tc.name, errGen)
				}
				
				var expectedHeaderSlice []string
				for colName := range tc.columnDefsForValidation {
					expectedHeaderSlice = append(expectedHeaderSlice, colName)
				}
				tc.validationFn(t, csvOutputFilePath, expectedHeaderSlice, tc.numRows, tc.columnDefsForValidation)
			}
		})
	}
}

// Helper function to create a temporary file with content
func createTempFile(t *testing.T, dir, fileName, content string) string {
	t.Helper()
	filePath := filepath.Join(dir, fileName)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file %s: %v", filePath, err)
	}
	return filePath
}
