package password

import (
	"encoding/csv"
	"fmt"
	"os"
)

const csvFilename = "passwords.csv"

// SavePasswordToCSV appends the password to a CSV file
func SavePasswordToCSV(siteName, username, password string) error {
	file, err := os.OpenFile(csvFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{siteName, username, password}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}

// LoadPasswordsFromCSV reads all passwords from the CSV file
func LoadPasswordsFromCSV() ([][]string, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return [][]string{}, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Ensure all records have at least 3 columns (backward compatibility)
	for i, record := range records {
		for len(record) < 3 {
			record = append(record, "")
		}
		records[i] = record
	}

	return records, nil
}
