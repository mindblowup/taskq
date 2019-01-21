package main

import (
	"fmt"
	"github.com/nezamy/jsondb"
	"log"
	"os"
	"strconv"
)

var db *jsondb.Driver

func initDB() {
	var err error
	db, err = jsondb.New(".taskq")
	if err != nil {
		printError(err.Error())
	}
}

func addTaskLog(channel string, id string, q Queue) {
	if err := db.Write(channel, id, q); err != nil {
		printError(err.Error())
	}
}

func addErrorLog(log string) {
	addLog(*flagErrorLog, log)
	printError(log)
}

//func addCompleteLog (log interface{}){
//	addLog(*flagCompleteLog, log)
//}

func addLog(path string, cont interface{}) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(cont)
}

func printError(log string) {
	fmt.Println("\033[" + strconv.Itoa(0) + ";" + strconv.Itoa(31) + "m" + log + "\033[0m")
}
func printSuccess(log string) {
	fmt.Println("\033[" + strconv.Itoa(0) + ";" + strconv.Itoa(32) + "m" + log + "\033[0m")
}
func printInfo(log string) {
	fmt.Println("\033[" + strconv.Itoa(0) + ";" + strconv.Itoa(34) + "m" + log + "\033[0m")
}
