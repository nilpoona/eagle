/*
	Package eagle provides functions required to build REST API.

    It is things resource get API

	```
	type ThingsResource struct {
		eagle.ResourceImpl
	}

	func (tr *ThingsResource) Get(w http.ResponseWriter, r *http.Request) {
		type response {
			message string
		}
		resp := response{message: "something"}
		eagle.RenderJSON(w, http.StatusOK, &resp)
	}

	func main() {
		tr := &ThingsResource{}
		router := eagle.NewRouter()
		err = router.SetResource("/things", tr)
		if err != nil {
			panic(err)
		}
		if err := http.ListenAndServe(":8080", router); err != nil {
			panic(err)
		}
	}
	```
*/

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
	Use(middleware Middleware)
}

func NewRouter() *Mux {
	m := make(map[string]*resourceInfo)
	return &Mux{resources: m}
}
