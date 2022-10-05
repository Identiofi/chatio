# Identio - Chatio
A server and client implementation of a simple chat over websockets


## Workshop 

This is part of a workshop hosted by Identio, where the goal is to familiarize the attendees 
with Go, websockets, connections and consluting work. 

For more information: contact <a href="mailto:jimmy@identio.fi">Jimmy@identio.fi</a>


To get started, clone this repository

```
$ git clone https://github.com/Identiofi/chatio
```

## Prerequisits

1. Install go
2. Check that your install works

If running Visual Studio Code, install the Go plugin  
[go plugin for VSCode](https://marketplace.visualstudio.com/items?itemName=golang.go)

---
## Part 1
```
$ git checkout part-1
```

#### Introduction
In part 1 of this workshop, we will be creating a simple REST API that will be used to create and list users.  
This is a very simple API, and the same base will be used in the next part of the workshop.

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

```
$ cd server/
$ go run main.go
```

You can call the server using curl:

`$ curl -X POST -H "Content-Type: application/json" -d '{"message": "Hello World"}' http://localhost:8080`


---
Your task
---
Your task is to create a new enpoint that will return a list of registered users.
As of now, define a global variable `userList` ([see row 49](server/main.go)) that will hold a list of users.

```go
type User struct {
    Name string `json:"name"`
    ID   int    `json:"id"`
}
```

The endpoint should be available at `/users` and should return a list of users in the following format:

```json
{
    "users": [
        {
            "name": "John Doe",
            "id": 1234
        },
        {
            "name": "Jane Doe",
            "id": 5678
        }
    ]
}
```


// To summarize, you should create a new handler for the endpoint `/users` that will return a list of users,  
The list of users is stored in a global variable `userList` that is defined in the main method.


You can test your endpoint by running the following command:

`$ curl -X GET http://localhost:8080/users`

---
Register
---

Your second task is to create a new endpoint that will allow a user to register.  
The endpoint should be available at `/register` and should accept a request in the following format:

```json
{
    "name": "Mike Smith"
}
```

> Hint: look att the helloWorldhandler, which exepcts a format of `{ "message": "Hello world" }`

To easily decode the incoming message into the right type, use the `User` type defined in the previous task.  
(encoding and decoding json in go is very easy, and does not require all fields to be present in the message)


The handler should then generate a unique ID for the user, and add the user to the list of active users.

The endpoint should then return a response in the following format:

```json
{
    "id": 1234
}
```

> Hint: look at the CHEATSHEET for help with generating a unique ID and the `append` function for adding a user to the list of users.  
> CHEATSheet also contains information about how to return a response in the correct format.

Test your endpoint once again with curl
    
    `$ curl -X POST -H "Content-Type: application/json" -d '{"name": "Mike Smith"}' http://localhost:8080/register`

Feel free to test your endpoint with multiple users, and see if the list of users is updated correctly.

---
Unregister
---
Your third task is to create a new endpoint that will allow a user to unregister from the chat.
The endpoint should be available at `/unregister` and should accept a request in the following format:

```json
{
    "id": 1234
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

you will have to loop over your `userList` and find a user with the given ID. 
If there is no user with the given ID, return the error, else remove the user from the list.


Test your endpoint once again with curl,

`$ curl -X POST -H "Content-Type: application/json" -d '{"id": 1234}' http://localhost:8080/unregister`

---
Chat
---
Your fourth and final task is to create a new endpoint that will allow a user to send a message to the chat. The endpoint should be available at `/chat` and should receive the user ID as a query parameter, like this: `/chat?id=1234`

You can do this task in multiple steps, as it will be a bit more complex than the previous tasks.

1. Create a new endpoint `/chat` and define a handler for the endpoint.
2. The handler should accept a query parameter `id` that will contain the ID of the user sending the message.
   1. Have a look at the CHEATSHEET for help with parsing query parameters.
   2. Print the ID to the console to make sure that the ID is received correctly. (you can use the `fmt.Println` function for this)
3. Find the given user in the list of users, and print the user to the console to make sure that the user is found correctly.
4. Connect the user to the chat. (see the next section for more information)

---
We will cheat a little bit here, and use our already implemented server to handle the websocket connection. Fetch the server from the `part-1-solution` branch, you can get it with the following command:  

`git checkout origin/part-1-solution -- server/server.go`  
> _Note: you'll need to be in the project root directory for this to work_


Yuo should now have a new file called `server.go` in your server directory. Your task is to initialize
the chat server in the `main.go` file, and then connect to it from the `/chat` endpoint.

let's look at the `main.go` file below:

```go
// This is your main file (inside main.go) that you have worked on.
func main() {
    // create a new chat, c
    c := newChat()
    // notice the go keyword here, this is because we want to run the chat server 
    // in a separate goroutine. go routines are a way to run functions concurrently, 
    // in the background, like in a separate thread (but not really)
    go c.run()
    
    // .. your previous handlers

    //The new handler requires the chat struct as an argument, so we'll pass it in like this
    http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
        // .. your code goes here
    })

    // serve and listen
}
```

Note that we are initializing the handler for the `/chat` endpoint inside the `main` function. This is done because we need to pass the chat struct as an argument to the handler. Feel free to move your handler you just wrote to the `main` function.

The chat struct has a `connectUser(User, http.ResponseWriter, *http.Request)` method that you can use to connect a user to the chat.

//  To summarize, your task is to create a new handler for the `/chat` endpoint, and then connect the user to the chat using the `connectUser` method located in `server.go`. Before connecting to the chat, get the user from the `id` query parameter, check that the user exists in the list of active users.


> Note: When running the server, you'll have to specify all the files that you want to run, like so:  
> `$ go run main.go server.go` or as a short hand (all files) `$ go run .`


Once everything is done, you should be able to connect to the chat endpoint using curl, like so:

```$ curl -i --no-buffer -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" -H "Sec-WebSocket-Version: 13" http://localhost:8080/chat?id=1```

Curl will then wait for incoming messages from the server, you will not be able to send any messages using curl.. Let's fix that! Onto the client! 🚀

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

`$ git checkout origin/part-1-solution -- server/main.go server/server.go`

> Note: you'll need to be in the root directory for this to work. If you want to save your changes please commit them first.
> You can also stash your changes, if you do not wish to commit them. With: `$ git stash`

Now that we have a working chat server, we can start building our client.  

For the duration of the part, leave the server running, and open a new terminal window.

---
Client
---

Your first task is to create a new file called `main.go` in the `client` directory (you should have this file).

Once again, cd into the `client` directory, so that you can run the client from the same directory.

`$ cd client`

In this file, you'll need to create a new struct called `Client`, which will have the following fields:

```go
type Client struct {
    // the websocket connection
    conn *websocket.Conn
    // errChan is used to send errors to the main goroutine
    // from other goroutine
    errChan chan error
}
```

```
A goroutine is a lightweight thread managed by the Go runtime.
```

You'll also need to create a new method for initializing the client `newClient`, which should take a `url` as an argument, of type string. The method should then connect to the chat server, and store the connection in the `conn` field. Your method should return a pointer to the newly created client.  

You will also have to `make()` the new `errChan` field.

```go
func newClient(url string) (*Client, error) {
    // your code goes here (initialize connection)
    // handle errors

    c := &Client{
        conn:    conn,
        errChan: make(chan error),
    }

    fmt.Printf("Connected to server: %s\n", url)
    // your code goes here
    // return client
}
```

---
Writer
---

To send a message to the server, you can use the `WriteMessage` method on the `websocket.Conn` struct. The method takes no arguments, so you'll have to use the `bufio` package to read input from the console. Once again, let's do this in a few steps.

First, look for the placeholder method `write` in the `client.go` file, and implement it.  

```go
func (c *client) write() {
    // your code goes here
}
```

1. The `write` method should read input from the console, and for now send it back to the console. You can use the `fmt.Println` function to print to the console.

2. Next, you'll need to create a new goroutine that will run the `write` method. You can do this by calling the `go` keyword, like so: `go c.write()`.

3. Finally, you'll need to call the `write` method from the `main` method, like so: `c.write()`.

```go

func main() {
    // your code goes here

    // call the write method
    go c.write()

    // your code goes here
}
```

Test your code by running the client, and write a message in the console. You should see the message printed back to the console.

4. Change the `write` method to send the message to the server, instead of printing it to the console. You can use the `WriteMessage` method on the `websocket.Conn` struct. The method takes no arguments, so you'll have to use the `bufio` package to read input from the console.

> Hint: check the CHEATSHEET for how to read an input.

At this poit you should be able to send messages to the server, and see them printed to the console.
The server will also print a message to the console, which you can use to verify that your client is working correctly.

---
Read
---
Your client should then listen for incoming messages from the server, and print them to the console. 
You can use the `ReadMessage` method on the `websocket.Conn` struct to read incoming messages. 
`ReadMessage` returns a `messageType` and a `[]byte`, which you can then print to the console using `fmt.Println`.

1. Create a new method listen (already defined)

```go
func (c *Client) listen() {
    // your code goes here
}
```

2. Create an infinite loop, and call the `ReadMessage` method. Handle any error
3. Write a switch statement on the type returned from the `ReadMessage` method
   > Hint: Look for help in the CHEATSHEET
4. define a default statment in the switch statement, which writes an error to the error channel
   1. Return an error to the `errChan` field of the client, as this is an unknown message type
5. Print the message to the console


under the hood `string` is just a slice of bytes, so you can easily convert `[]byte` to a string using the `string` function.

> Hint: you can use `fmt.Fprintf(os.Stdout, ...)` to print to the console.

---
Stiching it all together
---

We should now have a way to read and write messages from the server. Let's put it all together.

We have helped you by providing a `run` method, but feel free to do it your own way. 
In the `main` method, do the following:
1. Initialize the client (using the `newClient` method)
2. Start the Writer in a new goroutine
3. Start the Listener in a new goroutine
4. if an error is written to the `errChan` field, print it to the console and exit the program
    - write an infinite loop on the main gorouting that waits for an error to be written to the `errChan` field
    - you can use the `os.Exit` function to exit the program
    - handle if the user press ctrl + c to exit the program (hint: check Run method)

```go
func main() {
    // your code goes here
}
```

As the final task, you should gracefully close the connection if an error occures or the user interupts the program, for this you can use the `Close` method which we provided.

Bonus: and handle any errors that might occur. We included a helper function called `fail` which will print an error message to the console,

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

To summarize, your task is to create a new struct called `Client`, which will have a `conn` field 
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
We have setup a public chat server, which you can connect to using the following url: `ws://chat.fly.dev/chat?id=X`

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
