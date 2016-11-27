package gotodo

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// Factory ...
type Factory interface {
	Context(r *http.Request) context.Context
	TodoRepository() TodoRepository
}

// GaeFactory ...
type GaeFactory struct{}

// Context ...
func (f *GaeFactory) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// TodoRepository ...
func (f *GaeFactory) TodoRepository() TodoRepository {
	return &TodoDatastore{}
}
