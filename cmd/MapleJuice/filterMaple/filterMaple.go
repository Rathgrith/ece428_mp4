package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Function to generate a valid filename from a key
func generateFilename(prefix, key string) string {
	// Sanitize key to remove special characters or spaces
	sanitizedKey := strings.ReplaceAll(key, " ", "_")
	sanitizedKey = strings.ReplaceAll(sanitizedKey, ",", "_")
	// Add other sanitizations as needed

	return fmt.Sprintf("%s_%s.csv", prefix, sanitizedKey)
}

func createCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	file.Close()
	return nil
}

// Maple function to filter CSV lines based on regex and output to separate files per key
func maple(inputFile, prefix, regexPattern string) error {
	// Open the input file
	inFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer inFile.Close()

	// Create a new CSV reader
	reader := csv.NewReader(inFile)

	// Compile the regular expression
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return fmt.Errorf("error compiling regex: %v", err)
	}

	// A map to track which keys (output files) have been encountered
	encounteredKeys := make(map[string]bool)

	// Process each record (line) from the CSV
	for {
		record, err := reader.Read()
		if err != nil {
			// EOF is expected when all records have been read
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading CSV record: %v", err)
		}

		// Convert record (slice of fields) to a single string
		line := strings.Join(record, ",")

		if regex.MatchString(line) {
			// Generate the output filename
			// Set the key to the attribute value that matched the regex
			outputFile := generateFilename(prefix, record[0])

			for _, attr := range record {
				if regex.MatchString(attr) {
					outputFile = generateFilename(prefix, attr)
					break
				}
			}

			// Check if this key has been encountered before
			if outputFile != "" {
				if _, exists := encounteredKeys[outputFile]; !exists {
					// If this is a new key, create the file and mark it as encountered
					if err := createCSV(outputFile); err != nil {
						return err
					}
					encounteredKeys[outputFile] = true
				}
			}

			// Append the line to the corresponding file
			if err := appendToCSV(outputFile, line, "1"); err != nil {
				return err
			}
		}
	}
	return nil
}

// Function to append a line to a CSV file
func appendToCSV(filename, key, value string) error {
	outFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	if err := writer.Write([]string{key, value}); err != nil {
		return fmt.Errorf("error writing CSV record: %v", err)
	}

	return nil
}

// func main() {
// 	inputFile := "map/test2.csv"
// 	prefix := "select_Anthony"
// 	regexPattern := "Anthony"

//		err := maple(inputFile, prefix, regexPattern)
//		if err != nil {
//			fmt.Println("Error executing maple function:", err)
//		} else {
//			fmt.Println("Maple function executed successfully.")
//		}
//	}
func main() {
	// Define command line flags
	inputFile := flag.String("input", "", "Path to the input CSV file")
	prefix := flag.String("prefix", "", "Prefix for the output CSV files")
	regexPattern := flag.String("regex", "", "Regular expression to match lines in the CSV")

	// Parse the command line flags
	flag.Parse()

	// Validate input
	if *inputFile == "" || *prefix == "" || *regexPattern == "" {
		fmt.Println("All flags (input, prefix, regex) are required")
		flag.Usage()
		return
	}

	// Execute the maple function
	err := maple(*inputFile, *prefix, *regexPattern)
	if err != nil {
		fmt.Println("Error executing maple function:", err)
	} else {
		fmt.Println("Maple function executed successfully.")
	}
}
