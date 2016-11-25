package gotodo

import "net/http"
import "github.com/gorilla/mux"

func init() {
	r := mux.NewRouter()
	handler := &APIHandler{
		factory: &GaeFactory{},
	}

	r.HandleFunc("/api/todos", handler.GetAllTodos).Methods("GET")
	r.HandleFunc("/api/todos/{id}", handler.GetTodo).Methods("GET")
	r.HandleFunc("/api/todos", handler.SaveTodo).Methods("POST")
	r.HandleFunc("/api/todos", handler.DeleteDoneTodos).Methods("DELETE")
	r.HandleFunc("/api/todos/{id}", handler.DeleteTodo).Methods("DELETE")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.Handle("/", r)
}
