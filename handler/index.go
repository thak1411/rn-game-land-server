package handler

import "net/http"

type IndexHandler struct{}

func (h *IndexHandler) ServeHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./view/index.html")
}

func NewIndex() *IndexHandler {
	return &IndexHandler{}
}
