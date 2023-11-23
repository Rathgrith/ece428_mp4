package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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
		handleJuice(*juiceExe, *numJuice, *sdfsPrefix, *sdfsOutput)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func handleMaple(exe string, numMaples int, prefix string, srcDir string) {
	// TODO: list all the machines containing source files
	// TODO: schedule the maple executable to run on the files
	// Implement Maple phase logic here
	fmt.Printf("Running Maple with exe: %s, num: %d, prefix: %s, srcdir: %s\n", exe, numMaples, prefix, srcDir)

	// Execute the maple executable
	cmd := exec.Command(exe, "-input", srcDir, "-prefix", prefix, "-regex", "Anthony")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}
	fmt.Printf("Maple executable output:\n%s\n", string(output))
	// TODO: put the file into sdfs
}

func handleJuice(exe string, numMaples int, prefix string, outDir string) {
	// Implement Juice phase logic here
	// TODO: list all the files in SDFS with the given prefix
	// TODO: schedule the juice executable to run on the files (once a file for each machine)
	fmt.Printf("Running Juice with exe: %s, num: %d, prefix: %s, output: %s\n", exe, numMaples, prefix, outDir)
	cmd := exec.Command(exe, "-prefix", prefix, "-output", outDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}
	fmt.Printf("Juice executable output:\n%s\n", string(output))
	// TODO: put the file into sdfs
}

// Implement other functions like handleJuice, partitioning, etc.
