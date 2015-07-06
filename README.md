Humble/Rest
=============

[![GoDoc](https://godoc.org/github.com/go-humble/rest?status.svg)](https://godoc.org/github.com/go-humble/rest)

Version X.X.X

A small go package for sending requests to a RESTful API and unmarshaling the response.
Rest sends requests using CRUD semantics. It supports requests with a Content-Type
of either application/x-www-form-urlencoded or application/json and parses json responses
from the server. Rest works great as a stand-alone package or in combination with other
packages in the [Humble Framework](https://github.com/go-humble/humble).

Rest is written in pure go. It feels like go, follows go idioms when possible, and
compiles with the go tools. But it is meant to be compiled to javascript and run
in the browser.


Browser Support
---------------

Rest works with IE9+ (with a
[polyfill for typed arrays](https://github.com/inexorabletash/polyfill/blob/master/typedarray.js))
and all other modern browsers. Rest compiles to javascript via [gopherjs](https://github.com/gopherjs/gopherjs)
and this is a gopherjs limitation.

Rest is regularly tested with the latest versions of Firefox, Chrome, and Safari on Mac OS.
Each major or minor release is tested with IE9+ and the latest versions of Firefox and Chrome
on Windows.


Installation
------------

Install rest like you would any other go package:

```bash
go get github.com/go-humble/rest
```

You will also need to install gopherjs if you don't already have it. The latest version is
recommended. Install gopherjs with:

```
go get -u github.com/gopherjs/gopherjs
```


Quickstart Guide
----------------

Rest follows CRUD semantics, so there are methods for
[`Create`](https://godoc.org/github.com/go-humble/rest/#Client.Create),
[`Read`](https://godoc.org/github.com/go-humble/rest/#Client.Read),
[`Update`](https://godoc.org/github.com/go-humble/rest/#Client.Update), and
[`Delete`](https://godoc.org/github.com/go-humble/rest/#Client.Delete).
In addition, there is a method for getting all of a specific type of resources, which we
call [`ReadAll`](https://godoc.org/github.com/go-humble/rest/#Client.ReadAll).

### The Model Interface

The `Create`, `Read`, `Update`, `Delete`, and `ReadAll` methods all expect `Model`
(or a slice of `Model`) as an argument.
[`Model`](https://godoc.org/github.com/go-humble/rest/#Model) is an interface
which all your models should implement. It consists of only two methods:

``` go
type Model interface {
	// ModelId returns a unique identifier for the model. It is used for determining
	// which URL to send a request to.
	ModelId() string
	// RootURL returns the url for the REST resource corresponding to this model.
	// If you want to send requests to the same server, it should look something
	// like "/todos". If you want to send requests to a different server, you can
	// include the entire domain in the url, e.g. "http://example.com/todos". Note
	// that the trailing slash should not be included.
	RootURL() string
}
```

Because of the way reflection is used to encode the data, a Model must have an
underlying type of either a struct or a pointer to a struct, and any fields you want
to be included in requests/responses need to be exported.

If you like, you can embed [`DefaultId`](https://godoc.org/github.com/go-humble/rest/#DefaultId)
to give your models an `Id` property and a `ModelId` method which simply returns it.

Here's a full example of a Todo type which implements `Model`:

``` go
type Todo struct {
	Title        string
	IsCompleted  bool
	rest.DefaultId
}

func (t *Todo) RootURL() {
	"/todos"
}
```

### Instantiating a Client

Before sending any requests, you need to instantiate a new client. Typically, you will only need
to do this once per application. It's sometimes a good idea to make this a top-level variable.

``` go
var client = rest.NewClient()
```

You can set the `ContentType` property of the client to change the Content-Type header of all
requests sent. By default, all new clients have a `ContentType` of `ContentURLEncoded`, which
corresponds to the header value "application/x-www-form-urlencoded". Alternatively, you can
set the `ContentType` to `ContentJSON` in order to send requests with the Content-Type header
"application/json". Rest will automatically encode the body of requests based on the header.

```go
// Create a new client which uses json instead of url encoding
var client = &rest.Client{
	ContentType: ContentJSON,
}
``` 

### Create

The [`Create`](https://godoc.org/github.com/go-humble/rest/#Client.Create) method sends a POST
request to the model's root url. It is used to create new models on the server. Create will
mutate the given model, setting its fields based on the response from the server.

``` go
todo := &Todo{
	Title: "Discover the meaning of life",
}
if err := client.Create(todo); err != nil {
	// Handle err
}
```

If you are using url encoding, the request sent to the server in the above example would look
like this:

```
POST /todos

Title=Discover%20the%20meaning%20of%20life
```

And the response from the server should look like this:

``` json
{
	"Id": "9fjq293n8fw8",
	"Title": "Discover the meaning of life",
	"IsCompleted": false
}
```

### Read

The [`Read`](https://godoc.org/github.com/go-humble/rest/#Client.Read) method sends a GET
request to the model's root url with the given id appended. It is used to get existing models
from the server. Read will mutate the given model, setting its fields based on the response
from the server.

``` go
// Create a new todo which is initially empty.
todo := Todo{}
// Call Read to get it from the server and fill in
// all the fields of todo.
if err := client.Read("9fjq293n8fw8", &todo); err != nil {
	// Handle err
}
fmt.Println(todo.Title)
// Output:
//  "Discover the meaning of life"
```

The example above would send a GET request to "/todos/9fjq293n8fw8". The response from
the server should look like this:

``` json
{
	"Id": "9fjq293n8fw8",
	"Title": "Discover the meaning of life",
	"IsCompleted": false
}
```

### Update

The [`Update`](https://godoc.org/github.com/go-humble/rest/#Client.Update) method sends a PATCH
request to the model's root url with the model id appended. It is used to update existing models
on the server. Update will mutate the given model, setting its fields based on the response from
the server.

``` go
// Assume todo is an existing model with a valid id.
// Here's how we would mark it as completed on the
// server.
todo.IsCompleted = true
if err := client.Update(todo); err != nil {
	// Handle err
}
```

If you are using url encoding, the request sent to the server in the above example would look
like this:

```
PATCH /todos/9fjq293n8fw8

IsCompleted=true
```

And the response from the server should look like this:

``` json
{
	"Id": "9fjq293n8fw8",
	"Title": "Discover the meaning of life",
	"IsCompleted": true
}
```

### Delete

The [`Delete`](https://godoc.org/github.com/go-humble/rest/#Client.Update) method sends a DELETE
request to the model's root url with the model id appended. It is used to delete existing models
from the server. Delete does not do anything with the response from the server and the given model
is not mutated.

``` go
// Assume todo is an existing model with a valid id.
if err := client.Delete(todo); err != nil {
	// Handle err
}
```

The example above would send a DELETE request to "/todos/9fjq293n8fw8". The response from
the server is not used by the rest package, but it should typically either be an empty json
response:

``` json
{}
```

Or the model that was deleted:

``` json
{
	"Id": "9fjq293n8fw8",
	"Title": "Discover the meaning of life",
	"IsCompleted": true
}
```

### ReadAll

The [`ReadAll`](https://godoc.org/github.com/go-humble/rest/#Client.ReadAll)
method sends a GET request to the model's root url. It is used to get all
existing models of a given type from the server. `ReadAll` accepts an
`interface{}` as an argument but will return an error if it is not a pointer to
a slice of some type that implements `Model`. (ReadAll does not accept
`[]rest.Model` because that would require an additional, non-free cast. In go,
a slice of some interface is not the same as a slice of some concrete type that
implements the interface, and you cannot convert freely from one to the other.)
`ReadAll` will mutate the given slice based on the response from the server,
growing or shrinking it as needed.

``` go
todos := []Todo{}
if err := client.ReadAll(&todos); err != nil {
	// Handle err
}
```

The example above would send a GET request to "/todos". The response from
the server should look like this:

``` json
[
	{
		"Id": "9fjq293n8fw8",
		"Title": "Discover the meaning of life",
		"IsCompleted": true
	},
	{
		"Id": "af948j2f0jdh",
		"Title": "Read a good book",
		"IsCompleted": false
	}
]
```

### Handling Errors

Whenever a non-2xx status code is returned by the server, the `Create`,
`Read`, `Update`, `Delete`, and `ReadAll` methods will return an
[`HTTPError`](https://godoc.org/github.com/go-humble/rest/#HTTPError),
which contains pretty much everything you need to know from the response.

``` go
type HTTPError struct {
	// URL is the url that the request was sent to
	URL string
	// Body is the body of the response
	Body []byte
	// StatusCode is the http status code of the response
	StatusCode int
}
```

It is up to the caller to decide what to do based on the type of error
that occurred. A typical flow might look something like this:

``` go
todo := Todo{}
if err := client.Read("9fjq293n8fw8", &todo); err != nil {
	// Check if the error was of type rest.HTTPError and if so,
	// if the status code was 404. This can be a pretty common
	// error and should be handled appropriately.
	if httpErr, ok := err.(rest.HTTPError); ok && httpErr.StatusCode == 404 {
		log.Println("Could not find todo with with the given id.")
		// NOTE: Probably want to render a message to the user
		// in this case.
		return
	}
	// Handle some other type of error.
	log.Fatal(err)
}
```

Testing
-------

Rest uses the [karma test runner](http://karma-runner.github.io/0.12/index.html) to test
the code running in actual browsers.

The tests require the following additional dependencies:

- [node.js](http://nodejs.org/) (If you didn't already install it above)
- [karma](http://karma-runner.github.io/0.12/index.html)
- [karma-qunit](https://github.com/karma-runner/karma-qunit)

Don't forget to also install the karma command line tools with `npm install -g karma-cli`.

You will also need to install a launcher for each browser you want to test with, as well as the
browsers themselves. Typically you install a karma launcher with `npm install -g karma-chrome-launcher`.
You can edit the config file at `karma/test-mac.conf.js` or create a new one (e.g. `karma/test-windows.conf.js`)
if you want to change the browsers that are tested on.

Once you have installed all the dependencies, start karma with `karma start karma/test-mac.conf.js` (or 
your customized config file, if applicable). Once karma is running, you can keep it running in between tests.

The tests communicate with a special idempotent test server, which does some basic validation and mocks
the kinds of responses we would expect from a real REST server. Before running the tests, you need to start
the test server with `go run test/server/main.go`.

Next you need to compile the test.go file to javascript so it can run in the browsers:

```
gopherjs build karma/go/rest_test.go -o karma/js/rest_test.js
```

Finally run the tests with `karma run karma/test-mac.conf.js` (changing the name of the config file if needed).

If you are on a unix-like operating system, you can recompile and run the tests in one go by running
the provided bash script: `./karma/test.sh`.


Contributing
------------

See [CONTRIBUTING.md](https://github.com/go-humble/rest/blob/master/CONTRIBUTING.md)


License
-------

Rest is licensed under the MIT License. See the [LICENSE](https://github.com/go-humble/rest/blob/master/LICENSE)
file for more information.
