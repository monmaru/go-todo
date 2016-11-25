package gotodo

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// Factory ...
type Factory interface {
	Context(r *http.Request) context.Context
	TodoRepo() TodoRepo
}

// GaeFactory ...
type GaeFactory struct {
}

// Context ...
func (f *GaeFactory) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// TodoRepo ...
func (f *GaeFactory) TodoRepo() TodoRepo {
	return &TodoDatastore{}
}
