package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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
	sdfsSrcDir := mapleCmd.String("srcdir", "", "SDFS source directory")

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
		handleMaple(*mapleExe, *numMaples, *sdfsIntermediatePrefix, *sdfsSrcDir)
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

func handleMaple(exe string, numMaples int, prefix string, srcDir string) {
	client, err := createMapleJuiceClient("localhost:50051") // Adjust the address as needed
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	// Prepare request
	req := &idl.MapleRequest{
		Exe:       exe,
		NumMaples: int32(numMaples),
		Prefix:    prefix,
		SrcDir:    srcDir,
	}

	// Make RPC call
	resp, err := client.ExecuteMaple(context.Background(), req)
	if err != nil {
		fmt.Printf("RPC failed: %v\n", err)
		return
	}

	fmt.Printf("Maple Response: %s\n", resp.GetOutput())
}

func handleJuice(exe string, numJuices int, prefix string, output string, deleteIntermediate string) {
	client, err := createMapleJuiceClient("localhost:50051") // Adjust the address as needed
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	// local logic to find the keys
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
	key := keys[0] // Using the first key for the RPC call

	fmt.Printf("Running Juice with exe: %s, num: %d, prefix: %s, output: %s\n", exe, numJuices, prefix, output)

	// Prepare request for RPC
	req := &idl.JuiceRequest{
		Exe:         exe,
		Key:         key,
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

// func executeMaple(exe string, numMaples int, prefix string, srcDir string) {
// 	cmd := exec.Command(exe, "-input", srcDir, "-prefix", prefix, "-regex", "Anthony")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Failed to execute command: %s\n", err)
// 		return
// 	}
// 	fmt.Printf("Maple executable output:\n%s\n", string(output))
// }

// func executeJuice(exe string, key string, prefix string, outDir string) {
// 	cmd := exec.Command(exe, "-key", key, "-prefix", prefix, "-output", outDir)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Failed to execute command: %s\n", err)
// 		return
// 	}
// 	fmt.Printf("Juice executable output:\n%s\n", string(output))
// }

// func handleMaple(exe string, numMaples int, prefix string, srcDir string) {
// 	// TODO: list all the machines containing source files
// 	// TODO: schedule the maple executable to run on the files
// 	// Implement Maple phase logic here
// 	fmt.Printf("Running Maple with exe: %s, num: %d, prefix: %s, srcdir: %s\n", exe, numMaples, prefix, srcDir)

// 	// Execute the maple executable
// 	executeMaple(exe, numMaples, prefix, srcDir)
// 	// TODO: put the file into sdfs
// }

// func handleJuice(exe string, numMaples int, prefix string, outDir string) {
// 	// Implement Juice phase logic here
// 	// TODO: list all the files in SDFS with the given prefix
// 	// TODO: schedule the juice executable to run on the files (once a file for each machine)
// 	files, err := os.ReadDir(".")
// 	if err != nil {
// 		fmt.Println("Error reading directory:", err)
// 		return
// 	}
// 	var keys []string
// 	for _, file := range files {
// 		if file.IsDir() {
// 			continue
// 		}
// 		// avoid slice bounds out of range
// 		if len(file.Name()) > len(prefix)+1 {
// 			if file.Name()[:len(prefix)] == prefix {
// 				keys = append(keys, file.Name()[len(prefix)+1:])
// 			}
// 		}
// 	}
// 	// currently local, should be sdfs
// 	fmt.Println("pending keys:")
// 	for _, key := range keys {
// 		fmt.Println("key:", key)
// 	}
// 	// Validate input
// 	if len(keys) == 0 {
// 		fmt.Println("No files found with the given prefix")
// 		return
// 	}
// 	key := keys[0]
// 	fmt.Println("processing key:", key)
// 	fmt.Printf("Running Juice with exe: %s, num: %d, prefix: %s, output: %s\n", exe, numMaples, prefix, outDir)
// 	executeJuice(exe, key, prefix, outDir)

// TODO: put the file into sdfs
// }
