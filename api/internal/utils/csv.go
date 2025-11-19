package utils

import (
	"encoding/csv"
	"fmt"
	"io"
)

func ReadCsv(r io.Reader, targetColumn string) ([]string, error) {
	reader := csv.NewReader(r)

	records, err := reader.ReadAll()

	if err != nil {
		return []string{}, fmt.Errorf("fail to read records from csv file")
	}

	if len(records) == 0 {
		return []string{}, fmt.Errorf("no records found for csv file")
	}

	// get column names
	cols := records[0]

	colToIndex := map[string]int{}
	for index, item := range cols {
		colToIndex[item] = index
	}

	targetIndex, isExist := colToIndex[targetColumn]
	if !isExist {
		return []string{}, fmt.Errorf("column: %s not found in the csv file", targetColumn)
	}

	results := []string{}

	for _, row := range records[1:] { // skip the first row
		for index, col := range row {
			if index == targetIndex {
				results = append(results, col)
			}
		}
	}

	return results, nil
}
