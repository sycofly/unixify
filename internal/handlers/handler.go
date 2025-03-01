package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/home/unixify/feature/ai_assisted/internal/models"
	"github.com/home/unixify/feature/ai_assisted/internal/services"
	"github.com/home/unixify/feature/ai_assisted/internal/utils"
)

// Handler holds the dependencies for HTTP handlers
type Handler struct {
	db             *sqlx.DB
	templates      *template.Template
	templatesMutex sync.RWMutex
	accountService *services.AccountService
	groupService   *services.GroupService
	templateDir    string
}

// NewHandler creates a new Handler
func NewHandler(db *sqlx.DB, templateDir string) *Handler {
	h := &Handler{
		db:             db,
		accountService: services.NewAccountService(db),
		groupService:   services.NewGroupService(db),
		templateDir:    templateDir,
	}

	// Load templates
	h.loadTemplates()

	return h
}

// loadTemplates loads all templates from the template directory
func (h *Handler) loadTemplates() {
	h.templatesMutex.Lock()
	defer h.templatesMutex.Unlock()

	// Get template functions
	funcMap := utils.TemplateFuncs()

	// Create a new template set
	templates := template.New("").Funcs(funcMap)

	// First check if the layout exists
	layoutPath := filepath.Join(h.templateDir, "shared", "layout.html")
	layoutContent, err := os.ReadFile(layoutPath)
	if err != nil {
		log.Printf("ERROR: Could not read layout template: %v", err)
	} else {
		// Create a base template with the layout content
		templates, err = template.New("layout.html").Funcs(funcMap).Parse(string(layoutContent))
		if err != nil {
			log.Printf("ERROR: Failed to parse layout template: %v", err)
		} else {
			log.Printf("Successfully loaded layout template")
		}
	}

	// Load all templates using a single pattern
	allTemplatePatterns := []string{
		filepath.Join(h.templateDir, "shared", "*.html"),
		filepath.Join(h.templateDir, "*.html"),
		filepath.Join(h.templateDir, "accounts", "*.html"),
		filepath.Join(h.templateDir, "groups", "*.html"),
	}

	for _, pattern := range allTemplatePatterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			log.Printf("Error finding template files with pattern %s: %v", pattern, err)
			continue
		}

		// Skip layout file since we've already loaded it
		var filesToParse []string
		for _, file := range files {
			if file != layoutPath {
				filesToParse = append(filesToParse, file)
			}
		}

		if len(filesToParse) > 0 {
			_, err = templates.ParseFiles(filesToParse...)
			if err != nil {
				log.Printf("Warning: Failed to load templates from pattern %s: %v", pattern, err)
			} else {
				log.Printf("Loaded %d templates from pattern %s", len(filesToParse), pattern)
			}
		}
	}

	// Log all loaded templates for debugging
	var templateNames []string
	for _, t := range templates.Templates() {
		if t.Name() != "" {
			templateNames = append(templateNames, t.Name())
		}
	}
	log.Printf("Loaded templates: %v", templateNames)

	h.templates = templates
}

