package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func routes() {
	http.Handle("/list", httpMiddleware(tokenMiddleware)(list))
	http.Handle("/clear", httpMiddleware(tokenMiddleware)(clear))
	http.Handle("/add-http-task", httpMiddleware(tokenMiddleware)(addHttpTask))
}

func list(w http.ResponseWriter, r *http.Request) {
	res := ListResponse{true, channels, make(map[string]string)}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func clear(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		collections := r.URL.Query().Get("collection")
		var list []string
		if collections != "" {
			list = strings.Split(collections, ",")
		} else {
			list = db.List()
		}

		results := make([]interface{}, len(list))
		for i, collection := range list {
			collection = strings.Trim(collection, " ")
			if err := db.Delete(collection, ""); err != nil {
				responseError(w, map[string]string{collection: err.Error()})
				return
			}
			msg := "All tasks in channel (" + collection + ") are deleted successfully."
			results[i] = msg
			printSuccess(msg)
		}

		responseSuccess(w, results)
	}
}

func addHttpTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		channel := getChannel(r)
		var bodyRes []Queue
		parseBody(r, &bodyRes)
		var results = make([]interface{}, len(bodyRes))

		HttpTaskValidation(w, bodyRes)

		for i, t := range bodyRes {
			t.Parse()
			t.Headers = r.Header
			t = AddTask(channel, t)
			results[i] = map[string]string{"id": t.Id, "name": t.Name, "url": t.Url}
		}

		responseSuccess(w, results)
		fmt.Println("Tasks in progress:", tasksTotal[channel]-tasksComplete[channel], "Tasks Complete", tasksComplete[channel])
	}
}

//func addCommandTask() {
//
//}

func HttpTaskValidation(w http.ResponseWriter, queues []Queue) {
	for _, t := range queues {
		if t.Name == "" {
			responseError(w, map[string]string{"name": "name is required"})
			return
		}

		if t.Url == "" {
			responseError(w, map[string]string{"url": "url not defined for " + t.Name})
			return
		}
	}
}
