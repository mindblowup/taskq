# TaskQ
Asynchronously task queues over HTTP request

## Usage
```
curl -X POST \
    'http://localhost:8001/add-http-task?token=25b01e511b5fc414032692658a4c1362d8c702cde8759f7fde612fe3e6355cdbc5b9b5a185c38f500083c5721ff59ae77718' \
    -H 'Content-Type: application/json' \
    -d '[
        {
            "name": "Send email to user10",
            "url": "http://yourwebservice.com/api/v1/send-mail",
            "data": {
                "id": "10",
                "email": "user1@example.com",
                "type": "welcome"
            }
        },
        {
            "name": "Send SMS to user10",
            "url": "http://yourwebservice.com/api/v1/send-sms",
            "data": {
                "id": "10",
                "phone": "+201234567890",
                "type": "welcome"
            }
        }
    ]'
```

```php
<?php

$request = new HttpRequest();
$request->setUrl('http://localhost:8001/add-http-task');
$request->setMethod(HTTP_METH_POST);

$request->setQueryData(array(
    'token' => '25b01e511b5fc414032692658a4c1362d8c702cde8759f7fde612fe3e6355cdbc5b9b5a185c38f500083c5721ff59ae77718'
));

$request->setHeaders(array(
    'Content-Type' => 'application/json'
));

$request->setBody('[
    {
        "name": "Send email to user10",
        "url": "http://yourwebservice.com/api/v1/send-mail",
        "data": {
            "id": "10",
            "email": "user1@example.com",
            "type": "welcome"
        }
    },
    {
        "name": "Send SMS to user10",
        "url": "http://yourwebservice.com/api/v1/send-sms",
        "data": {
            "id": "10",
            "phone": "+201234567890",
            "type": "welcome"
        }
    }
]');

try {
    $response = $request->send();

    echo $response->getBody();
} catch (HttpException $ex) {
    echo $ex;
}

```

```javascript
var request = require("request");

var options = {
    method: 'POST',
    url: 'http://localhost:8001/add-http-task',
    qs: { token: '25b01e511b5fc414032692658a4c1362d8c702cde8759f7fde612fe3e6355cdbc5b9b5a185c38f500083c5721ff59ae77718' },
    headers: {'Content-Type': 'application/json' },
    body: 
    [ 
        {
            name: "Send email to user10",
            url: "http://yourwebservice.com/api/v1/send-mail",
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
```
[
    {
        "name": "Update Payment Status for user10",
        "url": "http://yourwebservice.com/api/v1/update-payment-status/10",
        "options": {
            "method": "PUT",
            "repeat": 0,
            "every": 60*60*24,
            "startAt": 0,
            "retry" : 5,
            "retry_delay": 10
        }
    },
    {
        "name": "Send SMS to user1",
        "url": "http://yourwebservice.com/api/v1/send-sms",
        "data": {
            "id": "10",
            "phone": "+201234567890",
            "type": "welcome"
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