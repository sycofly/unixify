package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/home/unixify/feature/ai_assisted/internal/models"
	"github.com/home/unixify/feature/ai_assisted/internal/utils"
)

// RegisterAccountRoutes registers the account routes
func (h *Handler) RegisterAccountRoutes(r chi.Router) {
	r.Get("/accounts", h.ListAccounts)
	r.Get("/accounts/search", h.SearchAccounts)
	r.Get("/accounts/new", h.NewAccount)
	r.Post("/accounts", h.CreateAccount)
	r.Get("/accounts/{id}", h.ViewAccount)
	r.Get("/accounts/{id}/edit", h.EditAccount)
	r.Put("/accounts/{id}", h.UpdateAccount)
	r.Delete("/accounts/{id}", h.DeleteAccount)
}

// ListAccounts handles GET /accounts
func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page := utils.GetPageFromQuery(r.URL.Query(), 1)
	perPage := utils.GetPerPageFromQuery(r.URL.Query(), 10)

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Get accounts
	accounts, total, err := h.accountService.FindAll(page, perPage, "", currentUserID)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	// Create pagination
	pagination := utils.NewPagination(page, perPage, total)

	// Render template
	h.render(w, "accounts/index.html", map[string]interface{}{
		"Accounts":   accounts,
		"Pagination": pagination,
		"Query":      "",
	})
}

// SearchAccounts handles GET /accounts/search
func (h *Handler) SearchAccounts(w http.ResponseWriter, r *http.Request) {
	// Get search query
	query := r.URL.Query().Get("q")

	// Get pagination parameters
	page := utils.GetPageFromQuery(r.URL.Query(), 1)
	perPage := utils.GetPerPageFromQuery(r.URL.Query(), 10)

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Get accounts
	accounts, total, err := h.accountService.FindAll(page, perPage, query, currentUserID)
	if err != nil {
		http.Error(w, "Failed to search accounts", http.StatusInternalServerError)
		return
	}

	// Create pagination
	pagination := utils.NewPagination(page, perPage, total)

	// Render template
	h.render(w, "accounts/index.html", map[string]interface{}{
		"Accounts":   accounts,
		"Pagination": pagination,
		"Query":      query,
	})
}

// NewAccount handles GET /accounts/new
func (h *Handler) NewAccount(w http.ResponseWriter, r *http.Request) {
	// Get all groups for selection
	groups, _, err := h.groupService.FindAll(1, 100, "", "system")
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	// Render template
	h.render(w, "accounts/form.html", map[string]interface{}{
		"Account":   models.Account{Active: true},
		"AllGroups": groups,
	})
}

// CreateAccount handles POST /accounts
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Create account
	account := models.Account{
		Base: models.Base{
			ID: uuid.New().String(),
		},
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		FullName: r.FormValue("full_name"),
		Active:   r.FormValue("active") == "true",
	}

	// Hash password
	password := r.FormValue("password")
	if password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		account.Password = hashedPassword
	}

	// Get selected groups
	groupIDs := r.Form["groups[]"]
	if len(groupIDs) > 0 {
		for _, groupID := range groupIDs {
			account.Groups = append(account.Groups, models.Group{
				Base: models.Base{
					ID: groupID,
				},
			})
		}
	}

	// Save account
	err = h.accountService.Create(&account, currentUserID)
	if err != nil {
		http.Error(w, "Failed to create account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to account list
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}

// ViewAccount handles GET /accounts/{id}
func (h *Handler) ViewAccount(w http.ResponseWriter, r *http.Request) {
	// Get account ID from URL
	id := chi.URLParam(r, "id")

	// Get account
	account, err := h.accountService.FindByID(id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Get creator and updater names
	createdByName := "System" // Placeholder, should look up username
	updatedByName := "System" // Placeholder, should look up username

	// Render template
	h.render(w, "accounts/view.html", map[string]interface{}{
		"Account":       account,
		"CreatedByName": createdByName,
		"UpdatedByName": updatedByName,
	})
}

// EditAccount handles GET /accounts/{id}/edit
func (h *Handler) EditAccount(w http.ResponseWriter, r *http.Request) {
	// Get account ID from URL
	id := chi.URLParam(r, "id")

	// Get account
	account, err := h.accountService.FindByID(id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Get all groups for selection
	groups, _, err := h.groupService.FindAll(1, 100, "", "system")
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	// Render template
	h.render(w, "accounts/form.html", map[string]interface{}{
		"Account":   account,
		"AllGroups": groups,
	})
}

// UpdateAccount handles PUT /accounts/{id}
func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	// Get account ID from URL
	id := chi.URLParam(r, "id")

	// Get account
	account, err := h.accountService.FindByID(id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
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

	// Update account
	account.Username = r.FormValue("username")
	account.Email = r.FormValue("email")
	account.FullName = r.FormValue("full_name")
	account.Active = r.FormValue("active") == "true"

	// Update password if provided
	password := r.FormValue("password")
	if password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		account.Password = hashedPassword
	}

	// Update groups
	account.Groups = []models.Group{}
	groupIDs := r.Form["groups[]"]
	if len(groupIDs) > 0 {
		for _, groupID := range groupIDs {
			account.Groups = append(account.Groups, models.Group{
				Base: models.Base{
					ID: groupID,
				},
			})
		}
	}

	// Save account
	err = h.accountService.Update(account, currentUserID)
	if err != nil {
		http.Error(w, "Failed to update account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to account view
	http.Redirect(w, r, "/accounts/"+id, http.StatusSeeOther)
}

// DeleteAccount handles DELETE /accounts/{id}
func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Get account ID from URL
	id := chi.URLParam(r, "id")

	// Get the current user ID (from session/auth)
	currentUserID := "system" // Placeholder, should come from auth

	// Delete account
	err := h.accountService.Delete(id, currentUserID)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	// For HTMX requests, return empty response (success)
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Redirect to account list
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}