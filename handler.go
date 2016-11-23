package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type TodosHandler struct {
	factory Factory
}

func (t *TodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := t.factory.Context(r)
	repo := t.factory.TodoRepo()

	val, err := t.HandleTodos(r, c, repo)
	if err == nil {
		err = json.NewEncoder(w).Encode(val)
	}
	if err != nil {
		log.Errorf(c, "todo error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TodosHandler) HandleTodos(r *http.Request, c context.Context, repo TodoRepo) (interface{}, error) {
	switch r.Method {
	case "POST":
		todo, err := t.JSON2Todo(r.Body)
		if err != nil {
			return nil, err
		}
		return repo.Save(c, todo)
	case "GET":
		return repo.GetAllTodos(c)
	case "DELETE":
		return nil, repo.DeleteDoneTodos(c)
	}
	return nil, fmt.Errorf("method not implemented")
}

func (t *TodosHandler) JSON2Todo(r io.ReadCloser) (*Todo, error) {
	defer r.Close()
	var todo Todo
	err := json.NewDecoder(r).Decode(&todo)
	return &todo, err
}
