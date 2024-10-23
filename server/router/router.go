package router

import (
	"net/http"

	"github.com/cutlery47/key-value-storage/server/storage"
)

type Router struct {
	mux  *http.ServeMux
	ctrl Controller
}

func New(storage storage.Storage) *Router {
	ctrl := Controller{
		storage: storage,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/add", ctrl.HandleAdd)
	mux.HandleFunc("/api/v1/set", ctrl.HandleSet)
	mux.HandleFunc("/api/v1/get", ctrl.HandleGet)
	mux.HandleFunc("/api/v1/del", ctrl.HandleDel)

	return &Router{
		ctrl: ctrl,
		mux:  mux,
	}
}

func (r Router) Handler() *http.ServeMux {
	return r.mux
}

type Controller struct {
	storage storage.Storage
}

func (c Controller) HandleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (c Controller) HandleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (c Controller) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")

	entry, err := c.storage.Read(key)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	jsonEntry, err := storage.ToJSON(*entry)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonEntry)
}

func (c Controller) HandleDel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
