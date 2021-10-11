package middleware

import (
	"net/http"
)

func Middleware(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// inject code //
		f(w, r)
	}
}
