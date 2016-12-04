package gotodo

import (
	"google.golang.org/appengine/aetest"
	ds "google.golang.org/appengine/datastore"

	"testing"
)

func TestUserKey(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{UserID: "default"}

	if key := repo.userKey(ctx); key == nil {
		t.Fatalf("key should not be nil")
	}
}

func TestCreateTodo(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{UserID: "default"}

	expected := &Todo{
		Text: "todo",
		Done: false,
	}

	created, err := repo.CreateTodo(ctx, expected)

	if created == nil {
		t.Fatalf("created should not be nil")
	}

	if err != nil {
		t.Fatal(err)
	}

	result, err := repo.ReadTodo(ctx, created.ID)

	if err != nil {
		t.Fatal(err)
	}

	if expected.Text != result.Text {
		t.Fatalf("Failed to create")
	}
}

func TestUpdateTodo(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{UserID: "default"}

	todo := &Todo{
		Text: "todo",
		Done: false,
	}

	key := ds.NewIncompleteKey(ctx, kind, repo.userKey(ctx))
	if _, err := ds.Put(ctx, key, todo); err != nil {
		t.Fatal(err)
	}

	todo.Done = true

	updated, err := repo.UpdateTodo(ctx, todo)

	if updated == nil {
		t.Fatalf("updated should not be nil")
	}

	if err != nil {
		t.Fatal(err)
	}

	result, err := repo.ReadTodo(ctx, updated.ID)

	if err != nil {
		t.Fatal(err)
	}

	if !result.Done {
		t.Fatalf("Failed to update")
	}
}

func TestDeleteTodoAndReadTodo(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{UserID: "default"}
	todo := &Todo{
		Text: "todo",
		Done: false,
	}

	created, err := repo.CreateTodo(ctx, todo)
	if err != nil {
		t.Fatal(err)
	}

	if err = repo.DeleteTodo(ctx, created.ID); err != nil {
		t.Fatal(err)
	}

	todos, err := repo.ReadAllTodos(ctx)

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

	repo := &TodoDatastore{UserID: "default"}
	todo := &Todo{
		ID:   10,
		Text: "todo",
		Done: true,
	}

	key := ds.NewKey(ctx, kind, "", todo.ID, repo.userKey(ctx))
	if _, err := ds.Put(ctx, key, todo); err != nil {
		t.Fatal(err)
	}

	if err = repo.DeleteDoneTodos(ctx); err != nil {
		t.Fatal(err)
	}

	todos, err := repo.ReadAllTodos(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 0 {
		t.Fatalf("todos should be deleted")
	}
}

func TestReadAllTodos(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := &TodoDatastore{UserID: "default"}
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
		key := ds.NewKey(ctx, kind, "", todo.ID, repo.userKey(ctx))
		if _, err := ds.Put(ctx, key, todo); err != nil {
			t.Fatal(err)
		}
	}

	result, err := repo.ReadAllTodos(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(todos) {
		t.Fatalf("result length should be %#v", len(todos))
	}
}
