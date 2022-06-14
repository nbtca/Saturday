package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCsvFile reads a csv file and returns a slice of records
func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Unable to parse file as CSV for "+filePath, err)
	}
	return records, nil
}

// ReadCsvFileWithHeader reads a csv file and returns a map of header and data
func GetCsvMap(filePath string) ([]map[string]string, error) {
	records, err := ReadCsvFile(filePath)
	if err != nil {
		return nil, err
	}
	target := make([]map[string]string, len(records)-1)
	if len(records) < 2 {
		return nil, fmt.Errorf("file %s does not have enough records", filePath)
	}
	headers := records[0]
	for i := 1; i < len(records); i++ {
		target[i-1] = make(map[string]string)
		for j, k := range headers {
			target[i-1][k] = records[i][j]
		}
	}
	return target, nil
}
