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
	resource Resource
	isRegExp bool
}

type Mux struct {
	handler   http.Handler
	resources resourceInfoMap
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

func findResourceByRequestPath(resources resourceInfoMap, path string) (Resource, map[string]string) {
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
			return ri.resource, params
		} else {
			if path == pattern {
				return ri.resource, params
			}
		}
	}
	return nil, params
}

func (mux *Mux) SetResource(pattern string, resource Resource) error {
	p, isRegExp, err := genMatchPattern(pattern)
	if err != nil {
		return err
	}

	mux.resources[p] = &resourceInfo{
		resource: resource,
		isRegExp: isRegExp,
	}

	return nil
}

func (mux *Mux) handle(r *http.Request) (http.HandlerFunc, map[string]string) {
	method := r.Method
	path := r.URL.Path

	resource, params := findResourceByRequestPath(mux.resources, path)
	if resource == nil {
		// TODO: 404
	}

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
