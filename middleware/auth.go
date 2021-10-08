package middleware

import (
	"net/http"
)

func AuthToken(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return f
}
