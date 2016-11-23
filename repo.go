package todo

import (
	"time"

	"appengine"
	ds "appengine/datastore"
)

type TodoRepo interface {
	Save(c appengine.Context, t *Todo) (*Todo, error)
	GetAllTodos(c appengine.Context) ([]Todo, error)
	DeleteDoneTodos(c appengine.Context) error
}

type TodoDatastore struct {
}

func (store *TodoDatastore) TodoList(c appengine.Context) *ds.Key {
	return ds.NewKey(c, "TodoList", "default", 0, nil)
}

func (store *TodoDatastore) GetKey(c appengine.Context, t *Todo) *ds.Key {
	if t.ID == 0 {
		t.Created = time.Now()
		return ds.NewIncompleteKey(c, "Todo", store.TodoList(c))
	}
	return ds.NewKey(c, "Todo", "", t.ID, store.TodoList(c))
}

func (store *TodoDatastore) Save(c appengine.Context, t *Todo) (*Todo, error) {
	k, err := ds.Put(c, store.GetKey(c, t), t)
	if err != nil {
		return nil, err
	}
	t.ID = k.IntID()
	return t, nil
}

func (store *TodoDatastore) GetAllTodos(c appengine.Context) ([]Todo, error) {
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

func (store *TodoDatastore) DeleteDoneTodos(c appengine.Context) error {
	return ds.RunInTransaction(c, func(c appengine.Context) error {
		ks, err := ds.NewQuery("Todo").KeysOnly().Ancestor(store.TodoList(c)).Filter("Done=", true).GetAll(c, nil)
		if err != nil {
			return err
		}
		return ds.DeleteMulti(c, ks)
	}, nil)
}
