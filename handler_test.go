package gotodo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"
)

var dummyTodo = Todo{
	ID:      1111,
	Text:    "todo",
	Done:    true,
	Created: time.Now(),
}

var dummyTodos = []Todo{
	dummyTodo,
	dummyTodo,
}

type fakeRepository struct{}

func (repo *fakeRepository) CreateTodo(c context.Context, todo *Todo) (*Todo, error) {
	return &dummyTodo, nil
}

func (repo *fakeRepository) ReadTodo(c context.Context, id int64) (*Todo, error) {
	return &dummyTodo, nil
}

func (repo *fakeRepository) ReadAllTodos(c context.Context) ([]Todo, error) {
	return dummyTodos, nil
}

func (repo *fakeRepository) UpdateTodo(c context.Context, todo *Todo) (*Todo, error) {
	return &dummyTodo, nil
}

func (repo *fakeRepository) DeleteTodo(c context.Context, id int64) error {
	return nil
}

func (repo *fakeRepository) DeleteDoneTodos(c context.Context) error {
	return nil
}

type fakeFactory struct{}

func (f *fakeFactory) Context(r *http.Request) context.Context {
	return nil
}

func (f *fakeFactory) TodoRepository() TodoRepository {
	return &fakeRepository{}
}

func TestHandleGetTodo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/todos/1111111", nil)
	w := httptest.NewRecorder()

	Router(&fakeFactory{}).ServeHTTP(w, r)

	var todo Todo
	if err := json.NewDecoder(w.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dummyTodo, todo) {
		t.Fatalf("Not the same value")
	}
}

func TestHandleGetAllTodos(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/todos", nil)
	w := httptest.NewRecorder()

	Router(&fakeFactory{}).ServeHTTP(w, r)

	var todos []Todo
	if err := json.NewDecoder(w.Body).Decode(&todos); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dummyTodos, todos) {
		t.Fatalf("Not the same value")
	}
}

func TestHandlePutTodo(t *testing.T) {
	body, _ := json.Marshal(dummyTodo)
	r, _ := http.NewRequest("PUT", "/api/todos", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	defer r.Body.Close()

	Router(&fakeFactory{}).ServeHTTP(w, r)

	var todo Todo
	if err := json.NewDecoder(w.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dummyTodo, todo) {
		t.Fatalf("Not the same value")
	}
}

func TestHandlePostTodo(t *testing.T) {
	body, _ := json.Marshal(dummyTodo)
	r, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	defer r.Body.Close()

	Router(&fakeFactory{}).ServeHTTP(w, r)

	var todo Todo
	if err := json.NewDecoder(w.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dummyTodo, todo) {
		t.Fatalf("Not the same value")
	}
}

func TestHandleDeleteTodo(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/api/todos/1111111", nil)
	w := httptest.NewRecorder()

	Router(&fakeFactory{}).ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("Status Code is invalid")
	}
}

func TestHandleDeleteDoneTodos(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/api/todos", nil)
	w := httptest.NewRecorder()

	Router(&fakeFactory{}).ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("Status Code is invalid")
	}
}

func TestHandleRoot(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	Router(&fakeFactory{}).ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("Status Code is invalid")
	}
}
