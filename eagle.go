package eagle

import (
	"net/http"
)

type Resource interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Head(w http.ResponseWriter, r *http.Request)
}

type RestResource struct{}

func (resource *RestResource) Get(w http.ResponseWriter, r *http.Request)    {}
func (resource *RestResource) Post(w http.ResponseWriter, r *http.Request)   {}
func (resource *RestResource) Put(w http.ResponseWriter, r *http.Request)    {}
func (resource *RestResource) Delete(w http.ResponseWriter, r *http.Request) {}
func (resource *RestResource) Patch(w http.ResponseWriter, r *http.Request)  {}
func (resource *RestResource) Head(w http.ResponseWriter, r *http.Request)   {}

type Router interface {
	http.Handler
	SetResource(pattern string, resource Resource)
}

func NewRouter() *Mux {
	m := make(map[string]Resource)
	return &Mux{resources: m}
}
