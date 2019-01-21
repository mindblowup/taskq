package main

import (
	"fmt"
	"github.com/rs/xid"
	"time"
)

var (
	channels      = make(Channels)
	tasksTotal    = map[string]int64{"default": 0}
	tasksComplete = map[string]int64{"default": 0}
)

type Channels map[string]Tasks
type Tasks map[string]*Queue

func AddTask(channel string, queue Queue) Queue {
	guid := xid.New()
	if _, ok := channels[channel]; !ok {
		channels[channel] = make(Tasks)
		tasksTotal[channel] = 0
		tasksComplete[channel] = 0
	}
	queue.Id = guid.String()
	channels[channel][guid.String()] = &queue
	addTaskLog(channel, guid.String(), queue)
	tasksTotal[channel]++
	fmt.Println("add task", queue.Name)
	return queue
}

func removeTask(queue Tasks, channel string, id string) {
	delete(queue, id)
	if err := db.Delete(channel, id); err != nil {
		fmt.Println("Error", err)
	}
}

func runTasks() {
	for {
		for channel, _ := range channels {
			time.Sleep(time.Millisecond * 100)
			for key, val := range channels[channel] {
				if val.IsNeedToExecuteNow() {
					fmt.Println("Excuting :", val.Name)
					if val.exec() {
						fmt.Println(val.Name, "Done")
					}

					if val.IsExpired() {
						removeTask(channels[channel], channel, key)
						tasksComplete[channel]++
						fmt.Println("executed :", val.Name, "Tasks in progress:", tasksTotal[channel]-tasksComplete[channel], "Total Complete : ", tasksComplete[channel])
					}
				}

				//PrintMemUsage()
			}
		}
	}
}
