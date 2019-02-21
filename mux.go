package eagle

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type resourceInfoMap map[string]*resourceInfo

const pathParamPrefix = "EaglePathParam:"

type resourceInfo struct {
	resource   Resource
	isRegExp   bool
	middleware Middleware
}

type Mux struct {
	handler   http.Handler
	resources resourceInfoMap
}

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

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

func (mux *Mux) SetResourceWithMiddleware(pattern string, resource Resource, mw Middleware) error {
	p, isRegExp, err := genMatchPattern(pattern)
	if err != nil {
		return err
	}

	ri := &resourceInfo{
		resource:   resource,
		isRegExp:   isRegExp,
		middleware: mw,
	}

	mux.resources[p] = ri
	return nil
}

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
	path := r.URL.Path

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
	default:
		h = handleMethodNotAllowed
	}

	if ri.middleware != nil {
		h = ri.middleware(h)
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
