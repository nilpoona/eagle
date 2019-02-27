package eagle

import (
	"net/http"
	"strings"
)

func CORS(config CORSHeaders) Middleware {
	allowOrigins := config.AllowOrigins
	if len(allowOrigins) == 0 {
		allowOrigins = DefaultCORSHeaders.AllowOrigins
	}

	allowMethods := config.AllowMethods
	if len(allowMethods) == 0 {
		allowMethods = DefaultCORSHeaders.AllowMethods
	}

	am := strings.Join(allowMethods, ",")
	ah := strings.Join(config.AllowHeaders, ",")
	eh := strings.Join(config.ExposeHeaders, ",")
	ma := config.MaxAge
	ac := config.AllowCredentials
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowOrigin := ""

			for _, o := range allowOrigins {
				if o == "*" || ac {
					allowOrigin = o
					break
				}

				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
			}

			if r.Method != http.MethodOptions {
				w.Header().Add("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				if ac {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				if eh != "" {
					w.Header().Set("Access-Control-Expose-Headers", eh)
				}

				next.ServeHTTP(w, r)
				return
			}

			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
			w.Header().Add("Access-Control-Allow-Methods", am)
			w.Header().Add("Access-Control-Allow-Origin", allowOrigin)
			if ma > 0 {
				w.Header().Set("Access-Control-Max-Age", string(ma))
			}

			if ac {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if ah != "" {
				w.Header().Set("Access-Control-Allow-Headers", ah)
			} else {
				h := r.Header.Get("Access-Control-Request-Headers")
				if h != "" {
					w.Header().Set("Access-Control-Allow-Headers", h)
				}
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}
