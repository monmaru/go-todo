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

// Router ...
func Router(factory Factory) *mux.Router {
	r := mux.NewRouter()
	api := &API{factory: factory}

	r.HandleFunc("/api/todos", api.HandleGetAllTodos).Methods("GET")
	r.HandleFunc("/api/todos/{id}", api.HandleGetTodo).Methods("GET")
	r.HandleFunc("/api/todos", api.HandlePutTodo).Methods("PUT")
	r.HandleFunc("/api/todos", api.HandlePostTodo).Methods("POST")
	r.HandleFunc("/api/todos", api.HandleDeleteDoneTodos).Methods("DELETE")
	r.HandleFunc("/api/todos/{id}", api.HandleDeleteTodo).Methods("DELETE")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	return r
}

// API ...
type API struct {
	factory Factory
}

// HandleGetAllTodos ...
func (api *API) HandleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	api.commonHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		return repo.ReadAllTodos(c)
	})
}

// HandleGetTodo ...
func (api *API) HandleGetTodo(w http.ResponseWriter, r *http.Request) {
	api.commonHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		id, err := api.parseID(r)
		if err != nil {
			return nil, err
		}
		return repo.ReadTodo(c, id)
	})
}

// HandleDeleteDoneTodos ...
func (api *API) HandleDeleteDoneTodos(w http.ResponseWriter, r *http.Request) {
	api.commonHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		return nil, repo.DeleteDoneTodos(c)
	})
}

// HandleDeleteTodo ...
func (api *API) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	api.commonHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		id, err := api.parseID(r)
		if err != nil {
			return nil, err
		}
		return nil, repo.DeleteTodo(c, id)
	})
}

// HandlePutTodo ...
func (api *API) HandlePutTodo(w http.ResponseWriter, r *http.Request) {
	api.commonHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		todo, err := api.json2Todo(r.Body)
		if err != nil {
			return nil, err
		}
		return repo.UpdateTodo(c, todo)
	})
}

// HandlePostTodo ...
func (api *API) HandlePostTodo(w http.ResponseWriter, r *http.Request) {
	api.commonHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		todo, err := api.json2Todo(r.Body)
		if err != nil {
			return nil, err
		}
		return repo.CreateTodo(c, todo)
	})
}

func (api *API) commonHandler(
	w http.ResponseWriter, r *http.Request,
	fn func(c context.Context, repo TodoRepository) (interface{}, error)) {

	ctx := api.factory.Context(r)
	repo := api.factory.TodoRepository()

	val, err := fn(ctx, repo)
	if err == nil {
		err = json.NewEncoder(w).Encode(val)
	}

	if err != nil {
		log.Errorf(ctx, "todo error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) json2Todo(r io.ReadCloser) (*Todo, error) {
	defer r.Close()
	var todo Todo
	err := json.NewDecoder(r).Decode(&todo)
	return &todo, err
}

func (api *API) parseID(r *http.Request) (id int64, err error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["id"], 10, 64)
}
