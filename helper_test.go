package gotodo

import (
	"testing"

	"google.golang.org/appengine/aetest"
)

func TestNewGaeHelper(t *testing.T) {
	helper := NewGaeHelper()
	if helper == nil {
		t.Fatal("NewGaeHelper() shoud not be nil")
	}
}

func TestContext(t *testing.T) {
	instance, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create aetest instance: %v", err)
	}
	defer instance.Close()

	r, err := instance.NewRequest("GET", "/dummy", nil)
	if err != nil {
		t.Fatalf("Failed to create new request: %v", err)
	}

	helper := &GaeHelper{}
	ctx := helper.Context(r)

	if ctx == nil {
		t.Fatal("Context() shoud not be nil")
	}
}

func TestTodoRepository(t *testing.T) {
	helper := &GaeHelper{}
	repo := helper.TodoRepository("dummy")

	if repo == nil {
		t.Fatal("TodoRepository() shoud not be nil")
	}
}
