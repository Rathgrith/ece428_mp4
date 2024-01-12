# ECE428_MP3

## Project Layout

`cmd/`: Supports commands via gRPC requests, including `uploader`, `downloader`, `scheduler_server/client`, and `nodebooter`

`java/src` Hadoop mapreduce codes for testing.

`pkg/`: Library code that can be used in application code (and also be used inside the pkg package), including request handler for the datanode and namenode.

`scripts/`: Scripts to perform build, deploy, performance measure, etc operations.

## Overview

This repository currently contains MP4 code for this course. This repository consists of one major application:

  - `cmd/node/boot.go`: It can start nodemanager service.
  - `cmd/scheduler_server/Scheduler.go`: It can dispatch operation commands by starting the jobTracker to all servers and return aggregated results to the user.
  - `cmd/scheduler_client/SchedulerClient.go`: It parses the SQL command into maple/juice tasks to the server(resource manager)

## Getting started


Call the following commands to activate corresponding processes, and make sure you have MP3 SDFS running.

```go
go run ./cmd/node/boot.go
```

Then run a resourceManager on the master node by:

```go
go run cmd/scheduler_server/Scheduler.go 
```

Finally, run your SQL command at any machine by:

```go
go run cmd/scheduler_client/SchedulerClient.go -sql command -out_dir output.csv
```


## Deployment

- Simply use ssh to clone your repository using ssh.

  ```shell
  ssh dl58fa23-cs425-48XX.cs.illinois.edu
  ssh> cd ./go
  ssh> git clone https://gitlab.engr.illinois.edu/dl58/ece428_mp4.git
  ```
  

