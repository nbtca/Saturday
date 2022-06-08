package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func GetCsvMap(filePath string) ([]map[string]string, error) {
	records := ReadCsvFile(filePath)
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
