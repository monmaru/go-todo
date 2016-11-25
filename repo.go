package gotodo

import (
	"time"

	"golang.org/x/net/context"
	ds "google.golang.org/appengine/datastore"
)

// TodoRepo ...
type TodoRepo interface {
	SaveTodo(c context.Context, t *Todo) (*Todo, error)
	GetTodo(c context.Context, id int64) (*Todo, error)
	GetAllTodos(c context.Context) ([]Todo, error)
	DeleteTodo(c context.Context, id int64) error
	DeleteDoneTodos(c context.Context) error
}

// TodoDatastore ...
type TodoDatastore struct {
}

func (store *TodoDatastore) getTodoListKey(c context.Context) *ds.Key {
	return ds.NewKey(c, "TodoList", "default", 0, nil)
}

func (store *TodoDatastore) getTodoKey(c context.Context, t *Todo) *ds.Key {
	if t.ID == 0 {
		t.Created = time.Now()
		return ds.NewIncompleteKey(c, "Todo", store.getTodoListKey(c))
	}
	return ds.NewKey(c, "Todo", "", t.ID, store.getTodoListKey(c))
}

// SaveTodo ...
func (store *TodoDatastore) SaveTodo(c context.Context, t *Todo) (*Todo, error) {
	key, err := ds.Put(c, store.getTodoKey(c, t), t)
	if err != nil {
		return nil, err
	}
	t.ID = key.IntID()
	return t, nil
}

// GetTodo ...
func (store *TodoDatastore) GetTodo(c context.Context, id int64) (*Todo, error) {
	todo := &Todo{}
	key := ds.NewKey(c, "Todo", "", id, store.getTodoListKey(c))
	err := ds.Get(c, key, todo)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

// GetAllTodos ...
func (store *TodoDatastore) GetAllTodos(c context.Context) ([]Todo, error) {
	todos := []Todo{}
	keys, err := ds.NewQuery("Todo").Ancestor(store.getTodoListKey(c)).Order("Created").GetAll(c, &todos)
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
	key := ds.NewKey(c, "Todo", "", id, store.getTodoListKey(c))
	return ds.Delete(c, key)
}

// DeleteDoneTodos ...
func (store *TodoDatastore) DeleteDoneTodos(c context.Context) error {
	return ds.RunInTransaction(c, func(c context.Context) error {
		keys, err := ds.NewQuery("Todo").KeysOnly().Ancestor(store.getTodoListKey(c)).Filter("Done=", true).GetAll(c, nil)
		if err != nil {
			return err
		}
		return ds.DeleteMulti(c, keys)
	}, nil)
}
