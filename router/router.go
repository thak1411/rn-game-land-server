package router

import "net/http"

func New() *http.ServeMux {
	mux := http.NewServeMux()

	indexRouter := NewIndex()
	userRouter := NewUser()

	mux.Handle("/", indexRouter)
	mux.Handle("/user/", http.StripPrefix("/user", userRouter))
	return mux
}
