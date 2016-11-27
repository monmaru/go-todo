package gotodo

import "net/http"

func init() {
	http.Handle("/", Router(&GaeFactory{}))
}
