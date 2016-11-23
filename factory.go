package todo

import (
	"net/http"

	"appengine"
)

type Factory interface {
	Context(r *http.Request) appengine.Context
	TodoRepo() TodoRepo
}

type GaeFactory struct {
}

func (f *GaeFactory) Context(r *http.Request) appengine.Context {
	return appengine.NewContext(r)
}

func (f *GaeFactory) TodoRepo() TodoRepo {
	return &TodoDatastore{}
}
