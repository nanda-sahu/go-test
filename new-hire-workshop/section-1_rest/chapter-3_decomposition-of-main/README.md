# Chapter 3: Decomposition of Main

We now have a working REST server. However, as you've probably noticed, our `main.go` is becoming very bloated. In a real application, we tend to want to keep our main as simple as clean as possible. Usually it is just responsible for setting up everything we need and then running the application. Therefore, in this chapter we're going to decompose our `main.go` file into several packages to make things easier to navigate as we grow the server.

## Go Application Structure Conventions

By convention in Golang, there are three main locations for our code:

* `cmd` - The directory where the main function should live. (i.e. the command that runs the application)
* `internal` - A directory that contains all the code that should only be used for the *current* application.
* `pkg` - A directory that contains code that can be shared and re-used in *many* applications.

In the majority of the services we are going to write, almost all of the code we write will live in `internal`. This is because any code that would live in `pkg` is probably better suited to going into a common code library such as [go-gadgets](https://github.hpe.com/cloud/go-gadgets)

In the `internal` directory, we should be splitting our code up into packages, with a directory being responsible for containing all of a package’s code - effectively a 1:1 mapping. It is convention to name the directory the same as the package name. For a good starter on package naming, refer to [this blog.](https://go.dev/blog/package-names)

Unlike other languages where access controls are usually done using keywords such as `public` or `private`, in Go we declare variables, constants, or functions with a capital letter (e.g. `Hello`, `SendMessage`, `ParseUUID`) to allow other packages to use them. These are called [Exported Names](https://go.dev/tour/basics/3). The inverse is also true with variables, constants, and functions with a lowercase letter (e.g. `formatDate`, `validateEmail`, `createUUID`) cannot be used in other packages, and can only be “seen” from inside it’s own package.

Further on in this tutorial we will cover some more advanced concepts around how to structure our application using the Clean Architecture approach. For the moment, or application is too simple to warrant the overhead that this would bring. Instead for the moment we're just going to separate our application into our 'core' and 'rest' layers.

At the end of this chapter we should have a directory structure that looks something like this:
```
cmd
└── rest-server
    └── main.go
internal
├── orders
│   ├── order.go
│   └── orderstore.go
└── rest
    ├── handlers.go
    ├── models.go
    └── server.go
```

## `orders` Package

The orders package will be quite a simple one, as there is not much 'core' logic that we have in our REST server. In fact, at the moment we could just chuck the core logic in a `internal/orders/order.go` file and call it a day. However, let's also take the time to make our interaction with this package a bit cleaner.

Currently the `orders` package has a slice of pre-populated orders to use for our manual testing. In a real application we wouldn't want to just expose a slice from a package in such a manner. Not only is it not thread safe, but it removes any ability to add business logic around the modification of orders. Instead, we want to create an access point for the package that other packages can be used to interact with (and in the future abstract away).

## Exercise: Implement an OrderStore struct

The first exercise for this chapter is to write an OrderStore struct that controls access to the 'core' order data. This struct should have a receiver function with the following signature:
```go
func (o *OrderStore) GetOrders(customerID string) []Order
```
Whilst so far we've been dealing with value types for the Order structs (and therefore value semantics) for this function we will make it a pointer receiver function. This is because we intend to only have one `OrderStore` in the application, and instead to pass references around where it is needed.

A template for the `orders` package can be found [here](./materials/orders).

<details><summary><b>Hint</b></summary>

* You've already written the logic for the `GetOrders` function itself, we just need to encapsulate the orders slice with the `OrderStore` struct.
* Whilst Go structs do not have constructors, it is common to see functions in the form of `NewFoo() *Foo` - a `NewOrderStore` function would be a good place to initialise the `OrderStore` struct with the current hard coded list of 'starting' orders.

</details>

## `context.Context`

A concept that is ubiquitous in Go applications that send/receive requests to/from external services is the idea of a [`context`](https://pkg.go.dev/context). This context contains 'request scoped' information that is intended to be propogated through the system. Whilst there is a large amount that can be done with a context, the two main uses are to provide a timeout for the operation and store request scoped variables.

To set a lifetime for the context (and therefore the request that it is associated with) you can use the [`WithTimeout`](https://pkg.go.dev/context#WithTimeout) function, which returns both the context with the timeout set and a function to cancel the context before the timeout. The pattern of usage is to defer the cancel function so that if the function exits early the context is cancelled immediately, rather than waiting for the timeout.
```go
func slowOperationWithTimeout(ctx context.Context) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()  // releases resources if slowOperation completes before timeout elapses
	return slowOperation(ctx)
}
```

Contexts can be used to hold variables. However, this feature should be avoided in most cases, as it is always better to pass parameters to a function explicitly. The only time that the context should be used to pass variables is when you are working with middleware where you need to stick to a pre-defined interface that doesn't let you pass the variables explicitly. 

Examples of this can be seen in the [`restserver/middleware`](https://github.hpe.com/cloud/go-gadgets/blob/x/restserver/v0.0.1/x/restserver/middleware/logger.go#L42) package, where the `http.Handler` interface means that we can't pass the logger as an explicit parameter.

## Structured Logging

Before we refactor our next package, this is a good point to introduce logging. For most basic applications, the logic provided by the [`log`](https://pkg.go.dev/log) is sufficient to provide information about it's operation. However, for production applications we want to instead use structured logging.

The idea behind structured logging is that rather than embedding the important data in an error message, instead these details are added in a key-value map (which is usually logged out in JSON format), which makes our logging ingestion and searching much easier to do. For example, contrast the two approaches below:
```go
func FooUnstructured(inputVal string) {
    log.Printf("calling Foo with value %s", inputValue)
    // Output: "calling FooUnstructured with value Hello, World"
}

func FooStructured(inputValue string, logger logging.Logger) {
    logger.WithField("input-val", inputVal).Info("calling FooStructured")
    // Output: "{
    //     "input-val": "Hello, World",
    //     "message": "calling FooStructured"
    // }
}
```

Whilst for the moment we've been avoiding too many dependencies, no one wants to re-implement their logging logic for each application. Therefore, we are going to make use of the [`go-gadgets/x/logging`](https://github.hpe.com/cloud/go-gadgets/tree/main/x/logging) module to provide us with structured logging. This logger is also written to be compliant with the DSCC structured logging standards.

## Exercise: `rest` Package and REST Server

The second package that we're going to pull out from `main.go` is a `rest` package to hold all our REST logic. A template for this package can be found in the [`materials/rest`](./materials/rest) folder.

The first thing we are going to implement for this package is a basic REST server. This REST server will wrap the `http.Server` from the standard library, giving it a nice clean interface to use in our application. The server will expose `Start` and `Stop` methods, as well as allow routing configuration and a base context.

When a REST server handles a request, one of the the things it does in spawn a new `context.Context` for that request. If left unconfigured, a new 'background' context will be created. However, we can configure the function that the server uses to generate this context to use a context we provide. This allows us to specify a 'parent' context, which can be used to set a global request timeout, or cancel all in flight requests in the case of server shutdown.

For this exercise populate the `server.go` template with the required logic, migrating and building upon the REST server logic from `main.go` in the previous chapter.

<details><summary><b>Hint</b></summary>

* Whilst it can be tempting to make the server's `Start` function to call a goroutine to make it non-blocking, it is better to make the `Start` function synchronous. This allows the most control for the consumer, as they can always call `Start` in a goroutine themselves if they need to.
* By adding a logger to the server, we can log out startup and shutdown events, with any relevant information.
* We haven't yet refactored the HTTP handler yet, so the `NewMux` function can just return a blank mux for the time being.

</details>

## Exercise: REST Models and Conversion

We can now also migrate our REST models (structs and enums) to the new `rest` package. As they relate to the REST models, we can also migrate the conversion functions to this file. Because we are now namespaced in the `rest` package, we can also drop all the 'REST' prefixes.

## Exercise: ListOrderSummaries Handler Struct

The final part to migrate across is the REST handler for listing our order summaries. However, currently our handler is specified as a [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc), by having the specified signature. As part of the migration, we should convert it into a struct that implements the [`http.Handler`](https://pkg.go.dev/net/http#Handler) interface. This will allow us to have members on the struct for use in the handler, such as the logger and the `OrderStore`.

<details><summary><b>Hint</b></summary>

* Any common handler logic (e.g. the `writeError()` function in the sample solution of the previous chapter can also be moved into the `handlers.go` file for the moment.
* Now that we have a logger, it would be good to add logging of errors, particularly in situations where the causing error is not included in the response (i.e. when we return a 500 InternalServerError).

</details>

## Exercise: Refactoring `main.go`

To run our code and bring all these packages together, we need to write a new `main.go`. We want to create a main function that starts the HTTP server, and stops the server if it gets an interruption or a terminate signal. As we're using a Go version higher than 1.16, we can use the [`signal.NotifyContext`](https://pkg.go.dev/os/signal#NotifyContext) to handle listening for terminations. We can also then provide this context to the REST server to automatically cancel requests on shutdown.

Because all our code has now been nicely decomposed, our `main.go`' is effectively responsible to instantiating all our resource, linking them up and then just running it. In a production app this would usually be done with the aid of a [dependency injection](https://github.hpe.com/cloud/go-chapter/tree/master/guides/dependency-injection) library/tool, but for practice let's do it manually for this.

<details><summary><b>Hint</b></summary>

* We want to return non-zero exit codes on startup error, but we want to avoid littering our `main.go` with lots of `os.Exit()` calls.
* We will expect our `main.go` to be very brief now, as most of the logic are in the other packages.

</details>
