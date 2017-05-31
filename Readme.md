# Goat -- Easy Middlewares for Golang

Go has the **net/http** package for managing requests and response,
but sometimes we need some common code to handle repetitive code.
This is where **middlewares** comes in to the rescue.

**Goat** is simple middleware package which has some commonly used middlewares
and allows to add other middlewares because every middleware is an **http.Handler** or **http.HandlerFunc**

## Installation

```
go get -u github.com/retiredbatman/goat
```

## Features
  * Ready to use common middlewares 
  * Common Middlewares include *logging*, *global recovery*, *nocache headers* 
  * Easy to add middlewares to the middleware chain
  * http.Handler can be used as middleware
  * Writing New middlewares are very easy
  * Should with other golang web frameworks because these are basically http.Handler


## Middlewares Included
* Logger -> logs to the console 
* Recovery -> recovers from a panic globally , stops the app from crashing
* NoCache -> adds no-cache headers to prevent api responses getting cache by the browser
* Compression -> gzip compression of response data , currently supports gzip.DefaultCompression level
* Monitor -> simple metrics about the app like uptime , pid , responsecounts etc

## Usage

### Using Common Middlewares With DefaultServeMux

```go
import (
    "fmt"
    "net/http"

    "github.com/retiredbatman/goat"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Its a test handler")
}

func main() {
    mux := http.NewServeMux()

    commonMiddlewares := goat.CommonMiddlewares()
    mux.Handle("/", commonMiddlewares.ThenFunc(indexHandler))
    http.ListenAndServe(":8080", mux)
}
```


### Using Common Middlewares with Gorilla Mux

```go
import(
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/retiredbatman/goat"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Its a test handler")
}

func main() {
   router := mux.NewRouter()

    commonMiddlewares := goat.CommonMiddlewares()
    router.Handle("/", commonMiddlewares.ThenFunc(indexHandler))

    http.ListenAndServe(":8080", router)
}
```


### Appending Middlewares to the Middleware Chain

```go
import(
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/retiredbatman/goat"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Its a test handler")
}

func compressedIndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Its a test handler")
}
func main() {
   router := mux.NewRouter()

    commonMiddlewares := goat.CommonMiddlewares()
    compressionMiddlewareAdded := commonMiddlewares.Append(goat.Compression)
    router.Handle("/", commonMiddlewares.ThenFunc(indexHandler))
    router.Handle("/compressed", compressionMiddlewareAdded.ThenFunc(indexHandler))

    http.ListenAndServe(":8080", router)
}
```

### Writing your own Middleware

```go
//sampleMiddleware.go
import "net/http"


func SampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    //do something before next Handler

	    //call the next handler
	    next.ServeHTTP(w, r)
            //do something after next Handler
	})
}
```

## Future Plans

Write a helmet like middleware for golang
https://github.com/helmetjs/helmet



