package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"ece428_mp4/idl"

	"google.golang.org/grpc"
)

func main() {
	// Command line arguments
	mapleCmd := flag.NewFlagSet("maple", flag.ExitOnError)
	juiceCmd := flag.NewFlagSet("juice", flag.ExitOnError)

	// Maple command arguments
	mapleExe := mapleCmd.String("exe", "", "Executable for maple phase")
	numMaples := mapleCmd.Int("num", 1, "Number of maples")
	sdfsIntermediatePrefix := mapleCmd.String("prefix", "", "SDFS intermediate filename prefix")
	sdfsSrcDir1 := mapleCmd.String("srcdir", "", "SDFS source directory")
	sdfsSrcDir2 := mapleCmd.String("srcdir2", "", "SDFS source directory")
	regexPattern := mapleCmd.String("regex", "", "Regular expression to match lines in the CSV")
	joinColumn1 := mapleCmd.Int("col1", -1, "First column to join on")
	joinColumn2 := mapleCmd.Int("col2", -1, "Second column to join on")

	juiceExe := juiceCmd.String("exe", "", "Executable for juice phase")
	numJuice := juiceCmd.Int("num", 1, "Number of juice")
	sdfsPrefix := juiceCmd.String("prefix", "", "SDFS intermediate filename prefix")
	sdfsOutput := juiceCmd.String("output", "", "SDFS source directory")
	deleteIntermediate := juiceCmd.String("delete", "0", "Delete intermediate files after juice")

	// Juice command arguments

	if len(os.Args) < 2 {
		fmt.Println("maple or juice subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "maple":
		mapleCmd.Parse(os.Args[2:])
		handleMaple(*mapleExe, *numMaples, *sdfsIntermediatePrefix, *sdfsSrcDir1, *sdfsSrcDir2, *regexPattern, *joinColumn1, *joinColumn2)
	case "juice":
		juiceCmd.Parse(os.Args[2:])
		handleJuice(*juiceExe, *numJuice, *sdfsPrefix, *sdfsOutput, *deleteIntermediate)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func createMapleJuiceClient(address string) (idl.MapleJuiceServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return idl.NewMapleJuiceServiceClient(conn), nil
}

func handleMaple(exe string, numMaples int, prefix string, srcDir1 string, srcDir2 string, regex string, col1 int, col2 int) {
	// TODO: list all machines except self
	// for each of them schedule the maple executable to run on the files
	timeStart := time.Now()
	client, err := createMapleJuiceClient("localhost:50051") // Adjust the address as needed
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	// Prepare request
	req := &idl.MapleRequest{
		Exe:         exe,
		NumMaples:   int32(numMaples),
		Prefix:      prefix,
		SrcDir1:     srcDir1,
		SrcDir2:     srcDir2,
		Regex:       regex,
		JoinColumn1: int32(col1),
		JoinColumn2: int32(col2),
	}

	// Make RPC call
	resp, err := client.ExecuteMaple(context.Background(), req)
	if err != nil {
		fmt.Printf("RPC failed: %v\n", err)
		return
	}

	fmt.Printf("Maple Response: %s\n", resp.GetOutput())
	timeEnd := time.Now()
	fmt.Println("Time to complete the task:", timeEnd.Sub(timeStart))
}

func handleJuice(exe string, numJuices int, prefix string, output string, deleteIntermediate string) {
	// record the time to complete the task
	timeStart := time.Now()
	// local logic to find the keys
	// TODO: list all files with the given prefix in SDFS
	// TODO: schedule the juice executable to run on the files (once a file for each machine)
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	var keys []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if len(file.Name()) > len(prefix)+1 {
			if file.Name()[:len(prefix)] == prefix {
				keys = append(keys, file.Name()[len(prefix)+1:])
			}
		}
	}
	fmt.Println("pending keys:")
	for _, key := range keys {
		fmt.Println("key:", key)
	}
	if len(keys) == 0 {
		fmt.Println("No files found with the given prefix")
		return
	}
	// key := keys[0] // Using the first key for the RPC call
	// need goroutine to handle each keys
	for _, k := range keys {
		client, err := createMapleJuiceClient("localhost:50051")
		if err != nil {
			fmt.Printf("Failed to create client: %v\n", err)
			return
		}

		fmt.Printf("Running Juice with exe: %s, num: %d, prefix: %s, output: %s\n", exe, numJuices, prefix, output)

		// Prepare request for RPC
		req := &idl.JuiceRequest{
			Exe:         exe,
			Key:         k,
			Prefix:      prefix,
			OutDir:      output,
			DeleteInput: deleteIntermediate,
		}

		// Make RPC call
		resp, err := client.ExecuteJuice(context.Background(), req)
		if err != nil {
			fmt.Printf("RPC failed: %v\n", err)
			return
		}
		fmt.Printf("Juice Response: %s\n", resp.GetOutput())
	}
	timeEnd := time.Now()
	fmt.Println("Time to complete the task:", timeEnd.Sub(timeStart))
}
