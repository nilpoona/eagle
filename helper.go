package eagle

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	DefaultCORSHeaders = CORSHeaders{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
)

func PathParam(r *http.Request, k string) string {
	if v := r.Context().Value(pathParamKey(k)); v != nil {
		return v.(string)
	}

	return ""
}

func bindFormData(formData map[string][]string, v interface{}) error {
	typ := reflect.TypeOf(v).Elem()
	val := reflect.ValueOf(v).Elem()

	if typ.Kind() != reflect.Struct {
		return errors.New("must be a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tagValue := f.Tag.Get("form")
		if tagValue == "" {
			continue
		}

		formValue := formData[tagValue]
		if len(formValue) == 0 {
			continue
		}

		if formValue[0] == "" {
			return fmt.Errorf("form value does not exist key: %s", tagValue)
		}

		field := val.FieldByName(f.Name)
		if field.Interface() == nil {
			continue
		}

		kind := field.Type().Kind()
		isPtr := false
		if kind == reflect.Ptr {
			kind = field.Type().Elem().Kind()
			isPtr = true
		}

		fv := formValue[0]
		switch kind {
		case reflect.Int:
			n, err := strconv.Atoi(fv)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowInt(int64(n))
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowInt(int64(n))
			}
			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Int8:
			n, err := strconv.ParseInt(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowInt(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowInt(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := int8(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Int16:
			n, err := strconv.ParseInt(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowInt(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowInt(n)
			}
			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := int16(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Int32:
			n, err := strconv.ParseInt(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowInt(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowInt(n)
			}
			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := int32(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Int64:
			n, err := strconv.ParseInt(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowInt(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowInt(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Uint:
			n, err := strconv.ParseUint(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowUint(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowUint(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := uint(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Uint8:
			n, err := strconv.ParseUint(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowUint(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowUint(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := uint8(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Uint16:
			n, err := strconv.ParseUint(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowUint(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowUint(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := uint16(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Uint32:
			n, err := strconv.ParseUint(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowUint(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowUint(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := uint32(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Uint64:
			n, err := strconv.ParseUint(fv, 10, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowUint(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowUint(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := uint64(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Uintptr:
			n, err := strconv.ParseUint(fv, 10, 64)
			if err != nil || reflect.Zero(field.Type()).OverflowUint(n) {
				return fmt.Errorf("bind form data failed: %s", err)
			}
			field.Set(reflect.ValueOf(n).Convert(field.Type()))
		case reflect.Float32:
			n, err := strconv.ParseFloat(fv, 64)
			isOverflow := false
			if isPtr {
				isOverflow = reflect.Zero(field.Type().Elem()).OverflowFloat(n)
			} else {
				isOverflow = reflect.Zero(field.Type()).OverflowFloat(n)
			}

			if err != nil || isOverflow {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				n := float32(n)
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.Float64:
			n, err := strconv.ParseFloat(fv, 64)
			if err != nil || reflect.Zero(field.Type()).OverflowFloat(n) {
				return fmt.Errorf("bind form data failed: %s", err)
			}

			if isPtr {
				field.Set(reflect.ValueOf(&n))
			} else {
				field.Set(reflect.ValueOf(n).Convert(field.Type()))
			}
		case reflect.String:
			if isPtr {
				str := formValue[0]
				field.Set(reflect.ValueOf(&str))
			} else {
				field.Set(reflect.ValueOf(formValue[0]))
			}
		case reflect.Bool:
			b, err := strconv.ParseBool(fv)
			if err != nil {
				return fmt.Errorf("bind form data failed: %s", err)
			}
			field.Set(reflect.ValueOf(b).Convert(field.Type()))
		default:
			fmt.Println(field.Type().Kind())
		}
	}

	return nil
}

func Bind(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		enc := json.NewDecoder(r.Body)
		return enc.Decode(v)
	case strings.HasPrefix(contentType, "application/xml"):
		enc := xml.NewDecoder(r.Body)
		return enc.Decode(v)
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"), strings.HasPrefix(contentType, "multipart/form-data"):
		err := r.ParseForm()
		if err != nil {
			return err
		}
		return bindFormData(r.Form, v)
	default:
		return errors.New("unsupported media type")
	}
}

type CORSHeaders struct {
	AllowOrigins     []string `json:"Access-Control-Allow-Origin"`
	AllowMethods     []string `json:"Access-Control-Allow-Methods"`
	AllowHeaders     []string `json:"Access-Control-Allow-Headers"`
	ExposeHeaders    []string `json:"Access-Control-Expose-Headers"`
	MaxAge           uint     `json:"Access-Control-Max-Age"`
	AllowCredentials bool     `json:"Access-Control-Allow-Credentials"`
}

func RenderJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}
