package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Status bool              `json:"status"`
	Data   []interface{}     `json:"data"`
	Errors map[string]string `json:"errors"`
}

func runServer() {
	http.HandleFunc("/add-task", func(w http.ResponseWriter, r *http.Request) {
		if !isTokenValid(r) {
			responseError(w, map[string]string{"token": "invalid token"})
			return
		}

		if r.Method == "POST" {
			channel := getChannel(r)
			var bodyRes []Queue
			parseBody(r, &bodyRes)
			var results = make([]interface{}, len(bodyRes))

			for _, t := range bodyRes {
				if t.Name == "" {
					responseError(w, map[string]string{"name": "name is required"})
					return
				}

				if t.Url == "" {
					responseError(w, map[string]string{"url": "url not defined for " + t.Name})
					return
				}
			}

			for i, t := range bodyRes {
				t.Parse()
				t.Headers = r.Header
				t = AddTask(channel, t)
				results[i] = map[string]string{"id": t.Id, "name": t.Name, "url": t.Url}
			}

			responseSuccess(w, results)
			fmt.Println("Tasks in progress:", tasksTotal[channel]-tasksComplete[channel], "Tasks Complete", tasksComplete[channel])
		}

	})

	printSuccess("server is running on " + lisenAddress())
	printInfo("token: " + *flagSecretToken)
	fmt.Println("To add new tasks send POST request to ")
	fmt.Println(lisenAddress() + "/add-task?token=" + *flagSecretToken)
	log.Fatal(http.ListenAndServe(*flagListen, nil))
	PrintMemUsage()
}

func isTokenValid(r *http.Request) bool {
	t := r.URL.Query().Get("token")
	if *flagSecretToken == t {
		return true
	}
	return false
}

func getChannel(r *http.Request) string {
	channel := r.URL.Query().Get("channel")
	if channel == "" {
		channel = "default"
	}
	//fmt.Fprintf(w, "channel is %s\n", channel)
	return channel
}

func parseBody(r *http.Request, body interface{}) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		addErrorLog(err.Error())
	}
	if err = json.Unmarshal(b, &body); err != nil {
		addErrorLog(err.Error())
	}
}

func tokenGenerator() string {
	b := make([]byte, 50)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func lisenAddress() string {
	url := "http://localhost" + *flagListen
	if string([]rune(*flagListen)[0]) != ":" {
		url = "http://" + *flagListen
	}
	return url
}

func responseSuccess(w http.ResponseWriter, data []interface{}) {
	res := Response{true, data, map[string]string{}}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func responseError(w http.ResponseWriter, errors map[string]string) {
	res := Response{false, make([]interface{}, 0), errors}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(js)
}
