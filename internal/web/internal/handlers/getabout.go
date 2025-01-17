package handlers

import (
	"net/http"

	"github.com/akiidjk/styx/internal/web/internal/templates"
)

type AboutHandLer struct{}

func NewAboutHandler() *AboutHandLer {
	return &AboutHandLer{}
}

func (h *AboutHandLer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.About()
	err := templates.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
