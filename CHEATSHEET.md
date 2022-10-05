## This is a cheat sheet for this workshop

Basic go syntax 

### Create a new variable

Declare a variable and assign a value

```go
// define a new variable and assign a value
var name = "John"

// define a variable to but not assign a value
var id int
// or a shorthand way
number := 1
```

### Decoding JSON

```go
// define a struct to hold the JSON data
type Person struct {
    Name string
    Age  int
}

// define a variable to hold the JSON data
var p Person

// decode the JSON data into the variable
err := json.NewDecoder(r.Body).Decode(&p)
if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}
```

### Encoding JSON

```go
// define a struct to hold the JSON data
type Person struct {
    Name string
    Age  int
}

// define a variable to hold the JSON data
p := Person{
    Name: "John",
    Age:  42,
}

// encode the JSON data into the variable
err := json.NewEncoder(w).Encode(p)
if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}
```

### Generate random numbers

```go
// generate a random number between 0 and 100
var nr int
nr = rand.Intn(100)

// or 
nr := rand.Intn(100)
```

### Convert a int to a string

```go
// convert a int to a string
var nr int
nr = 42
var str string
str = strconv.Itoa(nr)
```

### Adding a new item to a slice

```go
// define a slice
var slice []string

// add a new item to the slice
slice = append(slice, "item")

// works for all types
var usrs []User
usrs = append(usrs, User{Name: "John", Age: 42})
```

### Creating a loop

```go
// create a loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// or if you want to loop over a slice
// where i is the index and item is the value
for i, item := range slice {
    fmt.Println(i, item)
}
```


### Removing an item from a slice

```go
// define a slice
var slice []string

// add a new item to the slice
slice = append(slice, "item")

// remove an item from the slice
slice = append(slice[:index], slice[index+1:]...)

// works for all types
var users []User{{Name: "John", Age: 42}, {Name: "Jane", Age: 42}}

// remove the user at index 1
users = append(users[:1], users[1+1:]...)
```

### Return a http message

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // set the content type
    w.Header().Set("Content-Type", "application/json")
    // write the message
    w.Write([]byte(`{"message": "Hello World"}`))

    // or
    fmt.Fprintf(w, `{"message": "Hello World"}`)

    // if you need to embed a variable in the message
    fmt.Fprintf(w, `{"message": "Hello %s"}`, name)

    // or
    fmt.Fprintf(w, `{"message": "Hello %s", "name": %s}`, name, name)
    
    // embeded variable types you might need
    // %s - string
    // %d - int
}
```

### Decode the body into an integer
    
```go
// define a variable to hold the JSON data
var id int

// decode the JSON data into the variable
err := json.NewDecoder(r.Body).Decode(&id)
if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}
```

### Get a query parameter

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // get the query parameter
    name := r.URL.Query().Get("name")
    // or
    name := r.FormValue("name")
}
```


### Opening a websocket connection

```go
conn, err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat", nil)
if err != nil {
    log.Fatal("dial:", err)
}
```

### Reading from the console
    
```go
// read line by line from the console
scanner := bufio.NewScanner(os.Stdin)

// buffer input until a new line is found
for scanner.Scan() {
    // get the text from the scanner
    text := scanner.Text()
    // do something with the text
}

```


### Writing to the websocket connection

```go
// write a message to the websocket connection
err := conn.WriteMessage(websocket.TextMessage, []byte("Hello World"))
if err != nil {
    log.Println("write:", err)
    return
}
```

### Reading from a websocket

```go
// t:     type of message
// msg:   message (in []byte format)
// error: error if any
t, msg, err := conn.ReadMessage()
if err != nil {
    log.Println("read:", err)
    return
}
text := string(msg)

// print the message to stdout (standard output)
fmt.Fprintf(os.Stdout, "%s\n", text)
```

### Switch statemnts

```go
// switch statements
switch t {
case websocket.TextMessage:
    // do something
case websocket.BinaryMessage:
    // do something
default:
    // do something
    // as we do not expect or handle this, this is a good place to return an error to our error channel
}
