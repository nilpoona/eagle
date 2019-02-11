package eagle

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type resourceInfo struct {
	resource Resource
	params   map[string]string
	pattern  string
}

type Mux struct {
	handler   http.Handler
	resources map[string]*resourceInfo
}

// /users/{id:[0-9]+}
func genMatchPattern(pattern string) (string, bool, error) {
	p := pattern
	isRegExp := false

	for {
		si := strings.Index(p, "{")
		ei := strings.Index(p, "}")

		if si != -1 && ei == -1 {
			return p, false, errors.New("The path parameter's } is not set")
		}

		if si == -1 && ei != -1 {
			return p, false, errors.New("The path parameter's { is not set")
		}

		if si == -1 && ei == -1 {
			break
		}

		param := p[si+1 : ei]

		separatorIndex := strings.Index(param, ":")
		if separatorIndex == -1 {
			return pattern, false, errors.New("delimiters are not set in path parameters")
		}

		if strings.Index(param, "}") != -1 || strings.Index(param, "{") != -1 {
			return pattern, false, errors.New("`{}` Can not be used as a regular expression pattern")
		}

		isRegExp = true

		label := param[0:separatorIndex]
		value := param[separatorIndex+1:]

		p = fmt.Sprintf("%s(?P<%s>%s)%s", p[0:si], label, value, p[ei+1:])
	}

	return p, isRegExp, nil
}

func (mux *Mux) SetResource(pattern string, resource Resource) {
	mux.resources[pattern] = &resourceInfo{
		resource: resource,
	}
}

func (mux *Mux) handle(r *http.Request) http.HandlerFunc {
	method := r.Method
	path := r.URL.Path

	resource := mux.resources[path].resource
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
