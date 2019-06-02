# TaskQ
Scheduling task queues for web service

## Features
- Scheduling the task by specific time
- Support all HTTP methods for web `GET, POST, PUT, PATCH, DELETE`
- Support channels
- Automatically remove task after complete
- Retrying execute task if failure
- Allow you repeat the task ever you want 
- Recovering uncompleted tasks when restart TaskQ

## Install
**Binary** :
> Download from [release](https://github.com/mindblowup/taskq/releases) page and download yours.   
Then run this
```bash
./taskq --secret="25b01e511b5fc414032692658a4c1362d8c702cde8759f7fde612fe3e6355c"
```

**From Source**
> You must have the `Go` environment installed
```bash
go get -u github.com/mindblowup/taskq
```

**For production**
> You can use pm2 or docker
```bash
pm2 start /path/to/taskq -- --secret="25b01e511b5fc414032692658a4c1362d8c702cde8759f7fde612fe3e6355c" --failure_callback="http://yourwebservice.com/api/v1/failure_callback"
```

## Configration
There no Configration file .
Run `./taskq -h` you will see the options like this 
```
Usage of ./taskq:
  -clear
        Clear all previous uncompleted tasks
  -error_log string
        error logs (default "./taskq_error.log")
  -failure_callback string
        Api url will be called (POST request) when task failure after retries.
  -listen string
        the http listen address (default ":8001")
  -retry int
        How many retry to execute the task when failed (default 3)
  -retry_delay int
        Number of seconds you need to delay and retry to executing the task when failed (default 30)
  -secret string
        Secret token because taskQ is end to end  (default "aa8fc052f014fd97942ddcc012760e4c67ac39cce325d6cb136f29ec10420a1ccc88e8b5d32608250c8770955dbbddb599ab")
  -timeout int
        How many times of seconds limit for waiting a response before cancelling the request. 0 means no limit (default 20)
```

## Usage
To add a task just send a `POST` request to **TaskQ** from anywhere 
You can add one task or more in the same request

```
POST http://localhost:8001/add-http-task?secret=25b01e511b5fc414032692658a4c1362d8c702cde8759f7fde612fe3e6355c
```
#### Body (Array of JSON)

| Name    | Type     | Required or optional |Description |
| --------|:--------:|:--------:|:--------|
| `name`  | `string` | `required` |name for task|
| `url` | `string` | `required` | full url for the webservice |
|`body` | `object`| `optional` | the body of request will be send to the webservice |
| `headers` | `object`| `optional` | the HTTP headers. may webservice required special headers like Authentication, Authorization ... |
| `options` | `object`| `optional` | see [Options](#options) for more details|



### PHP
```php
<?php
use TaskQ\TaskQ;
use TaskQ\HttpTask;
require './vendor/autoload.php';
$taskq = new TaskQ(':8001', '1f79ff70f7d2a26d4e1199b59ab8013d167298c02e5f2feb9910d21422a13e4a6ce86146df2b1968fc35542bac801469f66e');
$taskq->addHttpTask(function (HttpTask $tsk){
    $tsk->name('send_mail_to_user.10')
        ->method('POST')
        ->url('http://yourwebservice.com/api/v1/send-mail')
        ->data([
            'id' => 10,
            'email' => 'email@example.com',
            'type' => 'welcome'
        ]);
    return $tsk;
});

$taskq->addHttpTask(function (HttpTask $tsk){
    $tsk->name('Send SMS to user10')
        ->method('POST')
        ->url('http://yourwebservice.com/api/v1/send-sms')
        ->data([
            'id' => 10,
            'phone' => '+201234567890',
            'type' => 'welcome'
        ]);
    return $tsk;
});

$taskq->addHttpTask(function (HttpTask $tsk){
    $tsk->name('Update Payment Status for user10')
        ->method('PUT')
        ->url('https://yourwebservice.com/api/v1/update-payment-status/10')
        ->data([
            'id' => 10,
            'phone' => '+201234567890',
            'type' => 'welcome'
        ])->headers([
            // maybe you need to send headers to API you use.
            'Authorization' => 'Basic YWxhZGRpbjpvcGVuc2VzYW1l'
        ])->options(function (TaskOptions $opt){
            $opt->everyMonth(1) // execute every 1 month
              ->startAt(strtotime('tomorrow 2pm'))
              ->forever(); // repeat it forever
            return $opt;
        });
    return $tsk;
});
$response = $taskq->send();
if($taskq->hasErrors()){
    http_response_code(400);
    print_r($taskq->errors());
}
```


### Node.js
```javascript
var request = require("request");

var options = {
    method: 'POST',
    url: 'http://localhost:8001/add-http-task',
    qs: { token: '25b01e511b5fc414032692658a4c1362d8c702cde8759f7fde612fe3e6355c' },
    headers: {'Content-Type': 'application/json' },
    body: 
    [ 
        {
            // name for job (required)
            name: "send_mail_to_user.10",
            // full url for the webservice (required)
            url: "http://yourwebservice.com/api/v1/send-mail",
            //a body of request (optional)
            data: {
                id: "10",
                email: "user1@example.com",
                type: "welcome"
            }
        },
        {
            name: "Send SMS to user10",
            url: "http://yourwebservice.com/api/v1/send-sms",
            data: {
                id: "10",
                phone: "+201234567890",
                type: "welcome"
            }
        },
        {
            name: "Update Payment Status for user10",
            url: "https://yourwebservice.com/api/v1/update-payment-status/10",
            // you can add HTTP headers for every task 
            headers: {
                Authorization: "Basic YWxhZGRpbjpvcGVuc2VzYW1l"
            },
            options: {
                method: "PUT",
                repeat: 0, //forever
                every: 60*60*24*30, //every day
            }
        }
    ],
    json: true 
};

request(options, function (error, response, body) {
  if (error) throw new Error(error);

  console.log(body);
});

```

### Options
``` json
[
    {   
        "name": "Update Payment Status for user10",
        "url": "http://yourwebservice.com/api/v1/update-payment-status/10",
        "data": {},
        "options": {
            "method": "PUT",
            "repeat": 0,
            "every": 60*60*24,
            "startAt": 0,
            "retry": 5,
            "retry_delay": 10,
            "failure_callback": "http://yourwebservice.com/api/v1/failure_callback"
        }
    }
]
```

| Option             | Description                                                                   | Default  |
| -------------------|:------------------------------------------------------------------------------|:--------:|
| `method`           | Request method i.e **GET, HEAD, POST, PUT, DELETE** ... | POST     |
| `repeat`           | How many times repeat this task`repeat` of 0 means repeat forever             | 1        |
| `every`            | How many of seconds you need to delay and execute the next time. 60 means execute every 60 sec |    0 |
| `startAt`          | Schedule the task to start executing at a specific time.                      | 0 |
| `timout`           | Set a time limit for waiting a response before cancelling the request         | 30 |
| `retry`            | How many retry to execute the task when failed.                               | 3 |
| `retry_delay`      | Number of seconds you need to delay and retry to executing the task when failed | 30 |
| `failure_callback` | Url will be called (POST request) when task failure after retries.        | `""` |
