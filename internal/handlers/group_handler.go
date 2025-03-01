package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/home/unixify/feature/ai_assisted/internal/models"
	"github.com/home/unixify/feature/ai_assisted/internal/utils"
)

// RegisterGroupRoutes registers the group routes
func (h *Handler) RegisterGroupRoutes(r chi.Router) {
	r.Get("/groups", h.ListGroups)
	r.Get("/groups/search", h.SearchGroups)
	r.Get("/groups/new", h.NewGroup)
	r.Post("/groups", h.CreateGroup)
	r.Get("/groups/{id}", h.ViewGroup)
	r.Get("/groups/{id}/edit", h.EditGroup)
	r.Put("/groups/{id}", h.UpdateGroup)
	r.Delete("/groups/{id}", h.DeleteGroup)
}

// ListGroups handles GET /groups
func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page := utils.GetPageFromQuery(r.URL.Query(), 1)
	perPage := utils.GetPerPageFromQuery(r.URL.Query(), 10)

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Get groups
	groups, total, err := h.groupService.FindAll(page, perPage, "", currentUserID)
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	// Create pagination
	pagination := utils.NewPagination(page, perPage, total)

	// Render template
	h.render(w, "groups/index.html", map[string]interface{}{
		"Groups":     groups,
		"Pagination": pagination,
		"Query":      "",
	})
}

// SearchGroups handles GET /groups/search
func (h *Handler) SearchGroups(w http.ResponseWriter, r *http.Request) {
	// Get search query
	query := r.URL.Query().Get("q")

	// Get pagination parameters
	page := utils.GetPageFromQuery(r.URL.Query(), 1)
	perPage := utils.GetPerPageFromQuery(r.URL.Query(), 10)

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Get groups
	groups, total, err := h.groupService.FindAll(page, perPage, query, currentUserID)
	if err != nil {
		http.Error(w, "Failed to search groups", http.StatusInternalServerError)
		return
	}

	// Create pagination
	pagination := utils.NewPagination(page, perPage, total)

	// Render template
	h.render(w, "groups/index.html", map[string]interface{}{
		"Groups":     groups,
		"Pagination": pagination,
		"Query":      query,
	})
}

// NewGroup handles GET /groups/new
func (h *Handler) NewGroup(w http.ResponseWriter, r *http.Request) {
	// Get all accounts for selection
	accounts, _, err := h.accountService.FindAll(1, 100, "", "system")
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	// Render template
	h.render(w, "groups/form.html", map[string]interface{}{
		"Group":       models.Group{},
		"AllAccounts": accounts,
	})
}

// CreateGroup handles POST /groups
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Create group
	group := models.Group{
		Base: models.Base{
			ID: uuid.New().String(),
		},
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	// Get selected members
	memberIDs := r.Form["members[]"]
	if len(memberIDs) > 0 {
		for _, memberID := range memberIDs {
			group.Members = append(group.Members, models.Account{
				Base: models.Base{
					ID: memberID,
				},
			})
		}
	}

	// Save group
	err = h.groupService.Create(&group, currentUserID)
	if err != nil {
		http.Error(w, "Failed to create group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to group list
	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}

// ViewGroup handles GET /groups/{id}
func (h *Handler) ViewGroup(w http.ResponseWriter, r *http.Request) {
	// Get group ID from URL
	id := chi.URLParam(r, "id")

	// Get group
	group, err := h.groupService.FindByID(id)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Get creator and updater names
	createdByName := "System" // Placeholder, should look up username
	updatedByName := "System" // Placeholder, should look up username

	// Render template
	h.render(w, "groups/view.html", map[string]interface{}{
		"Group":         group,
		"CreatedByName": createdByName,
		"UpdatedByName": updatedByName,
	})
}

// EditGroup handles GET /groups/{id}/edit
func (h *Handler) EditGroup(w http.ResponseWriter, r *http.Request) {
	// Get group ID from URL
	id := chi.URLParam(r, "id")

	// Get group
	group, err := h.groupService.FindByID(id)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Get all accounts for selection
	accounts, _, err := h.accountService.FindAll(1, 100, "", "system")
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	// Render template
	h.render(w, "groups/form.html", map[string]interface{}{
		"Group":       group,
		"AllAccounts": accounts,
	})
}

// UpdateGroup handles PUT /groups/{id}
func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	// Get group ID from URL
	id := chi.URLParam(r, "id")

	// Get group
	group, err := h.groupService.FindByID(id)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Parse form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Update group
	group.Name = r.FormValue("name")
	group.Description = r.FormValue("description")

	// Update members
	group.Members = []models.Account{}
	memberIDs := r.Form["members[]"]
	if len(memberIDs) > 0 {
		for _, memberID := range memberIDs {
			group.Members = append(group.Members, models.Account{
				Base: models.Base{
					ID: memberID,
				},
			})
		}
	}

	// Save group
	err = h.groupService.Update(group, currentUserID)
	if err != nil {
		http.Error(w, "Failed to update group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to group view
	http.Redirect(w, r, "/groups/"+id, http.StatusSeeOther)
}

// DeleteGroup handles DELETE /groups/{id}
func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	// Get group ID from URL
	id := chi.URLParam(r, "id")

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Delete group
	err := h.groupService.Delete(id, currentUserID)
	if err != nil {
		http.Error(w, "Failed to delete group", http.StatusInternalServerError)
		return
	}

	// For HTMX requests, return empty response (success)
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Redirect to group list
	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}