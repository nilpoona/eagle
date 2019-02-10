package eagle

import (
	"net/http"
)

type Mux struct {
	handler   http.Handler
	resources map[string]Resource
}

func (mux *Mux) SetResource(pattern string, resource Resource) {
	mux.resources[pattern] = resource
}

func (mux *Mux) handle(r *http.Request) http.HandlerFunc {
	method := r.Method
	path := r.URL.Path

	resource := mux.resources[path]
	var h http.HandlerFunc
	switch method {
	case http.MethodPost:
		h = http.HandlerFunc(resource.Post)
	case http.MethodGet:
		h = http.HandlerFunc(resource.Get)
	case http.MethodPut:
		h = http.HandlerFunc(resource.Put)
	case http.MethodDelete:
		h = http.HandlerFunc(resource.Delete)
	case http.MethodPatch:
		h = http.HandlerFunc(resource.Patch)
	}
	// TODO: method not allowed

	return h
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := mux.handle(r)
	h.ServeHTTP(w, r)
}
