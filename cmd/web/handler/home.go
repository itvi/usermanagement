package handler

import "net/http"

// HomeHandler ...
type HomeHandler struct{}

func (h *HomeHandler) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, "./ui/html/index.html", nil, "")
	}
}
