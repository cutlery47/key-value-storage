package router

import (
	"fmt"
	"net/http"

	"github.com/cutlery47/key-value-storage/storage/internal/service"
)

type Router struct {
	mux  *http.ServeMux
	ctrl *Controller
}

func New(service *service.Service) *Router {
	ctrl := &Controller{
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

func (r *Router) Handler() *http.ServeMux {
	return r.mux
}

type Controller struct {
	service    *service.Service
	errHandler errHandler
}

func (c *Controller) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	key := r.PostFormValue("key")
	value := r.PostFormValue("value")
	expires_at := r.PostFormValue("expires_at")

	err := c.service.Add(key, value, expires_at)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	expires_at := r.Form.Get("expires_at")

	err := c.service.Set(key, value, expires_at)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")

	res, err := c.service.Get(key)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, res)
}

func (c *Controller) handleDel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")

	err := c.service.Delete(key)
	if err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
		return
	}

	w.WriteHeader(http.StatusOK)
}
