# Identio - Chatio
A server and client implementation of a simple chat over websockets


## Workshop 

This is part of a workshop hosted by Identio, where the goal is to familiarize the attendees 
with Go, websockets, connections and consluting work. 

For more information: contact <a href="mailto:jimmy@identio.fi">Jimmy@identio.fi</a>


To get started, clone this repository

```
$ git clone github.com/identiofi/chatio
```

---
## Part 1
```
$ git checkout part-1
```

#### Introduction
In part 1 of this workshop, we will be creating a simple server that will listen for incoming
connections and print out the messages that are sent to it.

The server will be listening on port 8080, and will be expecting a message in the following format:

```json
{
    "message": "Hello World"
}
```

The server will then print out the message to the console, and respond with the following message:

```json
{
    "message": "Message received"
}
```

This part is already implemented, and you can run the server by running the following command:

`$ make run-server` or `$ cd server && go run main.go`

You can call the server using curl:

`$ curl -X POST -H "Content-Type: application/json" -d '{"message": "Hello World"}' http://localhost:8080`


---
Your task
---
Your task is to create a new enpoint that will return a list of active users in the chat.
The list of users can be stored in a global variable for now, and you can use the following
struct to represent a user:

```go
type user struct {
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
Register
---

Your second task is to create a new endpoint that will allow a user to register to the chat.
The endpoint should be available at `/register` and should accept a request in the following format:

```json
{
    "name": "Mike Smith"
}
```

The endpoint should then generate a unique ID for the user, and add the user to the list of active users.
The endpoint should then return a response in the following format:

```json
{
    "id": "1234"
}
```

Hint: You can use `rand.Intn()` to generate a random number.

Test your endpoint once again with curl, this time figure it out yourself.

---
Unregister
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
if the user was successfully removed, or

```json
{
    "message": "User not found"
}
```
if the user was not found.

Test your endpoint once again with curl, you know how it's done by now.

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

## Part 2

Great job! You're done with part 1 of the workshop. To summarize, you've created a simple server that can accept incoming connections, and you've created a few endpoints that allow users to register and unregister from the chat.

In part 2 of the workshop, we will be creating a simple CLI tool that will connect to the server, and allow users to send messages to the server.
For this part, we'll be using the `github.com/gorilla/websocket` package, which is a websocket client library for Go.

We also need a working chat server, if you did not complete part 1, you can checkout a working solution from the `part-1-solution` branch.

`$ git checkout part-1-solution -- server/main.go server/server.go`

> Note: you'll need to be in the root directory for this to work. If you want to save your changes please commit them first.
> You can also stash your changes, if you do not wish to commit them. With: `$ git stash`

Now that we have a working chat server, we can start building our client.  

For the duration of the part, leave the server running, and open a new terminal window.

---
Client
---

Your first task is to create a new file called `main.go` in the `client` directory (you should have this file).
In this file, you'll need to create a new struct called `client`, which will have the following fields:

```go
type client struct {
    // the websocket connection
    conn *websocket.Conn
    // errChan is used to send errors to the main goroutine
    // from the listen goroutine
    errChan chan error
}
```

You'll also need to create a new method for initializing the client `newClient`, which will take a `url` as an argument. The method should then connect to the chat server, and store the connection in the `conn` field. Your method should return a pointer to the newly created client.  

You will also have to `make()` the new `errChan` field.

```go
func newClient(url string) (*client, error) {
    // your code goes here
    c := &client{
        conn:    conn,
        errChan: make(chan error),
    }
    // your code goes here
}
```

---
Write
---

To send a message to the server, you can use the `WriteMessage` method on the `websocket.Conn` struct. 
The method takes no arguments, so you'll have to use the `bufio` package to read input from the console. 
You can then send the input to the server using the `WriteMessage` method.  
Create a new `Scanner` using the `bufio.NewScanner` function, and then use the `Scan` method to read 
input from the console. You can then get the input using the `Text` method on the `Scanner` struct.

```go
func (c *client) write() {
    // your code goes here
}
```
For the sake of debugging, you can use `fmt.Println` to print the input to the console, 
to make sure that you're reading the input correctly. 
> Hint: you can use `fmt.Fprintf(os.Stdout, ...)` to print to the console.

At this poit you should be able to send messages to the server, and see them printed to the console.
The server will also print a message to the console, which you can use to verify that your client is working correctly.

---
Read
---
Your client should then listen for incoming messages from the server, and print them to the console. 
You can use the `ReadMessage` method on the `websocket.Conn` struct to read incoming messages. 
`ReadMessage` returns a `messageType` and a `[]byte`, which you can then print to the console using `fmt.Println`.

under the hood `string` is just a slice of bytes, so you can easily convert `[]byte` to a string using the `string` function.


```go
func (c *client) listen() {
    // your code goes here
}
```

> Hint: you can use `fmt.Fprintf(os.Stdout, ...)` to print to the console.

---
Stiching it all together
---

Go requires a `main` function, so you'll need to create one, which should create a new `client` struct, 
and then start the write and listen goroutines.

```go
func main() {
    // your code goes here
}
```

To help you we have created a `run` method, which will run the client, start new goroutines for the reader and writer,
and handle any errors that might occur. We included a helper function called `fail` which will print an error message to the console,
and a close method which will close the connection and exit the program.

```go
// run starts the client
func (c *client) run() {
	// interrupt may be sent by the OS when the user presses Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// start listen and write in separate goroutines
	// allows us to read and write concurrently
	go c.listen()
	go c.write()

	// wait for an error or interrupt, this is done on the main thread
	// so that we can block until the program is terminated
	for {
		select {
		case err := <-c.errChan:
			c.Close()
			fail("error: %v\n", err)
		case <-interrupt:
			c.Close()
			return
		}
	}
}

func (c *client) Close() {
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func fail(msg string, o ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, o...)
	os.Exit(1)
}
```

---
---

To summarize, your task is to create a new struct called `client`, which will have a `conn` field 
of type `*websocket.Conn`. You'll also need to create a new method 
called `newClient`, which will take a `url` as an argument. 
The method should then connect to the chat server, and store the connection in the `conn` field.
You'll also need to create a new method called `send` on the `Client` struct, which will send a message to the server.
Finally, you'll need to create a new method called `listen` on the `Client` struct, which will listen 
for incoming messages from the server, and print them to the console.

Once everything is done, you should be able to run the client, and connect to the chat server.
You can do this by running the following command:

`$ go run .`

You should then be able to send messages to the server, and see them printed to the console.

---
Solution
---

Alright good job! You're done with part 2 of the workshop. If you want to see the solution, you can checkout the `part-2-solution` branch.

`$ git checkout part-2-solution`

### Bonus, and the fun part!

Now that we have a working chat server, and a working client, we can start chatting with each other! 🎉
We have setup a public chat server, which you can connect to using the following url: `ws://chat.identio.fi`

**Here are some things you can try:**

Medium difficulty:
- Use the kingpin library to parse command line arguments for the client
- Allow the user to specify the ID before connecting to the server (not specifing the ID through the URL)

Hard difficulty:
- Modify the server to allow multiple rooms
- Modify the client to allow the user to specify the room they want to connect to
  

---
Conclusion
---

Great job! You've now created a simple chat server, and a simple chat client.
You can now use this as a starting point for your own projects, and build something awesome!
Thanks for attending the workshop, and I hope you had fun!