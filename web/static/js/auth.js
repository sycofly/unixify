// auth.js - Authentication-related JavaScript functions

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
 * Get the authentication token
 * @returns {string|null} JWT token or null if not authenticated
 */
function getAuthToken() {
    return localStorage.getItem('auth_token');
}

/**
 * Add authentication header to fetch options
 * @param {Object} options - Fetch options
 * @returns {Object} Options with authorization header
 */
function addAuthHeader(options = {}) {
    const token = getAuthToken();
    if (!token) return options;
    
    if (!options.headers) {
        options.headers = {};
    }
    
    options.headers['Authorization'] = `Bearer ${token}`;
    return options;
}

/**
 * Authenticated fetch function
 * @param {string} url - URL to fetch
 * @param {Object} options - Fetch options
 * @returns {Promise} Fetch promise
 */
async function authFetch(url, options = {}) {
    const authOptions = addAuthHeader(options);
    const response = await fetch(url, authOptions);
    
    // If unauthorized, redirect to login page
    if (response.status === 401) {
        logout();
        window.location.href = '/login';
        throw new Error('Authentication required');
    }
    
    return response;
}

/**
 * Logout the current user
 */
function logout() {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_info');
    window.location.href = '/login';
}

/**
 * Initialize authentication UI components
 */
function initAuthUI() {
    // Add logout handler to logout button if it exists
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', function(e) {
            e.preventDefault();
            logout();
        });
    }
    
    // Show/hide elements based on authentication status
    const authElements = document.querySelectorAll('[data-auth-required]');
    const noAuthElements = document.querySelectorAll('[data-auth-forbidden]');
    const userNameElements = document.querySelectorAll('[data-user-name]');
    
    const isLoggedIn = isAuthenticated();
    const user = getCurrentUser();
    
    // Show elements that require authentication
    authElements.forEach(el => {
        el.style.display = isLoggedIn ? '' : 'none';
    });
    
    // Show elements that should be hidden when authenticated
    noAuthElements.forEach(el => {
        el.style.display = isLoggedIn ? 'none' : '';
    });
    
    // Update elements with user name
    if (user) {
        userNameElements.forEach(el => {
            el.textContent = user.username;
        });
    }
    
    // Check for role-based elements
    if (user) {
        const roleElements = document.querySelectorAll('[data-role]');
        roleElements.forEach(el => {
            const requiredRoles = el.dataset.role.split(',');
            const hasRequiredRole = requiredRoles.includes(user.role);
            el.style.display = hasRequiredRole ? '' : 'none';
        });
    }
}

// Check authentication on page load
document.addEventListener('DOMContentLoaded', function() {
    // Initialize auth UI
    initAuthUI();
    
    // Redirect to login if authentication is required but user is not logged in
    const requiresAuth = document.body.hasAttribute('data-requires-auth');
    if (requiresAuth && !isAuthenticated()) {
        window.location.href = '/login';
    }
});