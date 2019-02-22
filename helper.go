package eagle

import (
	"encoding/json"
	"net/http"
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

func RenderJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}
