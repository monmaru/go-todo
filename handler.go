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

// APIHandler ...
type APIHandler struct {
	factory Factory
}

// GetAllTodos ...
func (h *APIHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	h.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		return repo.GetAllTodos(c)
	})
}

// GetTodo ...
func (h *APIHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	h.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		id, err := h.parseID(r)

		if err != nil {
			return nil, err
		}

		return repo.GetTodo(c, id)
	})
}

// DeleteDoneTodos ...
func (h *APIHandler) DeleteDoneTodos(w http.ResponseWriter, r *http.Request) {
	h.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		return nil, repo.DeleteDoneTodos(c)
	})
}

// DeleteTodo ...
func (h *APIHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	h.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		id, err := h.parseID(r)

		if err != nil {
			return nil, err
		}

		return nil, repo.DeleteTodo(c, id)
	})
}

// SaveTodo ...
func (h *APIHandler) SaveTodo(w http.ResponseWriter, r *http.Request) {
	h.commonJSONHandler(w, r, func(c context.Context, repo TodoRepo) (interface{}, error) {
		todo, err := h.json2Todo(r.Body)
		if err != nil {
			return nil, err
		}
		return repo.SaveTodo(c, todo)
	})
}

func (h *APIHandler) commonJSONHandler(
	w http.ResponseWriter, r *http.Request,
	f func(c context.Context, repo TodoRepo) (interface{}, error)) {

	c := h.factory.Context(r)
	repo := h.factory.TodoRepo()

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

func (h *APIHandler) json2Todo(r io.ReadCloser) (*Todo, error) {
	defer r.Close()
	var todo Todo
	err := json.NewDecoder(r).Decode(&todo)
	return &todo, err
}

func (h *APIHandler) parseID(r *http.Request) (i int64, err error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["id"], 10, 64)
}
