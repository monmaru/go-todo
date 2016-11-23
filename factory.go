package todo

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

type Factory interface {
	Context(r *http.Request) context.Context
	TodoRepo() TodoRepo
}

type GaeFactory struct {
}

func (f *GaeFactory) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func (f *GaeFactory) TodoRepo() TodoRepo {
	return &TodoDatastore{}
}
