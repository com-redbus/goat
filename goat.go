package goat

import (
	"net/http"
)

//Middleware type which accepts a http.Handler and returns a http.Handler
type Middleware func(http.Handler) http.Handler

//MiddlewareChain struct contains array of all the middlewares that are in the chain
type MiddlewareChain struct {
	middlewares []Middleware
}

//New func accepts any number of middleware type func which are used to create a new middleware chain
func New(middlewares ...Middleware) MiddlewareChain {
	var m []Middleware
	m = append(m, middlewares...)
	return MiddlewareChain{
		middlewares: m,
	}
}

//Then func to added handler to the middleware chain
func (mc MiddlewareChain) Then(handler http.Handler) http.Handler {
	if handler == nil {
		handler = http.DefaultServeMux
	}
	for i := range mc.middlewares {
		handler = mc.middlewares[len(mc.middlewares)-i-1](handler)
	}

	return handler
}

//ThenFunc func in similiar to Then accept this accepts a http.HandlerFunc rather than http.Handler
func (mc MiddlewareChain) ThenFunc(handlerFunc http.HandlerFunc) http.Handler {
	return mc.Then(handlerFunc)
}

//Append func creates a new middleware without touching the original middleware chain
func (mc MiddlewareChain) Append(middlewares ...Middleware) MiddlewareChain {
	var newMiddlewares []Middleware
	newMiddlewares = append(newMiddlewares, mc.middlewares...)
	newMiddlewares = append(newMiddlewares, middlewares...)
	return MiddlewareChain{
		middlewares: newMiddlewares,
	}
}

//AppendToChain func append a middleware to the current middleware chain
func (mc MiddlewareChain) AppendToChain(middlewares ...Middleware) MiddlewareChain {
	mc.middlewares = append(mc.middlewares, middlewares...)
	return mc
}

//CommonMiddlewares func for crearting a few common middlewares like logger, nocache header and recovery
func CommonMiddlewares() MiddlewareChain {
	mc := New(NoCache, Recovery, Logger)
	return mc
}
