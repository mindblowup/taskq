package main

import "net/http"

type middleware func(next http.HandlerFunc) http.HandlerFunc

func httpMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

func tokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := r.URL.Query().Get("token")
		if *flagSecretToken != t {
			responseError(w, map[string]string{"token": "invalid token"})
			return
		}
		next.ServeHTTP(w, r)
	}
}
