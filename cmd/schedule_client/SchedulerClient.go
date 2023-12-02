package main

import (
	"context"
	"ece428_mp4/idl"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type SQLQuery struct {
	Operation   string
	Datasets    []string
	Condition   string
	JoinColumn1 string
	JoinColumn2 string
}

func parseSQLQuery(sql string) (SQLQuery, error) {
	parts := strings.Fields(strings.ToLower(sql))
	// fmt.Println(parts[3])
	// fmt.Println(parts[5])
	if len(parts) < 6 {
		return SQLQuery{}, errors.New("invalid query format")
	}

	query := SQLQuery{}
	if parts[2] == "from" && parts[4] == "where" {
		// Filter operation
		query.Operation = "filter"
		query.Datasets = []string{parts[3]}
		query.Condition = strings.Join(parts[5:], " ")
	} else if len(parts) >= 7 && strings.Contains(parts[3], ",") && parts[5] == "where" {
		// Join operation
		query.Operation = "join"
		dataset1 := strings.Split(parts[3], ",")[0]
		query.Datasets = []string{dataset1, parts[4]}
		query.JoinColumn1 = strings.Split(parts[6], ".")[1]
		query.JoinColumn2 = strings.Split(parts[8], ".")[1]
	} else {
		return SQLQuery{}, errors.New("unsupported query type")
	}

	return query, nil
}

func main() {
	serverAddr := flag.String("server_addr", "fa23-cs425-4801.cs.illinois.edu:50051", "The server address in the format of host:port")
	numJobs := flag.Int("numJobs", 1, "Number of maple tasks")
	prefix := flag.String("prefix", "inter", "Prefix for intermediate files")
	srcDir1 := flag.String("src_dir1", "./test1.csv", "Source directory 1 for maple/juice task")
	srcDir2 := flag.String("src_dir2", "", "Source directory 2 for juice task (optional)")
	regex := flag.String("regex", "Bloom", "Regex pattern for filter (optional)")
	joinColumn1 := flag.String("join_col1", "-1", "Join column for the first dataset (optional)")
	joinColumn2 := flag.String("join_col2", "-1", "Join column for the second dataset (optional)")
	outDir := flag.String("out_dir", "output.csv", "Output directory for juice task")
	sqlQuery := flag.String("sql", "", "SQL-like query for filter operation")

	flag.Parse()

	// Parse the SQL query if provided
	var query SQLQuery
	var err error
	if *sqlQuery != "" {
		query, err = parseSQLQuery(*sqlQuery)
		if err != nil {
			//log.Fatalf("Invalid SQL query: %v", err)
		}
		if query.Operation == "filter" {
			*srcDir1 = query.Datasets[0]
			*regex = query.Condition
		}
		if query.Operation == "join" {
			*srcDir1 = query.Datasets[0]
			*srcDir2 = query.Datasets[1]
			*joinColumn1 = query.JoinColumn1
			*joinColumn2 = query.JoinColumn2
		}
	}

	// Convert join column strings to integers
	joinCol1, err := strconv.Atoi(*joinColumn1)
	if err != nil {
		joinCol1 = -1
	}
	joinCol2, err := strconv.Atoi(*joinColumn2)
	if err != nil {
		joinCol2 = -1
	}
	// Set up a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := idl.NewMapleJuiceSchedulerClient(conn)

	// Contact the server and print out its response
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// query.Operation = "filter"
	fmt.Printf("Sending query: %+v\n", query)
	if query.Operation == "filter" {
		mapleResp, err := client.EnqueueTask(ctx, &idl.TaskRequest{
			TaskType: "maple",
			Exe:      "filterMaple.exe",
			NumJobs:  int32(*numJobs),
			Prefix:   *prefix,
			SrcDir1:  *srcDir1,
			Regex:    *regex,
			DestFile: *outDir,
		})
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		fmt.Printf("Server response: %s\n", mapleResp.GetMessage())
		juiceResp, err := client.EnqueueTask(ctx, &idl.TaskRequest{
			TaskType: "juice",
			Exe:      "./filterJuice.exe",
			NumJobs:  int32(*numJobs),
			Prefix:   *prefix,
			DestFile: *outDir,
		})
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		fmt.Printf("Server response: %s\n", juiceResp.GetMessage())
	} else if query.Operation == "join" {
		joinmResp, err := client.EnqueueTask(ctx, &idl.TaskRequest{
			TaskType:    "maple",
			Exe:         "./joinMaple.exe",
			NumJobs:     int32(*numJobs),
			SrcDir1:     *srcDir1,
			SrcDir2:     *srcDir2,
			Prefix:      *prefix,
			JoinColumn1: int32(joinCol1),
			JoinColumn2: int32(joinCol2),
			DestFile:    *outDir,
		})
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		fmt.Printf("Server response: %s\n", joinmResp.GetMessage())
		joinjResp, err := client.EnqueueTask(ctx, &idl.TaskRequest{
			TaskType: "juice",
			Exe:      "./joinJuice.exe",
			NumJobs:  int32(*numJobs),
			Prefix:   *prefix,
			DestFile: *outDir,
		})
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		fmt.Printf("Server response: %s\n", joinjResp.GetMessage())
	}
}
