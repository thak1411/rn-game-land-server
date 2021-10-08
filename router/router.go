package router

import "net/http"

func New() *http.ServeMux {
	mux := http.NewServeMux()

	userRouter := NewUser()

	mux.Handle("/api/user/", http.StripPrefix("/api/user", userRouter))
	return mux
}
