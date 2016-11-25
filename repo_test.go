package gotodo

import (
	"google.golang.org/appengine/aetest"
	ds "google.golang.org/appengine/datastore"

	"testing"
)

func TestGetTodoListKey(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{}

	if key := repo.getTodoListKey(ctx); key == nil {
		t.Fatalf("key should not be nil")
	}
}

func TestGetTodoKey(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{}

	for _, id := range []int64{0, 1} {
		todo := &Todo{
			ID:   id,
			Text: "todo",
			Done: false,
		}

		if key := repo.getTodoKey(ctx, todo); key == nil {
			t.Fatalf("key should not be nil")
		}
	}
}

func TestSave(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}

	defer done()

	repo := &TodoDatastore{}

	expected := &Todo{
		ID:   10,
		Text: "todo",
		Done: false,
	}

	saved, err := repo.SaveTodo(ctx, expected)

	if saved == nil {
		t.Fatalf("saved should not be nil")
	}

	if err != nil {
		t.Fatal(err)
	}

	result, err := repo.GetTodo(ctx, expected.ID)

	if err != nil {
		t.Fatal(err)
	}

	if expected.Text != result.Text {
		t.Fatalf("GetAllTodos Failed")
	}
}

func TestTodoList(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{}

	if key := repo.getTodoListKey(ctx); key == nil {
		t.Fatalf("TodoList Failed")
	}
}

func TestDeleteTodoAndGetTodo(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{}
	todo := &Todo{
		ID:   10,
		Text: "todo",
		Done: false,
	}

	if _, err := ds.Put(ctx, repo.getTodoKey(ctx, todo), todo); err != nil {
		t.Fatal(err)
	}

	if err = repo.DeleteTodo(ctx, todo.ID); err != nil {
		t.Fatal(err)
	}

	todos, err := repo.GetAllTodos(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 0 {
		t.Fatalf("todos should be deleted")
	}
}

func TestDeleteDoneTodos(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{}
	todo := &Todo{
		ID:   10,
		Text: "todo",
		Done: true,
	}

	if _, err := ds.Put(ctx, repo.getTodoKey(ctx, todo), todo); err != nil {
		t.Fatal(err)
	}

	if err = repo.DeleteDoneTodos(ctx); err != nil {
		t.Fatal(err)
	}

	todos, err := repo.GetAllTodos(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 0 {
		t.Fatalf("todos should be deleted")
	}
}

func TestGetAllTodos(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{}
	todos := []*Todo{
		&Todo{
			ID:   0,
			Text: "todo",
			Done: false,
		},
		&Todo{
			ID:   0,
			Text: "todo",
			Done: true,
		},
	}

	for _, todo := range todos {
		if _, err := ds.Put(ctx, repo.getTodoKey(ctx, todo), todo); err != nil {
			t.Fatal(err)
		}
	}

	result, err := repo.GetAllTodos(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(todos) {
		t.Fatalf("result length should be %#v", len(todos))
	}
}
