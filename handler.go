package gotodo

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// TodosHandler ...
type TodosHandler struct {
	factory Factory
}

// GetAllTodos ...
func (t *TodosHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	t.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		return repo.GetAllTodos(c)
	})
}

// GetTodo ...
func (t *TodosHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	t.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		id, err := t.parseID(r)

		if err != nil {
			return nil, err
		}

		return repo.GetTodo(c, id)
	})
}

// DeleteDoneTodos ...
func (t *TodosHandler) DeleteDoneTodos(w http.ResponseWriter, r *http.Request) {
	t.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		return nil, repo.DeleteDoneTodos(c)
	})
}

// DeleteTodo ...
func (t *TodosHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	t.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		id, err := t.parseID(r)

		if err != nil {
			return nil, err
		}

		return nil, repo.DeleteTodo(c, id)
	})
}

// SaveTodo ...
func (t *TodosHandler) SaveTodo(w http.ResponseWriter, r *http.Request) {
	t.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		todo, err := t.json2Todo(r.Body)
		if err != nil {
			return nil, err
		}
		return repo.SaveTodo(c, todo)
	})
}

func (t *TodosHandler) commonJSONHandler(
	w http.ResponseWriter, r *http.Request,
	f func(c context.Context, repo TodoRepo) (interface{}, error)) {

	c := t.factory.Context(r)
	repo := t.factory.TodoRepo()

	val, err := f(c, repo)
	if err == nil {
		err = json.NewEncoder(w).Encode(val)
	}

	if err != nil {
		log.Errorf(c, "todo error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TodosHandler) json2Todo(r io.ReadCloser) (*Todo, error) {
	defer r.Close()
	var todo Todo
	err := json.NewDecoder(r).Decode(&todo)
	return &todo, err
}

func (t *TodosHandler) parseID(r *http.Request) (i int64, err error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["id"], 10, 64)
}
