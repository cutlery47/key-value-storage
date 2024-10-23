package router

import (
	"fmt"
	"net/http"

	"github.com/cutlery47/key-value-storage/server/service"
)

type Router struct {
	mux  *http.ServeMux
	ctrl Controller
}

func New(service *service.Service) *Router {
	ctrl := Controller{
		service: service,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/add", ctrl.handleAdd)
	mux.HandleFunc("/api/v1/set", ctrl.handleSet)
	mux.HandleFunc("/api/v1/get", ctrl.handleGet)
	mux.HandleFunc("/api/v1/del", ctrl.handleDel)

	return &Router{
		ctrl: ctrl,
		mux:  mux,
	}
}

func (r Router) Handler() *http.ServeMux {
	return r.mux
}

type Controller struct {
	service    *service.Service
	errHandler errHandler
}

func (c Controller) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	expires_at := r.URL.Query().Get("expires_at")

	err := c.service.Add(key, value, expires_at)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
	}

	w.WriteHeader(http.StatusAccepted)

}

func (c Controller) handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	expires_at := r.URL.Query().Get("expires_at")

	err := c.service.Set(key, value, expires_at)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (c Controller) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")

	res, err := c.service.Get(key)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, res)

	// entry, err := c.storage.Read(key)
	// if err != nil {
	// 	if errors.Is(err, storage.ErrKeyNotFound) {
	// 		http.Error(w, err.Error(), http.StatusNotFound)
	// 	} else {
	// 		log.Println(err)
	// 		http.Error(w, "internal error", http.StatusInternalServerError)
	// 	}
	// 	return
	// }

	// jsonEntry, err := storage.ToJSON(*entry)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "internal error", http.StatusInternalServerError)
	// 	return
	// }

	// w.WriteHeader(http.StatusAccepted)
	// w.Write(jsonEntry)
}

func (c Controller) handleDel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")

	err := c.service.Delete(key)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
	}

	w.WriteHeader(http.StatusAccepted)
}
