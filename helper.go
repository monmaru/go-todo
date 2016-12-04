package gotodo

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

// Helper ...
type Helper interface {
	Context(r *http.Request) context.Context
	CurrentUser(ctx context.Context) *user.User
	TodoRepository(userID string) TodoRepository
}

// GaeHelper ...
type GaeHelper struct{}

// NewGaeHelper ...
func NewGaeHelper() *GaeHelper {
	return &GaeHelper{}
}

// Context ...
func (h *GaeHelper) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// CurrentUser ...
func (h *GaeHelper) CurrentUser(ctx context.Context) *user.User {
	return user.Current(ctx)
}

// TodoRepository ...
func (h *GaeHelper) TodoRepository(userID string) TodoRepository {
	return &TodoDatastore{UserID: userID}
}
