package handlers

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
)

// RegisterHomeRoutes registers the home routes
func (h *Handler) RegisterHomeRoutes(r chi.Router) {
	r.Get("/", h.Home)
}

// Home handles GET /
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Get total accounts
	_, totalAccounts, err := h.accountService.FindAll(1, 1, "", "system")
	if err != nil {
		totalAccounts = 0
	}

	// Get total groups
	_, totalGroups, err := h.groupService.FindAll(1, 1, "", "system")
	if err != nil {
		totalGroups = 0
	}

	// Render template
	h.render(w, "home.html", map[string]interface{}{
		"TotalAccounts": totalAccounts,
		"TotalGroups":   totalGroups,
	})
}