package todo

import (
	"time"

	"golang.org/x/net/context"
	ds "google.golang.org/appengine/datastore"
)

type TodoRepo interface {
	Save(c context.Context, t *Todo) (*Todo, error)
	GetAllTodos(c context.Context) ([]Todo, error)
	DeleteDoneTodos(c context.Context) error
}

type TodoDatastore struct {
}

func (store *TodoDatastore) TodoList(c context.Context) *ds.Key {
	return ds.NewKey(c, "TodoList", "default", 0, nil)
}

func (store *TodoDatastore) GetKey(c context.Context, t *Todo) *ds.Key {
	if t.ID == 0 {
		t.Created = time.Now()
		return ds.NewIncompleteKey(c, "Todo", store.TodoList(c))
	}
	return ds.NewKey(c, "Todo", "", t.ID, store.TodoList(c))
}

func (store *TodoDatastore) Save(c context.Context, t *Todo) (*Todo, error) {
	k, err := ds.Put(c, store.GetKey(c, t), t)
	if err != nil {
		return nil, err
	}
	t.ID = k.IntID()
	return t, nil
}

func (store *TodoDatastore) GetAllTodos(c context.Context) ([]Todo, error) {
	todos := []Todo{}
	ks, err := ds.NewQuery("Todo").Ancestor(store.TodoList(c)).Order("Created").GetAll(c, &todos)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(todos); i++ {
		todos[i].ID = ks[i].IntID()
	}
	return todos, nil
}

func (store *TodoDatastore) DeleteDoneTodos(c context.Context) error {
	return ds.RunInTransaction(c, func(c context.Context) error {
		ks, err := ds.NewQuery("Todo").KeysOnly().Ancestor(store.TodoList(c)).Filter("Done=", true).GetAll(c, nil)
		if err != nil {
			return err
		}
		return ds.DeleteMulti(c, ks)
	}, nil)
}
