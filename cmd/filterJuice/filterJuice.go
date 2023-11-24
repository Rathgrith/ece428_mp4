package main

import (
	"ece428_mp4/pkg"
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/theckman/go-flock"
)

// Function to read a CSV file and return its contents
func readCSV(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			// EOF is expected when all records have been read
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		// Assuming the key is the entire line and value is in the second column
		lines = append(lines, record[0])
	}

	return lines, nil
}

func juice(key, intermediatePrefix, outputFilename string) error {
	filename := fmt.Sprintf("%s_%s", intermediatePrefix, key)
	fmt.Println("filename:", filename)

	lines, err := readCSV(filename)
	if err != nil {
		return fmt.Errorf("error processing file %s: %v", filename, err)
	}

	// Create a new file lock for the output file
	fileLock := flock.New(outputFilename + ".lock")

	locked, err := fileLock.TryLock()
	if err != nil {
		return fmt.Errorf("error acquiring file lock: %v", err)
	}

	if !locked {
		return fmt.Errorf("unable to acquire lock on file: %s", outputFilename)
	}

	// Ensure the file lock is released
	defer fileLock.Unlock()

	// Now you can safely append to the CSV file
	err = pkg.AppendToCSV(outputFilename, lines)
	if err != nil {
		return fmt.Errorf("error appending to CSV: %v", err)
	}

	return nil
}

func main() {
	// TODO: Read the intermediate files from SDFS to local

	// Define command line flags
	key := flag.String("key", "", "Keys to process")
	intermediatePrefix := flag.String("prefix", "", "Prefix for the intermediate CSV files")
	outputFilename := flag.String("output", "", "Filename for the output CSV file")
	deleteIntermediate := flag.String("delete", "0", "Delete intermediate files after processing")

	// Parse the command line flags
	flag.Parse()
	// Find files with the given prefix locally
	if *key == "" || *intermediatePrefix == "" || *outputFilename == "" {
		fmt.Println("All flags (key, prefix, output) are required")
		flag.Usage()
		return
	}
	// Execute the juice function
	err := juice(*key, *intermediatePrefix, *outputFilename)
	if err != nil {
		fmt.Println("Error executing juice function:", err)
	} else {
		fmt.Println("Juice function executed successfully.")
	}
	if *deleteIntermediate == "1" {
		fmt.Println("Deleting intermediate files")
		filename := fmt.Sprintf("%s_%s", *intermediatePrefix, *key)
		err := os.Remove(filename)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		}
	}
}
