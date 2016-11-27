package gotodo

import "testing"

func TestTodoRepository(t *testing.T) {
	f := &GaeFactory{}
	repo := f.TodoRepository()

	if repo == nil {
		t.Fatal("TodoRepository() shoud not be nil")
	}
}
