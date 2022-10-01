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


#### Your task

Your task is to create a new enpoint that will return a list of active users in the chat.
The list of users can be stored in a global variable for now, and you can use the following
struct to represent a user:

```go
type User struct {
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

You can test your endpoint by running the following command:

`$ curl -X GET http://localhost:8080/users`

---

Your second task is to create a new endpoint that will allow a user to register to the chat.
The endpoint should be available at `/register` and should accept a request in the following format:

```json
{
    "name": "John Doe"
}
```

The endpoint should then generate a unique ID for the user, and add the user to the list of active users.
The endpoint should then return a response in the following format:

```json
{
    "id": "1234"
}
```

Test your endpoint once again with curl, this time figure it out yourself.

---

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

Test your endpoint once again with curl, you know how it's done by now.


Alright, you're done with part 1 of the workshop. If you want to see the solution, you can checkout the `part-1-solution` branch.

`$ git checkout part-1-solution`

### Part 2
```
$ git checkout part-2
```