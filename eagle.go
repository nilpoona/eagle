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
	Options(w http.ResponseWriter, r *http.Request)
	Head(w http.ResponseWriter, r *http.Request)
	Trace(w http.ResponseWriter, r *http.Request)
}

type ResourceImpl struct{}

func (resource *ResourceImpl) Get(w http.ResponseWriter, r *http.Request)     {}
func (resource *ResourceImpl) Post(w http.ResponseWriter, r *http.Request)    {}
func (resource *ResourceImpl) Put(w http.ResponseWriter, r *http.Request)     {}
func (resource *ResourceImpl) Delete(w http.ResponseWriter, r *http.Request)  {}
func (resource *ResourceImpl) Patch(w http.ResponseWriter, r *http.Request)   {}
func (resource *ResourceImpl) Options(w http.ResponseWriter, r *http.Request) {}
func (resource *ResourceImpl) Head(w http.ResponseWriter, r *http.Request)    {}
func (resource *ResourceImpl) Trace(w http.ResponseWriter, r *http.Request)   {}

type Router interface {
	http.Handler
	SetResource(pattern string, resource Resource) error
	SetResourceWithMiddleware(pattern string, resource Resource, mw Middleware) error
}

func NewRouter() *Mux {
	m := make(map[string]*resourceInfo)
	return &Mux{resources: m}
}
