# MARC'S GO REST API Demo
## Demonstrating REST API using Go language

## HOW TO TEST

### Client Tool:
Use [HTTPie](https://httpie.org/).

### Code the main.go program to import and use the MarcGoRESTAPIDemo Go package.
```go
/*
 * Copyright (c) 2021.
 * Marc Concepcion
 * marcanthonyconcepcion@gmail.com
 */

package main

import (
	"MarcGoRESTAPIDemo"
)

func main() {
	compiler := MarcGoRESTAPIDemo.MakeSubscriberController(MarcGoRESTAPIDemo.MakeDatabaseRecords())
	compiler.ViewHandleRequests()
}
```

### Start the REST API web server
```
> go run main.go
```

### SQL Script to create the test database scheme
[CreateSubscribersDatabase.sql](resources/CreateSubsribersDatabase.sql)

## FUNCTIONAL TEST SAMPLES

### Requirement 1: Create a new subscriber user record

#### Demonstrates POST without ID and CREATE a specified single record
```
C:\>http post http://127.0.0.1:8080/subscribers?email_address=riseofskywalker@starwars.com"&"last_name=Palpatine"&"first_name=Rey
HTTP/1.1 200 OK
Content-Length: 133
Content-Type: text/plain; charset=utf-8

{
    "message": "Record created",
    "updates": {
        "email_address": "riseofskywalker@starwars.com",
        "first_name": "Rey",
        "last_name": "Palpatine"
    }
}
```

### Requirement 2-1: Fetch a subscriber user record

#### Demonstrates GET with ID and RETRIEVE a specified single record
```
C:\>http get http://127.0.0.1:8080/subscribers/1
HTTP/1.1 200 OK
Content-Length: 101
Content-Type: text/plain; charset=utf-8

{
    "email_address": "riseofskywalker@starwars.com",
    "first_name": "Rey",
    "index": 1,
    "last_name": "Palpatine"
}
```

### Requirement 2-2: Fetch all subscriber user records

#### Demonstrates GET without ID and RETRIEVE all records
```
C:\>http get http://127.0.0.1:8080/subscribers
HTTP/1.1 200 OK
Content-Length: 377
Content-Type: text/plain; charset=utf-8

[
    {
        "email_address": "riseofskywalker@starwars.com",
        "first_name": "Rey",
        "index": 1,
        "last_name": "Palpatine"
    },
    {
        "email_address": "marcanthonyconcepcion@gmail.com",
        "first_name": "Marc Anthony",
        "index": 2,
        "last_name": "Concepcion"
    },
    {
        "email_address": "marcanthonyconcepcion@email.com",
        "index": 3
    },
    {
        "email_address": "kevin.andrews@email.com",
        "first_name": "Kevin",
        "index": 4,
        "last_name": "Andrews"
    }
]
```

If there are no records in the database, the API shall return an *HTTP 204: No Content* status code.
```
C:\>http get http://127.0.0.1:8080/subscribers
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: text/plain; charset=utf-8

[]
```

### Requirement 3: Edit an existing subscriber user record

#### Demonstrates PUT with ID and UPDATE a specified single record
```
C:\>http put http://127.0.0.1:8080/subscribers/1?last_name=Skywalker
HTTP/1.1 200 OK
Content-Length: 74
Content-Type: text/plain; charset=utf-8

{
    "message": "Record updated",
    "updates": {
        "index": 1,
        "last_name": "Skywalker"
    }
}
```

### Requirement 4: Activate a subscriber user record.

#### Demonstrates PATCH with ID to update a subscriber record while remaining idempotent.
```
C:>http patch http://127.0.0.1:8080/subscribers/1?activation_flag=true
HTTP/1.1 200 OK
Content-Length: 52
Content-Type: text/plain; charset=utf-8

{
    "details": "Record #1 activated.",
    "status": "success"
}
```

### Requirement 5: Delete an existing subscriber user record

#### Demonstrates DELETE with ID and DELETE a specified single record
```
C:\>http delete http://127.0.0.1:8080/subscribers/1
HTTP/1.1 200 OK
Content-Length: 64
Content-Type: text/plain; charset=utf-8

{
    "details": "Deleted record of subscriber #1.",
    "status": "success"
}
```

### Error Test Case 1: Get a record of a subscriber who does not exist.
```
C:\>http get http://127.0.0.1:8080/subscribers/400
HTTP/1.1 404 Not Found
Content-Length: 57
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

{
    "details": "Subscriber does not exist.",
    "status": "error"
}
```

### Error Test Case 2: Call an API without the prescribed 'subscribers' model
```
C:>http get http://127.0.0.1:8080
HTTP/1.1 404 Not Found
Content-Length: 19
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

404 page not found
```

### Error Test Case 3: Call an API with a model that is not 'subscribers'
```
C:\>http get http://127.0.0.1:8080/notsubscribers/1/
HTTP/1.1 404 Not Found
Content-Length: 19
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

404 page not found
```

### Error Test Case 4: Call HTTP commands that are not being used by the API.
```
C:\>http trace http://127.0.0.1:8080/subscribers
HTTP/1.1 405 Method Not Allowed
Content-Length: 0
```

### Error Test Case 5: POST an already existing record
```
C:\>http post http://127.0.0.1:8080/subscribers?email_address=riseofskywalker@starwars.com"&"last_name=Palpatine"&"first_name=Rey
HTTP/1.1 500 Internal Server Error
Content-Length: 126
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

{
    "details": "Error 1062: Duplicate entry 'riseofskywalker@starwars.com' for key 'subscribers.email_address'",
    "status": "error"
}
```

### Error Test Case 5-1: POST with specified ID.
```
C:\>http post http://127.0.0.1:8080/subscribers/1?email_address=riseofskywalker@starwars.com"&"last_name=Palpatine"&"first_name=Rey
HTTP/1.1 405 Method Not Allowed
Content-Length: 0
```

### Error Test Case 5-2: POST without required parameters
```
C:\>http post http://127.0.0.1:8080/subscribers
HTTP/1.1 405 Method Not Allowed
Content-Length: 137
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

{
    "details": "HTTP command POST without providing parameters is not allowed. Please provide an acceptable HTTP command.",
    "status": "error"
}
```

### Error Test Case 5-3: PUT without specified ID.
```
C:\>http put http://127.0.0.1:8080/subscribers?last_name=Skywalker
HTTP/1.1 405 Method Not Allowed
Content-Length: 0
```

### Error Test Case 5-4: PUT without required parameters
```
C:\>http put http://127.0.0.1:8080/subscribers/1
HTTP/1.1 405 Method Not Allowed
Content-Length: 136
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

{
    "details": "HTTP command PUT without providing parameters is not allowed. Please provide an acceptable HTTP command.",
    "status": "error"
}
```

### Error Test Case 6: DELETE without specified ID
```
C:\>http delete http://127.0.0.1:8080/subscribers
HTTP/1.1 405 Method Not Allowed
Content-Length: 0
```

### Error Test Case 7: PATCH that does not activate the Subscriber user record.
```
C:\>http patch http://127.0.0.1:8080/subscribers/1?activation_flag=false
HTTP/1.1 400 Bad Request
Content-Length: 114
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

{
    "details": "Only activating a subscriber is allowed. Please set the activation_flag to 'true'.",
    "status": "error"
}
```

For more inquiries, please feel free to e-mail me at marcanthonyconcepcion@gmail.com.

Thank you.

:copyright: 2021 Marc Concepcion

# GO LANG NG GO!
