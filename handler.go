package gotodo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

// Router ...
func Router(helper Helper) *mux.Router {
	router := mux.NewRouter()
	app := &App{helper: helper}
	router.HandleFunc("/api/todos", app.HandleGetAllTodos).Methods("GET")
	router.HandleFunc("/api/todos/{id}", app.HandleGetTodo).Methods("GET")
	router.HandleFunc("/api/todos", app.HandlePutTodo).Methods("PUT")
	router.HandleFunc("/api/todos", app.HandlePostTodo).Methods("POST")
	router.HandleFunc("/api/todos", app.HandleDeleteDoneTodos).Methods("DELETE")
	router.HandleFunc("/api/todos/{id}", app.HandleDeleteTodo).Methods("DELETE")
	router.HandleFunc("/", app.HandleRoot).Methods("GET")
	return router
}

// App ...
type App struct {
	helper Helper
}

var _userWrapper = struct {
	LoginURL  func(c context.Context, dest string) (string, error)
	LogoutURL func(c context.Context, dest string) (string, error)
}{
	user.LoginURL,
	user.LogoutURL,
}

// HandleRoot ...
func (app *App) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := app.helper.Context(r)
	u := app.helper.CurrentUser(ctx)
	if u == nil {
		url, err := _userWrapper.LoginURL(ctx, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	url, _ := _userWrapper.LogoutURL(ctx, "/")
	templ := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	templ.Execute(w, map[string]string{"EMAIL": u.Email, "LOGOUT": url})
}

// HandleGetAllTodos ...
func (app *App) HandleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	app.commonJSONHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		return repo.ReadAllTodos(c)
	})
}

// HandleGetTodo ...
func (app *App) HandleGetTodo(w http.ResponseWriter, r *http.Request) {
	app.commonJSONHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		id, err := app.parseID(r)
		if err != nil {
			return nil, err
		}
		return repo.ReadTodo(c, id)
	})
}

// HandleDeleteDoneTodos ...
func (app *App) HandleDeleteDoneTodos(w http.ResponseWriter, r *http.Request) {
	app.commonJSONHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		return nil, repo.DeleteDoneTodos(c)
	})
}

// HandleDeleteTodo ...
func (app *App) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	app.commonJSONHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		id, err := app.parseID(r)
		if err != nil {
			return nil, err
		}
		return nil, repo.DeleteTodo(c, id)
	})
}

// HandlePutTodo ...
func (app *App) HandlePutTodo(w http.ResponseWriter, r *http.Request) {
	app.commonJSONHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		todo, err := app.json2Todo(r.Body)
		if err != nil {
			return nil, err
		}
		return repo.UpdateTodo(c, todo)
	})
}

// HandlePostTodo ...
func (app *App) HandlePostTodo(w http.ResponseWriter, r *http.Request) {
	app.commonJSONHandler(w, r, func(c context.Context, repo TodoRepository) (interface{}, error) {
		todo, err := app.json2Todo(r.Body)
		if err != nil {
			return nil, err
		}
		return repo.CreateTodo(c, todo)
	})
}

func (app *App) commonJSONHandler(
	w http.ResponseWriter, r *http.Request,
	fn func(c context.Context, repo TodoRepository) (interface{}, error)) {

	ctx := app.helper.Context(r)
	u := app.helper.CurrentUser(ctx)
	if u == nil {
		log.Errorf(ctx, "User is not login: %#v", u)
		http.Error(w, "Please login!!", http.StatusInternalServerError)
		return
	}

	repo := app.helper.TodoRepository(toHash(u.ID))
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

func (app *App) json2Todo(r io.ReadCloser) (*Todo, error) {
	defer r.Close()
	var todo Todo
	err := json.NewDecoder(r).Decode(&todo)
	return &todo, err
}

func (app *App) parseID(r *http.Request) (id int64, err error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["id"], 10, 64)
}

func toHash(text string) string {
	salt := os.Getenv("SALT")
	bytes := sha256.Sum256([]byte(text + salt))
	return hex.EncodeToString(bytes[:])
}
