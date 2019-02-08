package main

import (
	"fmt"
	"github.com/rs/xid"
	"sync"
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
	if _, ok := channels[channel]; !ok {
		channels[channel] = make(Tasks)
		tasksTotal[channel] = 0
		tasksComplete[channel] = 0
	}

	if queue.Id == "" {
		guid := xid.New()
		queue.Id = guid.String()
		if err := db.Write(channel, queue.Id, queue); err != nil {
			addErrorLog(err.Error())
		}
	}
	channels[channel][queue.Id] = &queue
	tasksTotal[channel]++
	fmt.Println("Add task", queue.Name, "To Channel", channel)
	return queue
}

func removeTask(queue Tasks, channel string, id string) {
	delete(queue, id)
	if err := db.Delete(channel, id); err != nil {
		fmt.Println("Error", err)
	}
}

func runTasks() {
	var WG sync.WaitGroup
	for {
		count := len(channels)
		if count > 0 {
			WG.Add(count)
			for channel, _ := range channels {
				time.Sleep(time.Millisecond * 100)
				go runChannelTasks(channel, &WG)

			}
			WG.Wait()
		}

	}
}

func runChannelTasks(channel string, WG *sync.WaitGroup) {
	defer WG.Done()
	for key, val := range channels[channel] {
		if val.IsNeedToExecuteNow() {
			fmt.Println("-----------------------------------------------------------")
			fmt.Printf("Executing : %s (%s) \n", val.Name, channel)
			if val.exec() {
				printSuccess(val.Name + " (" + channel + ") is Executed")
			}

			if val.IsExpired() {
				removeTask(channels[channel], channel, key)
				tasksComplete[channel]++
				fmt.Printf("%s (%s) is Done - Tasks in progress(%d) - Total Complete(%d)", val.Name, channel, tasksTotal[channel]-tasksComplete[channel], tasksComplete[channel])
			}
		}
	}
}
