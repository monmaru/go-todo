package todo

import "testing"

func TestTodoRepo(t *testing.T) {
	f := &GaeFactory{}
	repo := f.TodoRepo()

	if repo == nil {
		t.Fatal("TodoRepo() shoud not be nil")
	}
}
