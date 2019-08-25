package eagle

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var _ Router = (*Mux)(nil)

type resourceInfoMap map[string]*resourceInfo

const pathParamPrefix = "EaglePathParam:"

type resourceInfo struct {
	resource   Resource
	isRegExp   bool
	middleware Middleware
}

type Mux struct {
	handler     http.Handler
	resources   resourceInfoMap
	middlewares []Middleware
}

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func genMatchPattern(pattern string) (string, bool, error) {
	p := strings.TrimRight(pattern, "/")
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

func findResourceInfoByRequestPath(resources resourceInfoMap, path string) (*resourceInfo, map[string]string) {
	params := make(map[string]string)
	for pattern, ri := range resources {
		if ri.isRegExp {
			re := regexp.MustCompile(pattern)
			match := re.FindSubmatch([]byte(path))
			if len(match) == 0 {
				continue
			}
			for i, name := range re.SubexpNames() {
				if i != 0 && name != "" {
					params[name] = string(match[i])
				}
			}
			return ri, params
		} else {
			if path == pattern {
				return ri, params
			}
		}
	}
	return nil, params
}

// Middleware Add Middleware to mux. The added middleware applies to all resources
func (mux *Mux) Use(m Middleware) {
	mux.middlewares = append(mux.middlewares, m)
}

// SetResource
func (mux *Mux) SetResource(pattern string, resource Resource) error {
	p, isRegExp, err := genMatchPattern(pattern)
	if err != nil {
		return err
	}

	ri := &resourceInfo{
		resource:   resource,
		isRegExp:   isRegExp,
		middleware: nil,
	}

	mux.resources[p] = ri
	return nil
}

func (mux *Mux) handle(r *http.Request) (http.HandlerFunc, map[string]string) {
	method := r.Method
	path := strings.TrimRight(r.URL.Path, "/")

	ri, params := findResourceInfoByRequestPath(mux.resources, path)
	if ri == nil {
		return handleNotFound, make(map[string]string)
	}

	resource := ri.resource

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
	case http.MethodOptions:
		h = http.HandlerFunc(resource.Options)
	case http.MethodTrace:
		h = http.HandlerFunc(resource.Trace)
	default:
		h = handleMethodNotAllowed
	}

	if len(mux.middlewares) > 0 {
		middleware := ChainMiddleware(mux.middlewares...)
		h = middleware(h)
	}

	return h, params
}

func pathParamKey(k string) string {
	return fmt.Sprintf("%s%s", pathParamPrefix, k)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, params := mux.handle(r)
	for key, val := range params {
		r = r.WithContext(context.WithValue(r.Context(), pathParamKey(key), val))
	}
	h.ServeHTTP(w, r)
}
