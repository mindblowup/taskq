package main

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
)

var (
	flagListen          = flag.String("listen", ":8001", "the http listen address")
	flagSecretToken     = flag.String("token", tokenGenerator(), "Secret token becuse taskq is end to end ")
	flagTimeout         = flag.Int("timeout", 30, "How many times of seconds limit for waiting a response before cancelling the request. 0 means no limit")
	flagRetry           = flag.Int("retry", 3, "How many retry to execute the task when failed")
	flagRetryDelay      = flag.Int("retry_delay", 30, "Number of seconds you need to delay and retry to executing the task when failed")
	flagFailureCallback = flag.String("failure_callback", "", "Api url will be called (POST request) when task failure after retries.")
	flagErrorLog        = flag.String("error_log", "./taskq_error.log", "error logs")
	//flagCompleteLog     = flag.String("complete_log", "./taskq_complete.log", "comlplete log")
)

const (
	version = "v1.0.0"
	welcome = `
		 _____         _     ___              
		|_   _|_ _ ___| | __/ _ \             
		  | |/ _' / __| |/ / | | |            
		  | | (_| \__ \   <| |_| |            
		  |_|\__'_|___/_|\_\\__\_\  ` + version + `   

 Asynchronously task queues over http request
`
)

func main() {
	flag.Parse()
	initDB()
	go runTasks()
	fmt.Println("\033[" + strconv.Itoa(1) + ";" + strconv.Itoa(35) + "m" + welcome + "\033[0m")
	runServer()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
