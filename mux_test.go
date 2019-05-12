package eagle

import (
	"errors"
	"reflect"
	"testing"
)

func isSameErrorMessage(err, err2 error) bool {
	if err == nil && err2 == nil {
		return true
	}

	return err.Error() == err2.Error()
}

func TestFindResourceByRequestPath(t *testing.T) {
	type resource struct {
		ResourceImpl
	}
	type want struct {
		resourceInfo *resourceInfo
		params       map[string]string
	}

	tests := []struct {
		name            string
		path            string
		resourceInfoMap func() resourceInfoMap
		wants           want
	}{
		{
			name: "test1",
			path: "/users/2",
			resourceInfoMap: func() resourceInfoMap {
				r := &resource{}
				ri := &resourceInfo{
					resource: r,
					isRegExp: true,
				}
				return resourceInfoMap(map[string]*resourceInfo{
					"/users/(?P<id>[0-9]+)": ri,
				})
			},
			wants: want{
				resourceInfo: &resourceInfo{
					resource: &resource{},
					isRegExp: true,
				},
				params: map[string]string{
					"id": "2",
				},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			r, params := findResourceInfoByRequestPath(td.resourceInfoMap(), td.path)
			if !reflect.DeepEqual(r, td.wants.resourceInfo) {
				t.Errorf("findResourceInfoByRequestPath failed result %+v, expected: %+v", r, td.wants.resourceInfo)
			}

			if !reflect.DeepEqual(params, td.wants.params) {
				t.Errorf("findResourceInfoByRequestPath failed result %+v, expected: %+v", params, td.wants.params)
			}
		})
	}
}

func TestGenMatchPattern(t *testing.T) {
	type want struct {
		p        string
		isRegExp bool
		err      error
	}
	tests := []struct {
		name    string
		pattern string
		want    want
	}{
		{
			name:    "Regular expression match pattern can be obtained",
			pattern: "/users/{id:[0-9]+}",
			want: want{
				p:        "/users/(?P<id>[0-9]+)",
				isRegExp: true,
				err:      nil,
			},
		},
		{
			name:    "You can get the pattern passed as argument",
			pattern: "/users",
			want: want{
				p:        "/users",
				isRegExp: false,
				err:      nil,
			},
		},
		{
			name:    "{There is no error so",
			pattern: "/users/id:[0-9]+}",
			want: want{
				p:        "/users/id:[0-9]+}",
				isRegExp: false,
				err:      errors.New("The path parameter's { is not set"),
			},
		},
		{
			name:    "}There is no error so",
			pattern: "/users/{id:[0-9]+",
			want: want{
				p:        "/users/{id:[0-9]+",
				isRegExp: false,
				err:      errors.New("The path parameter's } is not set"),
			},
		},
		{
			name:    "Error because delimiter is not set",
			pattern: "/users/{id[0-9]+}",
			want: want{
				p:        "/users/{id[0-9]+}",
				isRegExp: false,
				err:      errors.New("delimiters are not set in path parameters"),
			},
		},
		{
			name:    "Error because it contains prohibited character string",
			pattern: "/users/{id:[0-9]{2}}",
			want: want{
				p:        "/users/{id:[0-9]{2}}",
				isRegExp: false,
				err:      errors.New("`{}` Can not be used as a regular expression pattern"),
			},
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			pattern, isRegExp, err := genMatchPattern(td.pattern)
			if pattern != td.want.p {
				t.Errorf("Error genMatchPattern result: %s, expected: %s", pattern, td.want.p)
			}

			if isRegExp != td.want.isRegExp {
				t.Errorf("Error genMatchPattern result: %T, expected: %T", isRegExp, td.want.isRegExp)
			}

			if !isSameErrorMessage(err, td.want.err) {
				t.Errorf("Error genMatchPattern result: %s, expected: %s", err, td.want.err)
			}
		})
	}
}
