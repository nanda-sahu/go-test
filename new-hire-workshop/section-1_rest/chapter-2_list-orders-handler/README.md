# Chapter 2: List Order Summaries

In the first chapter we stood up a simple HTTP server with a catch all endpoint that echoed the path sent. In this chapter we're going to modify this HTTP server to be the beginning of our REST server, and we're going to implement our first REST endpoint.

## REST and Internal Structs

As discussed in the section intro, we use OpenAPI to specify our REST interfaces. In order to be able to implement the specification, we need to generate some Go structs that represent the data models defined in the spec. These can be found in the `materials` directory [here](./materials/rest_structs.go) and can be copied over into your server's `main.go` file.

You'll notice that here we are using JSON annotations for our structs. These allow the Go struct to be serialised to a JSON string and vice vera automatically using the [`encoding/json`](https://pkg.go.dev/encoding/json) package. For an example of how to do this, you can refer to [this guide](https://gobyexample.com/json).

As well as representing the API spec in our code, we will also need an 'internal' representation of our data structures. In general, we don't want to couple the way the system handles data internally with how it is presented externally. This is to say if we change our external representation of our data, we won’t need to change how it is modeled internally. This also has the added benefit of allowing us to hide important information that we might not want the end-user to see. These data structures can be found in the `materials` directory [here](./materials/internal_structs.go) and can also be copied over to your `main.go` file.

Finally, for this chapter we want to mock out the existence of a database with some orders in it so that we can easily test our API. In the [`orders.go`](./materials/orders.go) file in the materials directory there is a slice that has two orders in it. You can copy this across to your `main.go` and use it in your handler, although you'll have to write the logic to convert the struct from the 'internal' representation to the 'rest' representation before returning it from the handler.

## REST Request Routing

In a REST server, an endpoint is not just the URL path, but the combination of the URL and the HTTP method. Whilst in the example server we used the standard library router, there are third party routers that are easier to use when registering REST endpoints. For this workshop we will be using the [`gorilla/mux`](https://github.com/gorilla/mux) router, which is a drop in replacement for the standard library router.

You can look at the documentation to get a sense of the full features of the library, for this exercise we can just use the basic route matching function. For example, if we wanted to configure the router to match the endpoint `GET /api/v1/example` we could use the below code:
```go
router := mux.NewRouter()
router.HandleFunc("/api/v1/example", ExampleHandler).Methods(http.MethodGet)
```

## Go HTTP Handlers

In Go's standard library a HTTP handler is any struct that implements the [http.Handler](https://pkg.go.dev/net/http#Handler) interface, or a function that has the same signature (i.e. takes in a ResponseWriter implementation and a pointer to a Request struct). The request struct can be used to access all the information from the request such as the URL, request body etc. and the response writer is what we use to respond to the request, usually with a status code and a body.

It's important to point out that these HTTP handlers don't have a return value, however when handling error cases you still want to have an early return to ensure that you don't write extra data to the response writer. This is usually done just after a body and header are written to the ResponseWriter using `Write()` and `WriteHeader()` functions.

## Query Parameters

We want our service to be multi-tenant from the start. Multi-tenancy means that the same instance of a service works for multiple different users without exposing a user's data to any other user in the system. In this case we don't want one customer to see the orders placed by another customer. Normally this kind of user information is contained in an session ID or authentication token (in DSCC we use [OAuth JWTs](https://jwt.io/introduction)), however, to keep this workshop simple the customer's ID will be a self reported query parameter.

Therefore, for a customer's request to list orders we would expect a query like:
```
GET www.example.com/api/v1/orders?customerID=customer1
```

We can fetch query parameters like this using the [`request.FormValue`](https://pkg.go.dev/net/http#Request.FormValue) function in the REST handler function to get the value of the query parameter. If the query parameter is not found, then an empty string will be returned. Therefore, the logic to use it looks like:
```go
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
    customerID := r.FormValue("customerID")
    if customerID == "" {
        // return a BadRequest error
        return
    }

    // rest of the handler logic...
}
```

## Exercise: Implement the GetOrderSummaries Endpoint 

For our first exercise, we want to implement the GetOrderSummaries endpoint for our REST server. We can do this by refactoring the `main.go` file that we had in the previous chapter, using the information and materials discussed above.

Whilst the handler logic for this exercise is fairly simple, it is a good idea to decompose the parts into functions that the handler calls, as this will help us with refactoring in future chapters.

Once it has been implemented, we should be able to list a customer's orders using the request:
```
curl localhost:8080/api/v1/orders?customerID=customer1
```
<details><summary><b>Hint</b></summary>

To complete this exercise, as well as copying over the materials you'll need to add a router, registering the route matching the API spec. Then you'll need to write a HTTP handler that does the following:

* Get the customer's ID from the request params - returning a 404 BadRequest error if the ID is blank. Returning the error message is exactly the same as returning an OrderSummary, except you populate and marshal the Error struct.
* Loop through the existing orders to get all the ones that match the customer's ID
* Convert thes orders to the REST model
* Serialise the order structs into JSON
* Write the response out along with a 200 header.

If any errors are encountered during the above, instead of writing the orders to the response, an Error struct should be populated and written, along with the appropriate error code (usually 400 if the request contains bad data or 500 if there is an unexpected error on the server side.

</details>
