package eagle

import "net/http"

func PathParam(r *http.Request, k string) string {
	if v := r.Context().Value(pathParamKey(k)); v != nil {
		return v.(string)
	}

	return ""
}
