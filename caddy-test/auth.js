// auth.js - Basic authentication functions

/**
 * Check if the user is authenticated
 * @returns {boolean} True if authenticated, false otherwise
 */
function isAuthenticated() {
    return !!localStorage.getItem('auth_token');
}

/**
 * Get the authenticated user info
 * @returns {Object|null} User info or null if not authenticated
 */
function getCurrentUser() {
    const userInfo = localStorage.getItem('user_info');
    return userInfo ? JSON.parse(userInfo) : null;
}

/**
 * Check if the user has edit permissions
 * @returns {boolean} True if user can edit, false otherwise
 */
function hasEditPermission() {
    const user = getCurrentUser();
    // In a real app, this would check specific permissions
    // For demo, we'll say all authenticated users can edit
    return !!user;
}

/**
 * Logout the current user
 */
function logout() {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_info');
    window.location.href = '/login.html';
}

/**
 * Initialize authentication UI
 */
function initAuthUI() {
    console.log('Initializing auth UI');
    
    // Add logout handler to any logout buttons
    const logoutBtns = document.querySelectorAll('[id^="logout-btn"]');
    logoutBtns.forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            logout();
        });
    });
    
    // Show guest mode alert if applicable
    const guestModeAlert = document.getElementById('guest-mode-alert');
    if (guestModeAlert && !isAuthenticated()) {
        guestModeAlert.style.display = 'block';
    }
    
    // Show user info if logged in
    const user = getCurrentUser();
    if (user) {
        // Add user info to navbar
        const navbarRight = document.querySelector('.ms-auto');
        if (navbarRight) {
            // Add username display
            const userInfo = document.createElement('span');
            userInfo.classList.add('navbar-text', 'me-3');
            userInfo.innerHTML = `<i class="bi bi-person-circle"></i> ${user.username}`;
            navbarRight.appendChild(userInfo);
            
            // Add logout button
            const logoutBtn = document.createElement('button');
            logoutBtn.id = 'logout-btn';
            logoutBtn.classList.add('btn', 'btn-outline-light', 'btn-sm');
            logoutBtn.innerHTML = '<i class="bi bi-box-arrow-right"></i> Logout';
            logoutBtn.addEventListener('click', logout);
            navbarRight.appendChild(logoutBtn);
        }
    }
    
    // Toggle edit controls based on authentication
    updateEditControls();
}

/**
 * Update UI elements based on authentication status
 */
function updateEditControls() {
    // If user is not authenticated, disable edit controls
    const canEdit = hasEditPermission();
    
    // Edit buttons usually have add, edit, delete, or modify in their class or ID
    const editButtons = document.querySelectorAll('.btn-edit, .btn-add, .btn-delete, [id*="edit"], [id*="add"], [id*="delete"]');
    
    editButtons.forEach(button => {
        if (canEdit) {
            button.classList.remove('disabled');
            button.removeAttribute('disabled');
            // Add tooltip indicating edit is enabled
            button.setAttribute('title', 'Click to edit');
        } else {
            button.classList.add('disabled');
            button.setAttribute('disabled', 'disabled');
            // Add tooltip explaining why editing is disabled
            button.setAttribute('title', 'Login required to edit');
        }
    });
    
    // For form inputs, textareas, selects, etc.
    const formControls = document.querySelectorAll('input:not([type="submit"]), textarea, select');
    
    formControls.forEach(control => {
        if (canEdit) {
            control.removeAttribute('disabled');
            control.removeAttribute('readonly');
        } else {
            // Use readonly instead of disabled for better appearance
            control.setAttribute('readonly', 'readonly');
        }
    });
    
    // Add visual indicator
    const editModeIndicator = document.getElementById('edit-mode-indicator');
    if (editModeIndicator) {
        editModeIndicator.textContent = canEdit ? 'Edit Mode' : 'Read-Only Mode';
        editModeIndicator.className = canEdit ? 'badge bg-success' : 'badge bg-secondary';
    }
}

// Check authentication on page load
document.addEventListener('DOMContentLoaded', function() {
    console.log('Auth.js loaded');
    
    // Initialize auth UI
    initAuthUI();
    
    // Check if page requires authentication
    const requiresAuth = document.body.getAttribute('data-requires-auth') === 'true';
    console.log('Page requires auth:', requiresAuth);
    
    if (requiresAuth && !isAuthenticated()) {
        console.log('Redirecting to login page');
        window.location.href = '/login.html';
    }
});