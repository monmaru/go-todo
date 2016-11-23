package todo

import "net/http"

func init() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/api/todos", &TodosHandler{
		factory: &GaeFactory{},
	})
}
