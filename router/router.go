package router

import "net/http"

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// chatRouter := NewChat()
	userRouter := NewUser()

	// mux.Handle("/ws/", http.StripPrefix("/ws", chatRouter))
	mux.Handle("/api/user/", http.StripPrefix("/api/user", userRouter))
	return mux
}
