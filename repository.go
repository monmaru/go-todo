package gotodo

import (
	"time"

	"golang.org/x/net/context"
	ds "google.golang.org/appengine/datastore"
)

// TodoRepository ...
type TodoRepository interface {
	CreateTodo(c context.Context, todo *Todo) (*Todo, error)
	ReadTodo(c context.Context, id int64) (*Todo, error)
	ReadAllTodos(c context.Context) ([]Todo, error)
	UpdateTodo(c context.Context, todo *Todo) (*Todo, error)
	DeleteTodo(c context.Context, id int64) error
	DeleteDoneTodos(c context.Context) error
}

const (
	parent = "TodoList"
	kind   = "Todo"
)

// TodoDatastore ...
type TodoDatastore struct{}

func (store *TodoDatastore) parentKey(c context.Context) *ds.Key {
	return ds.NewKey(c, parent, "default", 0, nil)
}

// UpdateTodo ...
func (store *TodoDatastore) UpdateTodo(c context.Context, todo *Todo) (*Todo, error) {
	key := ds.NewKey(c, kind, "", todo.ID, store.parentKey(c))
	key, err := ds.Put(c, key, todo)
	if err != nil {
		return nil, err
	}
	todo.ID = key.IntID()
	return todo, nil
}

// CreateTodo ...
func (store *TodoDatastore) CreateTodo(c context.Context, todo *Todo) (*Todo, error) {
	todo.Created = time.Now()
	key := ds.NewIncompleteKey(c, kind, store.parentKey(c))
	key, err := ds.Put(c, key, todo)
	if err != nil {
		return nil, err
	}
	todo.ID = key.IntID()
	return todo, nil
}

// ReadTodo ...
func (store *TodoDatastore) ReadTodo(c context.Context, id int64) (*Todo, error) {
	todo := &Todo{}
	key := ds.NewKey(c, kind, "", id, store.parentKey(c))
	if err := ds.Get(c, key, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

// ReadAllTodos ...
func (store *TodoDatastore) ReadAllTodos(c context.Context) ([]Todo, error) {
	todos := []Todo{}
	keys, err := ds.NewQuery(kind).Ancestor(store.parentKey(c)).Order("Created").GetAll(c, &todos)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(todos); i++ {
		todos[i].ID = keys[i].IntID()
	}
	return todos, nil
}

// DeleteTodo ...
func (store *TodoDatastore) DeleteTodo(c context.Context, id int64) error {
	key := ds.NewKey(c, kind, "", id, store.parentKey(c))
	return ds.Delete(c, key)
}

// DeleteDoneTodos ...
func (store *TodoDatastore) DeleteDoneTodos(c context.Context) error {
	return ds.RunInTransaction(c, func(c context.Context) error {
		keys, err := ds.NewQuery(kind).KeysOnly().Ancestor(store.parentKey(c)).Filter("Done=", true).GetAll(c, nil)
		if err != nil {
			return err
		}
		return ds.DeleteMulti(c, keys)
	}, nil)
}
