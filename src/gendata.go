/*
Copyright © 2025 Thearas thearas850@gmail.com

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
package src

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv" // Ensured strconv is imported once
	"strings"
	// "time" // Removed as it's unused

	"time" // Added for date parsing

	"github.com/brianvoe/gofakeit/v6" 
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ColumnInfo holds detailed information about a parsed DDL column.
type ColumnInfo struct {
	Name      string
	BaseType  string // e.g., VARCHAR, DECIMAL, INT
	Length    int    // For VARCHAR(N), CHAR(N)
	Precision int    // For DECIMAL(P,S) -> P
	Scale     int    // For DECIMAL(P,S) -> S
	FullType  string // The originally parsed full type string
}

// Config mirrors the cmd.GendataConfig struct to avoid circular dependencies.
type Config struct {
	InputDDLDir  string
	OutputCsvDir string
	NumRows      int
}

// GenerateData is the main function for the data generation logic.
// It takes Config to understand where to read DDLs from, where to write CSVs,
// and how many rows to generate.
func GenerateData(config Config) error {
	logrus.Infof("Starting data generation with config: %+v", config)

	ddlPattern := filepath.Join(config.InputDDLDir, "*.sql")
	ddlFiles, err := filepath.Glob(ddlPattern)
	if err != nil {
		logrus.Errorf("Error listing DDL files using pattern '%s': %v", ddlPattern, err)
		return err
	}

	if len(ddlFiles) == 0 {
		logrus.Warnf("No DDL files found in '%s'", config.InputDDLDir)
		return nil
	}

	logrus.Infof("Found %d DDL files to process.", len(ddlFiles))

	for _, ddlFilePath := range ddlFiles {
		baseName := strings.TrimSuffix(filepath.Base(ddlFilePath), filepath.Ext(ddlFilePath))
		parts := strings.Split(baseName, ".")
		if len(parts) < 2 {
			logrus.Warnf("Skipping DDL file with unexpected name format: %s", ddlFilePath)
			continue
		}
		tableNamePrefix := baseName // Default for names like <db>.<table>.sql
		if len(parts) > 2 && (parts[len(parts)-1] == "table" || parts[len(parts)-1] == "view" || parts[len(parts)-1] == "materialized_view") {
			tableNamePrefix = strings.Join(parts[:len(parts)-1], ".") // <db>.<table>
		}
		
		statsFileName := tableNamePrefix + ".stats.yaml"
		statsFilePath := filepath.Join(config.InputDDLDir, statsFileName)

		csvFileName := tableNamePrefix + ".csv"
		csvFilePath := filepath.Join(config.OutputCsvDir, csvFileName)

		logrus.Infof("Processing DDL: %s", ddlFilePath)
		logrus.Infof("  Stats file: %s", statsFilePath)
		logrus.Infof("  Output CSV: %s", csvFilePath)

		columns, err := parseDDL(ddlFilePath)
		if err != nil {
			logrus.Errorf("Error parsing DDL file %s: %v. Skipping table.", ddlFilePath, err)
			continue
		}
		if len(columns) == 0 {
			logrus.Warnf("No columns parsed from DDL %s. Skipping table.", ddlFilePath)
			continue
		}
		
		// Derive dbName and tableName from tableNamePrefix for parseStats
		var dbName, actualTableName string
		partsDbTable := strings.SplitN(tableNamePrefix, ".", 2)
		if len(partsDbTable) == 2 {
			dbName = partsDbTable[0]
			actualTableName = partsDbTable[1]
		} else {
			actualTableName = tableNamePrefix
			logrus.Warnf("Could not derive database name from DDL path '%s' for stats lookup. Using empty dbName for stats lookup.", ddlFilePath)
		}

		tableStats, err := parseStats(statsFilePath, dbName, actualTableName) // Corrected call
		if err != nil {
			// parseStats now returns error for file issues or YAML unmarshal issues.
			// It returns (nil, nil) if file not found, empty, or table/db not matched (these are logged as warnings by parseStats).
			logrus.Errorf("Error encountered while trying to parse stats file %s: %v. Proceeding without stats.", statsFilePath, err)
			// tableStats will be nil if an error occurred or if no specific stats were found, which is handled by generateTableData
		}

		err = generateTableData(csvFilePath, columns, tableStats, config.NumRows)
		if err != nil {
			logrus.Errorf("Error generating data for table from DDL %s: %v. Skipping table.", ddlFilePath, err)
			continue
		}
	}

	logrus.Info("Data generation process completed.")
	return nil
}

func parseDDL(ddlFilePath string) (map[string]ColumnInfo, error) { 
	logrus.Debugf("Parsing DDL file: %s", ddlFilePath)
	content, err := os.ReadFile(ddlFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read DDL file %s: %w", ddlFilePath, err)
	}

	columns := make(map[string]ColumnInfo) 
	colDefinitionRegex := regexp.MustCompile(
		"`?([a-zA-Z0-9_]+)`?" + 
			"\\s+" +
			"([a-zA-Z]+(?:\\s*\\(\\s*\\d+\\s*(?:,\\s*\\d+\\s*)?\\))?(?:\\s+UNSIGNED|\\s+ZEROFILL)?)",
	)

	typeWithLengthRegex := regexp.MustCompile(`([a-zA-Z]+)\s*\((\d+)\)`)                 
	typeWithPrecisionScaleRegex := regexp.MustCompile(`([a-zA-Z]+)\s*\((\d+)\s*,\s*(\d+)\)`) 

	createTableRegex := regexp.MustCompile(`(?is)CREATE(?:\s+TEMPORARY)?\s+TABLE(?:\s+IF\s+NOT\s+EXISTS)?\s*` + "`?" + `[^` + "`" + `.\s]+` + "`?" + `\.` + "`?" + `[^` + "`" + `.\s]+` + "`?" + `\s*\((.*)\)`)
	matches := createTableRegex.FindSubmatch(content)
	if len(matches) < 2 {
		createTableRegex = regexp.MustCompile(`(?is)CREATE(?:\s+TEMPORARY)?\s+TABLE(?:\s+IF\s+NOT\s+EXISTS)?\s*` + "`?" + `[^` + "`" + `.\s]+` + "`?" + `\s*\((.*)\)`)
		matches = createTableRegex.FindSubmatch(content)
		if len(matches) < 2 {
			logrus.Warnf("Could not find CREATE TABLE statement or column definitions in %s", ddlFilePath)
			return columns, nil
		}
	}
	columnDefsBlock := string(matches[1])

	blockCommentRegex := regexp.MustCompile(`/\*(?s:.*?)\*/`)
	columnDefsBlock = blockCommentRegex.ReplaceAllString(columnDefsBlock, "")
	lineCommentRegex := regexp.MustCompile(`--.*?(\n|$)`)
	columnDefsBlock = lineCommentRegex.ReplaceAllString(columnDefsBlock, "\n")

	defs := strings.Split(columnDefsBlock, "\n")

	for _, line := range defs {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") ||
			strings.HasPrefix(strings.ToUpper(line), "KEY") || strings.HasPrefix(strings.ToUpper(line), "INDEX") ||
			strings.HasPrefix(strings.ToUpper(line), "CONSTRAINT") || strings.HasPrefix(strings.ToUpper(line), "UNIQUE KEY") ||
			strings.HasPrefix(line, ")") {
			continue
		}

		match := colDefinitionRegex.FindStringSubmatch(line)
		if len(match) < 3 {
			logrus.Debugf("Line did not match column definition regex: %s", line)
			continue 
		}

		colName := match[1]
		fullTypeStr := strings.ToUpper(match[2]) 

		colInfo := ColumnInfo{
			Name:      colName,
			FullType:  fullTypeStr,
			Length:    0, 
			Precision: 0, 
			Scale:     0, 
		}

		if psMatch := typeWithPrecisionScaleRegex.FindStringSubmatch(fullTypeStr); len(psMatch) == 4 {
			colInfo.BaseType = psMatch[1]
			colInfo.Precision, _ = strconv.Atoi(psMatch[2])
			colInfo.Scale, _ = strconv.Atoi(psMatch[3])
		} else if lMatch := typeWithLengthRegex.FindStringSubmatch(fullTypeStr); len(lMatch) == 3 {
			baseTypeName := lMatch[1]
			numVal, _ := strconv.Atoi(lMatch[2])
			if baseTypeName == "DECIMAL" || baseTypeName == "NUMERIC" || baseTypeName == "FLOAT" || baseTypeName == "DOUBLE" { 
				colInfo.BaseType = baseTypeName
				colInfo.Precision = numVal
				colInfo.Scale = 0 
			} else { 
				colInfo.BaseType = baseTypeName
				colInfo.Length = numVal
			}
		} else {
			typeParts := strings.Fields(fullTypeStr) 
			if len(typeParts) > 0 {
				colInfo.BaseType = typeParts[0]
			} else {
				colInfo.BaseType = fullTypeStr 
			}
		}
		
		columns[colName] = colInfo
		logrus.Debugf("Parsed column: Name=%s, BaseType=%s, Length=%d, Precision=%d, Scale=%d, FullType=%s",
			colInfo.Name, colInfo.BaseType, colInfo.Length, colInfo.Precision, colInfo.Scale, colInfo.FullType)
	}

	if len(columns) == 0 {
		logrus.Warnf("No columns parsed from DDL: %s. Check DDL structure and parsing regexes.", ddlFilePath)
	} else {
		logrus.Infof("Parsed %d columns from %s", len(columns), ddlFilePath)
	}
	return columns, nil
}

// parseStats parses a .stats.yaml file, expecting a DBSchema structure,
// and returns the TableStats for the specified dbName and tableName.
func parseStats(statsFilePath, dbName, tableName string) (*TableStats, error) {
	logrus.Debugf("Parsing stats file: %s for db: %s, table: %s", statsFilePath, dbName, tableName)
	
	content, err := os.ReadFile(statsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("Stats file not found: %s. Proceeding without stats.", statsFilePath)
			return nil, nil // Not an error to proceed, but no stats available.
		}
		return nil, fmt.Errorf("failed to read stats file %s: %w", statsFilePath, err)
	}

	if len(content) == 0 {
		logrus.Warnf("Stats file %s is empty. Proceeding without stats.", statsFilePath)
		return nil, nil // Treat as no stats available
	}
	
	var dbSchema DBSchema // Assuming DBSchema is defined in the same package (src)
	err = yaml.Unmarshal(content, &dbSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats YAML from %s into DBSchema: %w", statsFilePath, err)
	}

	// Optional: Validate dbName if dbSchema.Name is populated and expected to match.
	// The current DBSchema yaml tag for Name is "db".
	if dbSchema.Name != "" && dbSchema.Name != dbName {
		logrus.Warnf("DB name in stats file ('%s') does not match expected DB name ('%s') for file %s. Stats might not be relevant.", dbSchema.Name, dbName, statsFilePath)
		// Depending on strictness, could return an error here or just proceed to find table.
	}

	for _, ts := range dbSchema.Stats { // DBSchema.Stats is []*TableStats
		if ts.Name == tableName {
			logrus.Infof("Successfully found and parsed stats for table %s.%s from %s", dbName, tableName, statsFilePath)
			return ts, nil
		}
	}

	logrus.Warnf("Table '%s' not found in stats file %s (DB schema '%s'). Proceeding without table-specific stats.", tableName, statsFilePath, dbSchema.Name)
	return nil, nil // Table not found, but not necessarily an error to proceed without its stats.
}


func generateTableData(csvFilePath string, columnDefs map[string]ColumnInfo, tableStats *TableStats, numRows int) error { 
	logrus.Debugf("Generating data for CSV: %s. Columns: %+v, NumRows: %d", csvFilePath, columnDefs, numRows)
	if tableStats != nil {
		logrus.Debugf("With Stats: %+v", tableStats)
	}


	if len(columnDefs) == 0 {
		logrus.Warnf("No column definitions found for %s, skipping CSV generation.", csvFilePath)
		return nil
	}

	file, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file %s: %w", csvFilePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	columnOrder := make([]string, 0, len(columnDefs))
	for colName := range columnDefs { 
		columnOrder = append(columnOrder, colName)
	}

	if err := writer.Write(columnOrder); err != nil {
		return fmt.Errorf("failed to write CSV header to %s: %w", csvFilePath, err)
	}
	
	for r := 0; r < numRows; r++ {
		record := make([]string, len(columnOrder))
		for i, colName := range columnOrder {
			colInfo := columnDefs[colName] 
			colType := strings.ToUpper(colInfo.BaseType) 
			var val string

			var currentColumnStats *ColumnStats
			if tableStats != nil {
				for _, cs := range tableStats.Columns {
					if cs.Name == colName {
						currentColumnStats = cs
						break
					}
				}
			}

			switch {
			case colType == "VARCHAR", colType == "CHAR":
				lowerColName := strings.ToLower(colInfo.Name)
				length := colInfo.Length
				if length <= 0 { length = 10 } // Default length if not specified or invalid for LetterN

				// Heuristic-based generation
				generatedByCategory := true
				switch {
				case strings.Contains(lowerColName, "email"):
					val = gofakeit.Email()
				case strings.Contains(lowerColName, "phone") || strings.Contains(lowerColName, "mobile"):
					val = gofakeit.PhoneFormatted()
				case strings.Contains(lowerColName, "city"):
					val = gofakeit.City()
				case strings.Contains(lowerColName, "country"):
					val = gofakeit.Country()
				case strings.Contains(lowerColName, "zip") || strings.Contains(lowerColName, "postal"):
					val = gofakeit.Zip()
				case strings.Contains(lowerColName, "url"):
					val = gofakeit.URL()
				case strings.Contains(lowerColName, "uuid") || (colName == "id" && length >=36) : // Also check for typical ID column name
					val = gofakeit.UUID()
				case strings.Contains(lowerColName, "street"): // Must be before general "address"
					val = gofakeit.Street() // Corrected function name
				case strings.Contains(lowerColName, "address"):
					val = gofakeit.Address().Address
				case strings.Contains(lowerColName, "company") || strings.Contains(lowerColName, "organization"):
					val = gofakeit.Company()
				case strings.Contains(lowerColName, "job") || strings.Contains(lowerColName, "occupation"):
					val = gofakeit.JobTitle()
				case strings.Contains(lowerColName, "first") && strings.Contains(lowerColName, "name"):
					val = gofakeit.FirstName()
				case strings.Contains(lowerColName, "last") && strings.Contains(lowerColName, "name"):
					val = gofakeit.LastName()
				case strings.Contains(lowerColName, "full") && strings.Contains(lowerColName, "name"):
					val = gofakeit.Name()
				case strings.Contains(lowerColName, "name"): // General name, if no "first/last/full"
					val = gofakeit.Name()
				case strings.Contains(lowerColName, "title") && !strings.Contains(lowerColName, "job"): // Avoid job title
					val = gofakeit.BookTitle() 
				case strings.Contains(lowerColName, "subject"):
					val = gofakeit.Sentence(gofakeit.Number(3, 7)) // 3 to 7 words
				// Description/Comment like fields for VARCHAR/CHAR should be shorter than TEXT
				case strings.Contains(lowerColName, "desc") || strings.Contains(lowerColName, "comment") || 
				     strings.Contains(lowerColName, "note") || strings.Contains(lowerColName, "detail"):
					val = gofakeit.Sentence(gofakeit.Number(5, 15)) // 5 to 15 words
				default:
					generatedByCategory = false
					// Fallback to random letters if no specific heuristic matches
					// Cap length for LetterN to avoid excessive strings if DDL length is very large
					letterNLength := length
					if letterNLength > 500 { letterNLength = 500 } 
					val = gofakeit.LetterN(uint(letterNLength))
				}

				// Truncate if necessary, only for VARCHAR/CHAR where length is defined
				if colInfo.Length > 0 && len(val) > colInfo.Length {
					// For UUIDs or specific formats, truncation might break them.
					// However, DDL length constraint is king.
					// For UUID, if length is < 36, this will truncate it.
					if generatedByCategory && (strings.Contains(lowerColName, "uuid") || strings.Contains(lowerColName, "phone")) {
						// Maybe warn if a formatted string is being truncated due to DDL length
						logrus.Warnf("Formatted value for column '%s' might be truncated due to DDL length %d. Value: '%s'", colName, colInfo.Length, val)
					}
					val = val[:colInfo.Length]
				}

			case colType == "TEXT", colType == "BLOB", colType == "CLOB", colType == "STRING":
				lowerColName := strings.ToLower(colInfo.Name)
				if strings.Contains(lowerColName, "desc") || strings.Contains(lowerColName, "comment") || 
				   strings.Contains(lowerColName, "note") || strings.Contains(lowerColName, "detail") || 
				   strings.Contains(lowerColName, "text") || strings.Contains(lowerColName, "message") ||
				   strings.Contains(lowerColName, "content") {
					val = gofakeit.Paragraph(gofakeit.Number(1,3), gofakeit.Number(2,5), gofakeit.Number(10,20), ". ")
				} else if colName == "id" || strings.HasSuffix(strings.ToLower(colName), "_id") {
					val = gofakeit.UUID()
				}

			case colType == "INT", colType == "BIGINT", colType == "SMALLINT", colType == "TINYINT":
				minStatVal, maxStatVal := 0, 1000000 // Default range
				useStatRange := false
				if currentColumnStats != nil && currentColumnStats.Min != "" && currentColumnStats.Max != "" {
					parsedMin, errMin := strconv.Atoi(currentColumnStats.Min)
					parsedMax, errMax := strconv.Atoi(currentColumnStats.Max)
					if errMin == nil && errMax == nil && parsedMin <= parsedMax {
						minStatVal, maxStatVal = parsedMin, parsedMax
						useStatRange = true
					} else {
						logrus.Warnf("Invalid Min/Max ('%s', '%s') for INT column '%s'. Using default range.", currentColumnStats.Min, currentColumnStats.Max, colName)
					}
				}
				
				if useStatRange {
					val = fmt.Sprintf("%d", gofakeit.Number(minStatVal, maxStatVal))
				} else { // Fallback or heuristic based generation
					if colName == "id" || strings.HasSuffix(strings.ToLower(colName), "_id") {
						val = fmt.Sprintf("%d", gofakeit.Number(1, 100000))
					} else if strings.Contains(strings.ToLower(colName), "age") {
						val = fmt.Sprintf("%d", gofakeit.Number(1, 100))
					} else if strings.Contains(strings.ToLower(colName), "year") {
						val = fmt.Sprintf("%d", gofakeit.Year())
					} else {
						val = fmt.Sprintf("%d", gofakeit.Number(minStatVal, maxStatVal)) // Use default range
					}
				}

			case colType == "DECIMAL", colType == "NUMERIC", colType == "FLOAT", colType == "DOUBLE":
				minFStatVal, maxFStatVal := 0.0, 100000.0 // Default range
				useFloatStatRange := false
				if currentColumnStats != nil && currentColumnStats.Min != "" && currentColumnStats.Max != "" {
					parsedMin, errMin := strconv.ParseFloat(currentColumnStats.Min, 64)
					parsedMax, errMax := strconv.ParseFloat(currentColumnStats.Max, 64)
					if errMin == nil && errMax == nil && parsedMin <= parsedMax {
						minFStatVal, maxFStatVal = parsedMin, parsedMax
						useFloatStatRange = true
					} else {
						logrus.Warnf("Invalid Min/Max ('%s', '%s') for NUMERIC/FLOAT column '%s'. Using default range.", currentColumnStats.Min, currentColumnStats.Max, colName)
					}
				}

				numToFormat := gofakeit.Float64Range(minFStatVal, maxFStatVal)
				if !useFloatStatRange { // If not using stats, calculate max based on precision
					maxCalculated := 100000.0
					if colInfo.Precision > 0 && colInfo.Precision > colInfo.Scale {
						power := float64(colInfo.Precision - colInfo.Scale)
						maxCalculated = 1.0
						for i := 0; i < int(power); i++ { maxCalculated *= 10 }
						if maxCalculated > 1e12 { maxCalculated = 1e12 }
					} else if colInfo.Precision > 0 && colInfo.Scale == 0 {
						maxCalculated = 1.0
						for i := 0; i < colInfo.Precision; i++ { maxCalculated *= 10 }
						if maxCalculated > 1e12 { maxCalculated = 1e12 }
					}
					// If stats were not used, regenerate number within calculated range if smaller
					if maxCalculated < maxFStatVal { // Check if precision-based max is more restrictive
						numToFormat = gofakeit.Float64Range(minFStatVal, maxCalculated) // minFStatVal is likely 0 if no stats
					}
				}
				
				format := "%f"
				if colType == "DECIMAL" || colType == "NUMERIC" {
					format = fmt.Sprintf("%%.%df", colInfo.Scale)
				} else if colInfo.Scale > 0 { // For FLOAT/DOUBLE if scale was parsed (e.g. from FLOAT(P,S))
					format = fmt.Sprintf("%%.%df", colInfo.Scale)
				}
				val = fmt.Sprintf(format, numToFormat)

				if strings.Contains(strings.ToLower(colName), "price") || strings.Contains(strings.ToLower(colName), "amount") {
					priceFormat := "%.2f"
					if colInfo.Scale > 0 && (colType == "DECIMAL" || colType == "NUMERIC") {
						priceFormat = fmt.Sprintf("%%.%df", colInfo.Scale)
					}
					val = fmt.Sprintf(priceFormat, gofakeit.Price(minFStatVal, maxFStatVal))
				}
				
			case colType == "DATE":
				minTime, maxTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()
				useDateStatRange := false // Corrected: remove underscore to match previous declaration if it was a typo, or ensure consistency.
				                        // The previous build error implies these specific variable names were "declared and not used".
										// Let's assume the declaration was `useDateStatRange` and `useDateTimeStatRange` (no underscore).
				if currentColumnStats != nil && currentColumnStats.Min != "" && currentColumnStats.Max != "" {
					parsedMin, errMin := parseDateFlexible(currentColumnStats.Min)
					parsedMax, errMax := parseDateFlexible(currentColumnStats.Max)
					if errMin == nil && errMax == nil && (parsedMin.Before(parsedMax) || parsedMin.Equal(parsedMax)) {
						minTime, maxTime = parsedMin, parsedMax
						useDateStatRange = true
					} else {
						logrus.Warnf("Invalid Min/Max ('%s', '%s') for DATE column '%s'. Using default range.", currentColumnStats.Min, currentColumnStats.Max, colName)
					}
				}
				if useDateStatRange { 
					val = gofakeit.DateRange(minTime, maxTime).Format("2006-01-02")
				} else {
					val = gofakeit.Date().Format("2006-01-02")
				}


			case colType == "DATETIME", colType == "TIMESTAMP":
				minTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
				maxTime := time.Now()
				useDateTimeStatRange := false // Corrected: remove underscore
				if currentColumnStats != nil && currentColumnStats.Min != "" && currentColumnStats.Max != "" {
					parsedMin, errMin := parseDateTimeFlexible(currentColumnStats.Min)
					parsedMax, errMax := parseDateTimeFlexible(currentColumnStats.Max)
					if errMin == nil && errMax == nil && (parsedMin.Before(parsedMax) || parsedMin.Equal(parsedMax)) {
						minTime, maxTime = parsedMin, parsedMax
						useDateTimeStatRange = true
					} else {
						logrus.Warnf("Invalid Min/Max ('%s', '%s') for DATETIME/TIMESTAMP column '%s'. Using default range.", currentColumnStats.Min, currentColumnStats.Max, colName)
					}
				}
				if useDateTimeStatRange { 
					val = gofakeit.DateRange(minTime, maxTime).Format("2006-01-02 15:04:05")
				} else {
					val = gofakeit.Date().Format("2006-01-02 15:04:05")
				}

			case colType == "BOOLEAN", colType == "BOOL":
				val = strconv.FormatBool(gofakeit.Bool())
			case colType == "UUID": // Added case for UUID type
				val = gofakeit.UUID()
			default:
				logrus.Warnf("Unsupported SQL BaseType '%s' for column '%s' in %s. Generating a random word.", colType, colName, csvFilePath)
				val = gofakeit.Word()
			}
			record[i] = val
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write data row to CSV %s: %w", csvFilePath, err)
		}
	}
	
	logrus.Infof("Successfully generated %d rows into %s", numRows, csvFilePath)
	return nil
}

// Gendata is a placeholder function for the gendata command's logic.
// This function might be removed or refactored if GenerateData is called directly
// from the cmd package.
func Gendata() {
	fmt.Println("src.Gendata called (placeholder) - this might be deprecated")
}

// parseDateFlexible tries to parse a date string using common layouts.
func parseDateFlexible(dateStr string) (time.Time, error) {
	layouts := []string{"2006-01-02", "2006/01/02", "2006-1-2", "2006/1/2"} // Add more if needed
	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse date: %s with any known layouts", dateStr)
}

// parseDateTimeFlexible tries to parse a datetime string using common layouts.
func parseDateTimeFlexible(dateTimeStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006/01/02 15:04:05",
		"2006-01-02 15:04", "2006-01-02T15:04",
		"2006-01-02", // Also allow just date for datetime types if time is omitted
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, dateTimeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse datetime: %s with any known layouts", dateTimeStr)
}
