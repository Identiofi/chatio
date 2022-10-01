# Identio - Chatio
A server and client implementation of a simple chat over websockets


## Workshop 

This is part of a workshop hosted by Identio, where the goal is to familiarize the attendees 
with Go, websockets, connections and consluting work.

### Part 1
```
$ git checkout part-1
```

#### Introduction
In part 1 of this workshop, we will be creating a simple server that will listen for incoming
connections and print out the messages that are sent to it.

The server will be listening on port 8080, and will be expecting a message in the following format:

```
{
    "message": "Hello World"
}
```

The server will then print out the message to the console, and respond with the following message:

```
{
    "message": "Message received"
}
```

This part is already implemented, and you can run the server by running the following command:

`$ make run-server` or `$ cd server && go run main.go`

You can call the server using curl:

`$ curl -X POST -H "Content-Type: application/json" -d '{"message": "Hello World"}' http://localhost:8080`


<<<<<<< HEAD
---
Your task
---
=======
#### Your task

>>>>>>> caa2e3e (part-1: server)
Your task is to create a new enpoint that will return a list of active users in the chat.
The list of users can be stored in a global variable for now, and you can use the following
struct to represent a user:

```go
<<<<<<< HEAD
type user struct {
=======
type User struct {
>>>>>>> caa2e3e (part-1: server)
    Name string
    ID   string
}
```

The endpoint should be available at `/users` and should return a list of users in the following format:

```json
{
    "users": [
        {
            "name": "John Doe",
            "id": "1234"
        },
        {
            "name": "Jane Doe",
            "id": "5678"
        }
    ]
}
```

<<<<<<< HEAD
**You can test your endpoint by running the following command:**
=======
You can test your endpoint by running the following command:
>>>>>>> caa2e3e (part-1: server)

`$ curl -X GET http://localhost:8080/users`

---
<<<<<<< HEAD
Register
---
=======
>>>>>>> caa2e3e (part-1: server)

Your second task is to create a new endpoint that will allow a user to register to the chat.
The endpoint should be available at `/register` and should accept a request in the following format:

```json
{
<<<<<<< HEAD
    "name": "Mike Smith"
=======
    "name": "John Doe"
>>>>>>> caa2e3e (part-1: server)
}
```

The endpoint should then generate a unique ID for the user, and add the user to the list of active users.
The endpoint should then return a response in the following format:

```json
{
    "id": "1234"
}
```

<<<<<<< HEAD
Hint: You can use `rand.Intn()` to generate a random number. You'll have to convert your number to a string, using `strconv.Itoa()`

**Test your endpoint once again with curl, this time figure it out yourself.**

---
Unregister
---
=======
Test your endpoint once again with curl, this time figure it out yourself.

---

>>>>>>> caa2e3e (part-1: server)
Your third task is to create a new endpoint that will allow a user to unregister from the chat.
The endpoint should be available at `/unregister` and should accept a request in the following format:

```json
{
    "id": "1234"
}
```

The endpoint should then remove the user from the list of active users, and return a response in the following format:

```json
{
    "message": "User unregistered"
}
```
if the user was successfully removed, or

```json
{
    "message": "User not found"
}
```
if the user was not found.

**Test your endpoint once again with curl, you know how it's done by now.**

---
Chat
---

Your fourth and final task is to create a new endpoint that will allow a user to send a message to the chat. The endpoint should be available at `/chat` and should receive the user ID as a query parameter, like so: `/chat?id=1234`

The endpoint will connect the user to the chat, and will then listen for incoming messages from the user. The endpoint is a websocket endpoint. 

We will cheat a little bit here, and use our already implemented server to handle the websocket connection. Fetch the server from the `part-1-solution` branch, and run it with the following command:  

`git checkout part-1-solution -- server.go`  
> _Note: you'll need to be in the server directory for this to work_


Yuo should now have a new file called `server.go` in your server directory. Your task is to initialize
the chat server in the `main.go` file, and then connect to it from the `/chat` endpoint.

let's look at the `main.go` file:

```go
func main() {
    
    c := newChat()
    // notice the go keyword here, this is because we want to run the chat server 
    // in a separate goroutine. go routines are a way to run functions concurrently, 
    // in the background, like in a separate thread (but not really)
    go c.run()
    
    // .. your previous handlers

    //The new handler requires the chat struct as an argument, so we'll pass it in like this
    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        // .. your code goes here
    })

    // serve and listen
}
```

The chat struct has a `connectUser(usr user,  w http.ResponseWriter, r *http.Request)` method that you can use to connect a user to the chat. The method takes a user struct, a `http.ResponseWriter` and a `http.Request` as arguments.


furthermore, the connectUser method requires the request to be of type `http.MethodGet`, so you'll need to add a check for that in your handler. As the handler is not of type `http.MethodPost`, you cannot
send a request body to the endpoint. Instead, you'll have to send the user ID as a query parameter, like so: `/chat?id=1234`

to get the query parameters from the request, you can use the `r.URL.Query()` method, which returns a `map[string][]string`. You can then get the value of the `id` parameter by using the `Get` method on the map, like so: `r.URL.Query().Get("id")`

//  
To summarize, your task is to create a new handler for the `/chat` endpoint, and then connect the user to the chat using the `connectUser` method located in `server.go`. Before connecting the user, you'll have to check that the request method is `http.MethodGet`, and that the `id` query parameter is set, and that the user exists in the list of active users.


Once everything is done, you should be able to connect to the chat endpoint using curl, like so:

```$ curl -i --no-buffer -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" -H "Sec-WebSocket-Version: 13" http://localhost:8080/chat?id=1```

Curl will then wait for incoming messages from the server, you will not be able to send any messages using curl.. Let's fix that! Onto the client! 🚀


> Note: When running the server, you'll have to specify all the files that you want to run, like so:  
> `$ go run main.go server.go` or as a short hand (all files) `$ go run .`

---
Solution
---

Alright good job! You're done with part 1 of the workshop. If you want to see the solution, you can checkout the `part-1-solution` branch.

`$ git checkout part-1-solution`

### Part 2
```
$ git checkout part-2
```