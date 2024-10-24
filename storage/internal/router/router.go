package router

import (
	"fmt"
	"net/http"

	"github.com/cutlery47/key-value-storage/storage/internal/service"
	"github.com/sirupsen/logrus"
)

// responsible for routing http-request
// to the specific handler
type Router struct {
	ctrl *Controller

	mux *http.ServeMux
	log *logrus.Logger
}

func New(service *service.Service, infoLog, errLog *logrus.Logger) *Router {
	errHandler := errHandler{
		errLog: errLog,
	}

	ctrl := &Controller{
		service:    service,
		errHandler: errHandler,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/add", ctrl.handleAdd)
	mux.HandleFunc("/api/v1/set", ctrl.handleSet)
	mux.HandleFunc("/api/v1/get", ctrl.handleGet)
	mux.HandleFunc("/api/v1/del", ctrl.handleDel)

	return &Router{
		ctrl: ctrl,
		mux:  mux,
		log:  infoLog,
	}
}

func (r *Router) Handler() http.Handler {
	return WithLogging(r.mux, r.log)
}

// responsible for parsing and packing http-requests/responses
// passes received data down to the service layer
type Controller struct {
	service    *service.Service
	errHandler errHandler
}

func (c *Controller) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key, value, expiresAt := c.parsePostForm(r)

	if err := c.service.Add(key, value, expiresAt); err != nil {
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

	key, value, expiresAt := c.parsePostForm(r)

	if err := c.service.Set(key, value, expiresAt); err != nil {
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

	if err := c.service.Delete(key); err != nil {
		status, msg := c.errHandler.Handle(err)
		http.Error(w, msg, status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c Controller) parsePostForm(r *http.Request) (key, value, expiresAt string) {
	key = r.PostFormValue("key")
	value = r.PostFormValue("value")
	expiresAt = r.PostFormValue("expires_at")

	return key, value, expiresAt
}
