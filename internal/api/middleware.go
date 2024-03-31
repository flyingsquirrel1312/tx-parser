package api

import "net/http"

func MethodMiddleware(handler http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func Post(handler http.HandlerFunc) http.HandlerFunc {
	return MethodMiddleware(handler, http.MethodPost)
}

func Get(handler http.HandlerFunc) http.HandlerFunc {
	return MethodMiddleware(handler, http.MethodGet)
}
