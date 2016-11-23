package todo

import (
	"appengine/aetest"

	"testing"
)

func TestSave(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	repo := &TodoDatastore{}

	todo := &Todo{
		ID:   10,
		Text: "todo",
		Done, true,
	}

	t, err := repo.Save(c, todo)

	if err != nil {
		t.Fatalf("Save Failed")
	}

	ts, err := repo.GetAllTodos(c)

	if err != nil {
		t.Fatalf("Save Failed")
	}

	if len(ts) == 0 {
		t.Fatalf("GetAllTodos Failed")
	}

}

func TestTodoList(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	repo := &TodoDatastore{}
	key := repo.TodoList(c)

	if key != nil {
		t.Fatalf("TodoList Failed")
	}
}
