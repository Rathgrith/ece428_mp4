package pkg

import (
	"encoding/csv"
	"fmt"
	"os"
)

// Append lines to a CSV file
func AppendToCSV(filename string, lines []string) error {
	outFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)

	for _, line := range lines {
		fmt.Println("Writing line to CSV:", line)
		if err := writer.Write([]string{line}); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}
	return nil
}
