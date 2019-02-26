package eagle

import (
	"encoding/json"
	"net/http"
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

func Bind(r *http.Request, v interface{}) error {
	enc := json.NewDecoder(r.Body)
	return enc.Decode(v)
}

type CORSHeaders struct {
	AllowOrigins     []string `json:"Access-Control-Allow-Origin"`
	AllowMethods     []string `json:"Access-Control-Allow-Methods"`
	AllowHeaders     []string `json:"Access-Control-Allow-Headers"`
	ExposeHeaders    []string `json:"Access-Control-Expose-Headers"`
	MaxAge           uint     `json:"Access-Control-Max-Age"`
	AllowCredentials bool     `json:"Access-Control-Allow-Credentials"`
}

func SetCORSHeaders(ch CORSHeaders) error {

	return nil
}

func RenderJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}
