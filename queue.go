package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Queue struct {
	Id            string
	Name          string                 `json:"name"`
	Data          map[string]interface{} `json:"data"`
	Url           string                 `json:"url"`
	Headers       http.Header
	CustomHeaders map[string]string `json:"headers"`
	Options       struct {
		// Request method. default is "POST"
		Method string
		// How many times repeat this task
		// `Repeat` of 0 means repeat forever. default is 1
		Repeat *int `json:"repeat,omitempty"`
		// How many of seconds you need to delay and execute the next time.
		// 60 means execute every 1 min. default is 0
		Every int64
		// Schedule the task to start executing at a specific time. default is 0
		StartAt int64
		// Set a time limit for waiting a response before cancelling the request.
		// i.e 30 means 30 sec, 0 means no limit default is 20
		Timeout *int `json:"timeout,omitempty"`
		// How many retry to execute the task when failed. defualt is 3
		Retry *int `json:"retry,omitempty"`
		// Number of seconds you need to delay and retry to executing the task when failed. defailt is 30
		RetryDelay *int `json:"retry_delay,omitempty"`
		// Api url will be called (POST request) when task failure after retries.
		FailureCallback *string `json:"failure_callback,omitempty"`
	} `json:"options"`
	LastExecute int64
	NextExecute int64
	CreatedAt   int64
	Repeat      int
	Tries       int
	Response    struct {
		Status string
		Body   string
	}
	//Failed bool
}

func (q *Queue) Parse() {
	if q.Options.Method == "" {
		q.Options.Method = "POST"
	}

	if repeat := 1; q.Options.Repeat == nil {
		q.Options.Repeat = &repeat
	}

	now := time.Now().Unix()
	if q.Options.StartAt == 0 {
		q.Options.StartAt = now
	}

	if q.Options.Timeout == nil {
		q.Options.Timeout = flagTimeout
	}

	if q.Options.Retry == nil {
		q.Options.Retry = flagRetry
	}

	if q.Options.RetryDelay == nil {
		q.Options.RetryDelay = flagRetryDelay
	}

	if q.Options.FailureCallback == nil {
		q.Options.FailureCallback = flagFailureCallback
	}

	q.CreatedAt = now
	q.NextExecute = q.Options.StartAt

}

func (q *Queue) mergeHeaders(header http.Header) {
	headers := make(http.Header, len(header))
	for k, vv := range header {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		headers[k] = vv2
	}
	if len(q.CustomHeaders) > 0 {
		for k, v := range q.CustomHeaders {
			headers.Set(k, v)
		}
	}

	q.Headers = headers
}

func (q *Queue) exec() bool {
	now := time.Now().Unix()
	if q.Tries > 0 {
		addErrorLog("---------------------------------------------")
		addErrorLog("Retry " + q.Name)
	}
	if q.sendRequest("execute") {
		//addCompleteLog(q)
		q.LastExecute = now
		//q.Failed = false
		q.setNextExecute(now)
		return true

	} else {
		// error
		if *q.Options.Retry == q.Tries {
			q.setNextExecute(now)
			// call flag_failure_callback
			if *q.Options.FailureCallback != "" {
				q.sendRequest("failure")
			}
		} else {
			q.NextExecute = func() int64 {
				return now + int64(*q.Options.RetryDelay)
			}()
		}
		q.Tries++

		return false
	}
}

func (q *Queue) IsExpired() bool {
	return q.NextExecute == 0
}

func (q *Queue) IsNeedToExecuteNow() bool {
	now := time.Now().Unix()
	return !q.IsExpired() && q.NextExecute <= now
}

func (q *Queue) sendRequest(reqType string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	jsonValue, err := json.Marshal(q.Data)

	if err != nil {
		log.Fatalln(err)
	}
	url := q.Url
	method := q.Options.Method
	if reqType == "failure" {
		url = *q.Options.FailureCallback
		method = "POST"
		jsonValue, _ = json.Marshal(q)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	req.Header = q.Headers
	client := &http.Client{
		Timeout: time.Duration(*q.Options.Timeout) * time.Second,
	}
	res, err := client.Do(req)
	if res != nil {
		defer func() {
			if err := res.Body.Close(); err != nil {
				q.Response.Body = err.Error()
				q.Response.Status = res.Status
				addErrorLog("---------------------------------------------")
				addErrorLog(err.Error())
			}
		}()

		body, _ := ioutil.ReadAll(res.Body)
		success := res.StatusCode < 299
		if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
			success = false
		}
		if !success {
			q.Response.Body = string(body)
			q.Response.Status = res.Status

			addErrorLog("---------------------------------------------")
			addErrorLog("failure " + q.Name)
			addErrorLog("Response Status " + res.Status)
			addErrorLog("Request : " + q.Options.Method + " " + q.Url)
			requestBody, _ := json.Marshal(q.Data)
			addErrorLog("Request Body : " + string(requestBody))
			addErrorLog("Response : " + string(body))
		}
		return success
	}
	if err != nil {
		q.Response.Body = err.Error()
		addErrorLog("---------------------------------------------")
		addErrorLog("failure " + q.Name)
		addErrorLog(err.Error())
		//panic(err)
	}
	return false
}

func (q *Queue) setNextExecute(now int64) {

	// if is one time execute or specific number
	if *q.Options.Repeat == 1 || *q.Options.Repeat > 0 && q.Repeat+1 == *q.Options.Repeat || *q.Options.Repeat == 0 && *q.Options.Retry == q.Tries {
		// delete
		q.NextExecute = 0
	} else {
		// if more than one time execute
		q.NextExecute = func() int64 {
			q.Repeat++
			return now + q.Options.Every
		}()
	}
	q.Tries = 0
}
