package todo

import "time"

type Todo struct {
	ID      int64     `json:"id" datastore:"-"`
	Text    string    `json:"text" datastore:",noindex"`
	Done    bool      `json:"done"`
	Created time.Time `json:"created"`
}