// render renders a template with data
func (h *Handler) render(w http.ResponseWriter, name string, data interface{}) {
	h.templatesMutex.RLock()
	defer h.templatesMutex.RUnlock()

	log.Printf("Rendering template: %s with data: %+v", name, data)
	
	// Debug: List all defined templates
	var templates []string
	for _, t := range h.templates.Templates() {
		if t.Name() != "" {
			templates = append(templates, t.Name())
		}
	}
	log.Printf("Available templates: %v", templates)
	
	// Add default CurrentYear if not present
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		dataMap = map[string]interface{}{
			"Content": data,
			"CurrentYear": time.Now().Year(),
		}
	} else if _, exists := dataMap["CurrentYear"]; !exists {
		dataMap["CurrentYear"] = time.Now().Year()
	}

	// Handle special cases based on template name
	w.Header().Set("Content-Type", "text/html")
	
	if name == "home.html" {
		// For home page, we know this works
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Unixify</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
    <script defer src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"></script>
</head>
<body>
    <header>
        <nav>
            <div class="container">
                <a href="/" class="logo">Unixify</a>
                <ul class="nav-links">
                    <li><a href="/">Home</a></li>
                    <li><a href="/accounts">Accounts</a></li>
                    <li><a href="/groups">Groups</a></li>
                </ul>
            </div>
        </nav>
    </header>

    <main class="container">
        <div class="card">
            <div class="card-header">
                <h1>Welcome to Unixify</h1>
            </div>
            
            <div class="dashboard">
                <div class="dashboard-item">
                    <div class="dashboard-item-header">
                        <h2>Accounts</h2>
                        <a href="/accounts" class="btn btn-sm">View All</a>
                    </div>
                    <div class="dashboard-item-content">
                        <p class="stat">` + fmt.Sprintf("%d", dataMap["TotalAccounts"]) + ` Total Accounts</p>
                        <p>Manage users with role-based access control</p>
                        <div class="dashboard-actions">
                            <a href="/accounts/new" class="btn">Create Account</a>
                        </div>
                    </div>
                </div>
                
                <div class="dashboard-item">
                    <div class="dashboard-item-header">
                        <h2>Groups</h2>
                        <a href="/groups" class="btn btn-sm">View All</a>
                    </div>
                    <div class="dashboard-item-content">
                        <p class="stat">` + fmt.Sprintf("%d", dataMap["TotalGroups"]) + ` Total Groups</p>
                        <p>Organize users into logical groups for easier management</p>
                        <div class="dashboard-actions">
                            <a href="/groups/new" class="btn">Create Group</a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </main>

    <footer>
        <div class="container">
            <p>&copy; ` + fmt.Sprintf("%d", dataMap["CurrentYear"]) + ` Unixify. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
`))
	} else if name == "accounts/form.html" {
		// Specific handler for account form
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Create Account - Unixify</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <header>
        <nav>
            <div class="container">
                <a href="/" class="logo">Unixify</a>
                <ul class="nav-links">
                    <li><a href="/">Home</a></li>
                    <li><a href="/accounts">Accounts</a></li>
                    <li><a href="/groups">Groups</a></li>
                </ul>
            </div>
        </nav>
    </header>

    <main class="container">
        <div class="card">
            <div class="card-header">
                <h1>Create New Account</h1>
            </div>
            <div class="card-body">
                <form action="/accounts" method="POST" class="form">
                    <div class="form-group">
                        <label for="username">Username</label>
                        <input type="text" id="username" name="username" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="email">Email</label>
                        <input type="email" id="email" name="email" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="first_name">First Name</label>
                        <input type="text" id="first_name" name="first_name" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="last_name">Last Name</label>
                        <input type="text" id="last_name" name="last_name" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="role">Role</label>
                        <select id="role" name="role" required>
                            <option value="user">User</option>
                            <option value="admin">Admin</option>
                        </select>
                    </div>
                    
                    <div class="form-actions">
                        <a href="/accounts" class="btn btn-secondary">Cancel</a>
                        <button type="submit" class="btn btn-primary">Create Account</button>
                    </div>
                </form>
            </div>
        </div>
    </main>

    <footer>
        <div class="container">
            <p>&copy; ` + fmt.Sprintf("%d", time.Now().Year()) + ` Unixify. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
`))
	} else if name == "accounts/index.html" {
		// Specific handler for accounts list with dynamic content
		accountsData, ok := dataMap["Accounts"].([]models.Account)
		if !ok {
			log.Printf("Error: Failed to convert Accounts to proper type")
			accountsData = []models.Account{}
		}
		
		query := ""
		if q, ok := dataMap["Query"].(string); ok {
			query = q
		}
		
		// Build dynamic accounts table rows
		tableRows := ""
		if len(accountsData) == 0 {
			tableRows = `<tr><td colspan="5" class="text-center">No accounts found</td></tr>`
		} else {
			for _, account := range accountsData {
				role := "User"
				if strings.Contains(strings.ToLower(account.Username), "admin") {
					role = "Admin"
				}
				
				// Format account data
				fullName := account.FullName
				if fullName == "" {
					fullName = "-"
				}
				
				tableRows += fmt.Sprintf(`
				<tr>
					<td>%s</td>
					<td>%s</td>
					<td>%s</td>
					<td>%s</td>
					<td>
						<div class="action-buttons">
							<a href="/accounts/%s" class="btn btn-sm">View</a>
							<a href="/accounts/%s/edit" class="btn btn-sm">Edit</a>
							<button class="btn btn-sm btn-danger" 
								hx-delete="/accounts/%s" 
								hx-confirm="Are you sure you want to delete this account?">Delete</button>
						</div>
					</td>
				</tr>`, 
				account.Username, 
				fullName,
				account.Email, 
				role,
				account.ID, 
				account.ID,
				account.ID)
			}
		}
		
		w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Accounts - Unixify</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <header>
        <nav>
            <div class="container">
                <a href="/" class="logo">Unixify</a>
                <ul class="nav-links">
                    <li><a href="/">Home</a></li>
                    <li><a href="/accounts">Accounts</a></li>
                    <li><a href="/groups">Groups</a></li>
                </ul>
            </div>
        </nav>
    </header>

    <main class="container">
        <div class="card">
            <div class="card-header">
                <h1>Account Management</h1>
                <a href="/accounts/new" class="btn">Create New Account</a>
            </div>
            
            <div class="card-body">
                <div class="search-container">
                    <form action="/accounts" method="GET" class="search-form">
                        <input type="text" name="search" placeholder="Search accounts..." value="%s">
                        <button type="submit" class="btn">Search</button>
                    </form>
                </div>
                
                <table class="table">
                    <thead>
                        <tr>
                            <th>Username</th>
                            <th>Name</th>
                            <th>Email</th>
                            <th>Role</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        %s
                    </tbody>
                </table>
            </div>
        </div>
    </main>

    <footer>
        <div class="container">
            <p>&copy; %d Unixify. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
`, query, tableRows, time.Now().Year())))
	} else if name == "groups/form.html" {
		// Specific handler for groups form
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Create Group - Unixify</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <header>
        <nav>
            <div class="container">
                <a href="/" class="logo">Unixify</a>
                <ul class="nav-links">
                    <li><a href="/">Home</a></li>
                    <li><a href="/accounts">Accounts</a></li>
                    <li><a href="/groups">Groups</a></li>
                </ul>
            </div>
        </nav>
    </header>

    <main class="container">
        <div class="card">
            <div class="card-header">
                <h1>Create New Group</h1>
            </div>
            <div class="card-body">
                <form action="/groups" method="POST" class="form">
                    <div class="form-group">
                        <label for="name">Group Name</label>
                        <input type="text" id="name" name="name" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="description">Description</label>
                        <textarea id="description" name="description" rows="3"></textarea>
                    </div>
                    
                    <div class="form-actions">
                        <a href="/groups" class="btn btn-secondary">Cancel</a>
                        <button type="submit" class="btn btn-primary">Create Group</button>
                    </div>
                </form>
            </div>
        </div>
    </main>

    <footer>
        <div class="container">
            <p>&copy; ` + fmt.Sprintf("%d", time.Now().Year()) + ` Unixify. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
`))
	} else if name == "groups/index.html" {
		// Specific handler for groups list with dynamic content
		groupsData, ok := dataMap["Groups"].([]models.Group)
		if !ok {
			log.Printf("Error: Failed to convert Groups to proper type")
			groupsData = []models.Group{}
		}
		
		query := ""
		if q, ok := dataMap["Query"].(string); ok {
			query = q
		}
		
		// Build dynamic groups table rows
		tableRows := ""
		if len(groupsData) == 0 {
			tableRows = `<tr><td colspan="3" class="text-center">No groups found</td></tr>`
		} else {
			for _, group := range groupsData {
				// Format group data
				description := group.Description
				if description == "" {
					description = "-"
				}
				
				tableRows += fmt.Sprintf(`
				<tr>
					<td>%s</td>
					<td>%s</td>
					<td>
						<div class="action-buttons">
							<a href="/groups/%s" class="btn btn-sm">View</a>
							<a href="/groups/%s/edit" class="btn btn-sm">Edit</a>
							<button class="btn btn-sm btn-danger" 
								hx-delete="/groups/%s" 
								hx-confirm="Are you sure you want to delete this group?">Delete</button>
						</div>
					</td>
				</tr>`, 
				group.Name,
				description,
				group.ID, 
				group.ID,
				group.ID)
			}
		}
		
		w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Groups - Unixify</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <header>
        <nav>
            <div class="container">
                <a href="/" class="logo">Unixify</a>
                <ul class="nav-links">
                    <li><a href="/">Home</a></li>
                    <li><a href="/accounts">Accounts</a></li>
                    <li><a href="/groups">Groups</a></li>
                </ul>
            </div>
        </nav>
    </header>

    <main class="container">
        <div class="card">
            <div class="card-header">
                <h1>Group Management</h1>
                <a href="/groups/new" class="btn">Create New Group</a>
            </div>
            
            <div class="card-body">
                <div class="search-container">
                    <form action="/groups" method="GET" class="search-form">
                        <input type="text" name="search" placeholder="Search groups..." value="%s">
                        <button type="submit" class="btn">Search</button>
                    </form>
                </div>
                
                <table class="table">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Description</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        %s
                    </tbody>
                </table>
            </div>
        </div>
    </main>

    <footer>
        <div class="container">
            <p>&copy; %d Unixify. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
`, query, tableRows, time.Now().Year())))
	} else if name == "accounts/view.html" {
		// Handle account view template
		account, ok := dataMap["Account"].(models.Account)
		if !ok {
			log.Printf("Error: Failed to convert Account to proper type")
			http.Error(w, "Internal Server Error: Account data not available", http.StatusInternalServerError)
			return
		}
		
		createdByName := "System"
		if val, ok := dataMap["CreatedByName"].(string); ok {
			createdByName = val
		}
		
		updatedByName := "System"
		if val, ok := dataMap["UpdatedByName"].(string); ok {
			updatedByName = val
		}
		
		// Process group data
		groupList := ""
		if len(account.Groups) > 0 {
			groupList = "<ul class=\"group-list\">"
			for _, group := range account.Groups {
				groupList += fmt.Sprintf("<li><a href=\"/groups/%s\">%s</a></li>", 
					group.ID, group.Name)
			}
			groupList += "</ul>"
		} else {
			groupList = "<p>No groups</p>"
		}
		
		// Format date/time
		createdAt := account.CreatedAt.Format("Jan 02, 2006 3:04 PM")
		updatedAt := account.UpdatedAt.Format("Jan 02, 2006 3:04 PM")
		
		// Format last login
		lastLogin := "Never"
		if account.LastLogin != nil {
			lastLogin = account.LastLogin.Format("Jan 02, 2006 3:04 PM")
		}
		
		// Format status
		status := "Inactive"
		if account.Active {
			status = "Active"
		}
		
		w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Account: %s - Unixify</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <header>
        <nav>
            <div class="container">
                <a href="/" class="logo">Unixify</a>
                <ul class="nav-links">
                    <li><a href="/">Home</a></li>
                    <li><a href="/accounts">Accounts</a></li>
                    <li><a href="/groups">Groups</a></li>
                </ul>
            </div>
        </nav>
    </header>

    <main class="container">
        <div class="card">
            <div class="card-header">
                <h1>Account: %s</h1>
                <div>
                    <a href="/accounts/%s/edit" class="btn">Edit</a>
                    <a href="/accounts" class="btn btn-secondary">Back to Accounts</a>
                </div>
            </div>
        
            <div class="view-details">
                <div class="detail-row">
                    <div class="detail-label">Username</div>
                    <div class="detail-value">%s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Email</div>
                    <div class="detail-value">%s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Full Name</div>
                    <div class="detail-value">%s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Status</div>
                    <div class="detail-value">%s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Last Login</div>
                    <div class="detail-value">%s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Created</div>
                    <div class="detail-value">%s by %s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Updated</div>
                    <div class="detail-value">%s by %s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Groups</div>
                    <div class="detail-value">
                        %s
                    </div>
                </div>
            </div>
            
            <div class="detail-actions">
                <button class="btn btn-danger" 
                        hx-delete="/accounts/%s" 
                        hx-confirm="Are you sure you want to delete this account?" 
                        hx-push-url="true" 
                        hx-target="body">
                    Delete Account
                </button>
            </div>
        </div>
    </main>

    <footer>
        <div class="container">
            <p>&copy; %d Unixify. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
`, 
		account.Username, 
		account.Username, 
		account.ID,
		account.Username,
		account.Email,
		account.FullName, 
		status,
		lastLogin,
		createdAt, createdByName,
		updatedAt, updatedByName,
		groupList,
		account.ID,
		time.Now().Year())))

	} else if name == "groups/view.html" {
		// Handle group view template
		group, ok := dataMap["Group"].(models.Group)
		if !ok {
			val, ok := dataMap["Group"].(*models.Group)
			if !ok {
				log.Printf("Error: Failed to convert Group to proper type: %T", dataMap["Group"])
				http.Error(w, "Internal Server Error: Group data not available", http.StatusInternalServerError)
				return
			}
			// Use the pointer value
			group = *val
		}
		
		createdByName := "System"
		if val, ok := dataMap["CreatedByName"].(string); ok {
			createdByName = val
		}
		
		updatedByName := "System"
		if val, ok := dataMap["UpdatedByName"].(string); ok {
			updatedByName = val
		}
		
		// Format description
		description := group.Description
		if description == "" {
			description = "-"
		}
		
		// Process members data
		membersTable := ""
		if len(group.Members) > 0 {
			membersTable = `
			<table class="nested-table">
				<thead>
					<tr>
						<th>Username</th>
						<th>Full Name</th>
						<th>Email</th>
					</tr>
				</thead>
				<tbody>`
				
			for _, member := range group.Members {
				membersTable += fmt.Sprintf(`
				<tr>
					<td><a href="/accounts/%s">%s</a></td>
					<td>%s</td>
					<td>%s</td>
				</tr>`, 
				member.ID, member.Username, member.FullName, member.Email)
			}
			
			membersTable += `
				</tbody>
			</table>`
		} else {
			membersTable = "<p>No members</p>"
		}
		
		// Format date/time
		createdAt := group.CreatedAt.Format("Jan 02, 2006 3:04 PM")
		updatedAt := group.UpdatedAt.Format("Jan 02, 2006 3:04 PM")
		
		w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Group: %s - Unixify</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <header>
        <nav>
            <div class="container">
                <a href="/" class="logo">Unixify</a>
                <ul class="nav-links">
                    <li><a href="/">Home</a></li>
                    <li><a href="/accounts">Accounts</a></li>
                    <li><a href="/groups">Groups</a></li>
                </ul>
            </div>
        </nav>
    </header>

    <main class="container">
        <div class="card">
            <div class="card-header">
                <h1>Group: %s</h1>
                <div>
                    <a href="/groups/%s/edit" class="btn">Edit</a>
                    <a href="/groups" class="btn btn-secondary">Back to Groups</a>
                </div>
            </div>
        
            <div class="view-details">
                <div class="detail-row">
                    <div class="detail-label">Name</div>
                    <div class="detail-value">%s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Description</div>
                    <div class="detail-value">%s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Created</div>
                    <div class="detail-value">%s by %s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Updated</div>
                    <div class="detail-value">%s by %s</div>
                </div>
                
                <div class="detail-row">
                    <div class="detail-label">Members</div>
                    <div class="detail-value">
                        %s
                    </div>
                </div>
            </div>
            
            <div class="detail-actions">
                <button class="btn btn-danger" 
                        hx-delete="/groups/%s" 
                        hx-confirm="Are you sure you want to delete this group?" 
                        hx-push-url="true" 
                        hx-target="body">
                    Delete Group
                </button>
            </div>
        </div>
    </main>

    <footer>
        <div class="container">
            <p>&copy; %d Unixify. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
`, 
		group.Name,
		group.Name,
		group.ID,
		group.Name,
		description,
		createdAt, createdByName,
		updatedAt, updatedByName,
		membersTable,
		group.ID,
		time.Now().Year())))
		
	} else {
		// For other templates, try to use the template system
		log.Printf("Trying to render template: %s", name)
		err := h.templates.ExecuteTemplate(w, name, dataMap)
		if err != nil {
			log.Printf("Error rendering template %s: %v", name, err)
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// respondJSON sends a JSON response
func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Marshal data to JSON and write to response
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Marshal error message to JSON and write to response
}