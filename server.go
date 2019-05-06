package main

import (
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

type ListResponse struct {
	Status bool              `json:"status"`
	Data   Channels          `json:"data"`
	Errors map[string]string `json:"errors"`
}
type ListChannel struct {
	Status bool              `json:"status"`
	Data   Tasks             `json:"data"`
	Errors map[string]string `json:"errors"`
}

func runServer() {
	routes()

	printSuccess("server is running on " + listenAddress())
	printInfo("secret: " + *flagSecretToken)
	//fmt.Println("To add new tasks send POST request to ")
	//fmt.Println(listenAddress() + "/add-http-task?token=" + *flagSecretToken)
	log.Fatal(http.ListenAndServe(*flagListen, nil))
}

func getChannel(r *http.Request) string {
	channel := r.URL.Query().Get("channel")
	if channel == "" {
		channel = "default"
	}
	return channel
}

func parseBody(r *http.Request, body interface{}) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		addErrorLog(err.Error())
	}
	if err = json.Unmarshal(b, &body); err != nil {
		fmt.Println(r.Body)
		addErrorLog(err.Error())
	}
}

func listenAddress() string {
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
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
