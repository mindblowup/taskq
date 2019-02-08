package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nezamy/jsondb"
	"time"
)

var (
	flagListen          = flag.String("listen", ":8001", "the http listen address")
	flagSecretToken     = flag.String("token", tokenGenerator(), "Secret token because taskQ is end to end ")
	flagTimeout         = flag.Int("timeout", 20, "How many times of seconds limit for waiting a response before cancelling the request. 0 means no limit")
	flagRetry           = flag.Int("retry", 3, "How many retry to execute the task when failed")
	flagRetryDelay      = flag.Int("retry_delay", 30, "Number of seconds you need to delay and retry to executing the task when failed")
	flagFailureCallback = flag.String("failure_callback", "", "Api url will be called (POST request) when task failure after retries.")
	flagErrorLog        = flag.String("error_log", "./taskq_error.log", "error logs")
	flagClear           = flag.Bool("clear", false, "Clear all previous uncompleted tasks")
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

==============================================
`
)

func main() {
	flag.Parse()
	printColor(1, 35, welcome)
	initDB()
	if *flagClear {
		clearTasks()
	}
	loadTasks()
	go runTasks()
	runServer()
}

var db *jsondb.Driver

func initDB() {
	var err error
	db, err = jsondb.New(".taskq")
	if err != nil {
		printError(err.Error())
	}
}

func loadTasks() {
	printInfo("Load previous uncompleted tasks ...")
	time.Sleep(time.Second)
	list := db.List()
	for _, collection := range list {
		records, _ := db.ReadAll(collection)
		for _, f := range records {
			queue := Queue{}
			if err := json.Unmarshal([]byte(f), &queue); err != nil {
				fmt.Println("Error", err)
			}
			AddTask(collection, queue)
		}
	}
	time.Sleep(time.Second)
}

func clearTasks() {
	list := db.List()
	for _, collection := range list {
		if err := db.Delete(collection, ""); err != nil {
			printError(err.Error())
		}
		printInfo("Clearing all tasks in channel (" + collection + ") ...")
	}

}

func tokenGenerator() string {
	b := make([]byte, 50)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//func PrintMemUsage() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
//	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
//	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
//	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
//	fmt.Printf("\tNumGC = %v\n", m.NumGC)
//}
//
//func bToMb(b uint64) uint64 {
//	return b / 1024 / 1024
//}
