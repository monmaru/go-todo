package gotodo

import (
	"testing"

	"google.golang.org/appengine/aetest"
)

func TestTodoRepository(t *testing.T) {
	factory := &GaeFactory{}
	repo := factory.TodoRepository()

	if repo == nil {
		t.Fatal("TodoRepository() shoud not be nil")
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

	factory := &GaeFactory{}
	ctx := factory.Context(r)

	if ctx == nil {
		t.Fatal("Context() shoud not be nil")
	}
}
