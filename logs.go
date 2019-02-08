package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

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
	printColor(0, 31, log)
}
func printSuccess(log string) {
	printColor(0, 32, log)
}
func printInfo(log string) {
	printColor(0, 34, log)
}

func printColor(fontType int, fontColor int, text string) {
	fmt.Printf("\033["+strconv.Itoa(fontType)+";"+strconv.Itoa(fontColor)+"m%s\033[0m \n", text)
}
