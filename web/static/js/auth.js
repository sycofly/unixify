// auth.js - Authentication-related JavaScript functions

/**
 * Check if the user is authenticated
 * @returns {boolean} True if authenticated, false otherwise
 */
function isAuthenticated() {
    // Check if user is logged in with a token
    const token = localStorage.getItem('auth_token');
    
    // If no token, but we have a guest token, user is authenticated as guest
    if (!token && localStorage.getItem('guest_token')) {
        return true;
    }
    
    return !!token;
}

/**
 * Check if the user is in guest mode
 * @returns {boolean} True if guest user, false otherwise
 */
function isGuestUser() {
    return !localStorage.getItem('auth_token') && !!localStorage.getItem('guest_token');
}

/**
 * Get the authenticated user info
 * @returns {Object|null} User info or null if not authenticated
 */
function getCurrentUser() {
    // First try to get regular user info
    const userInfo = localStorage.getItem('user_info');
    if (userInfo) {
        return JSON.parse(userInfo);
    }
    
    // If guest token exists, return guest user info
    if (localStorage.getItem('guest_token')) {
        return {
            id: 0,
            username: "guest",
            email: "guest@example.com",
            firstName: "Guest",
            lastName: "User",
            role: "guest",
            isGuest: true
        };
    }
    
    return null;
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
    // For regular users
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_info');
    
    // For guest users
    localStorage.removeItem('guest_token');
    
    // Create a new guest session before redirecting
    localStorage.setItem('guest_token', 'guest_session_' + Date.now());
    
    // Redirect to home page, not login (since we're using guest access)
    window.location.href = '/';
}

/**
 * Check if the user has edit permissions
 * @returns {boolean} True if user has edit permissions, false otherwise
 */
function hasEditPermission() {
    // If not authenticated or is guest user, always return false
    if (!isAuthenticated() || isGuestUser()) return false;
    
    // Get user info
    const user = getCurrentUser();
    if (!user) return false;
    
    // Guest users never have edit permissions
    if (user.role === 'guest') return false;
    
    // Check if user role has edit permissions
    // In this implementation, only 'admin' and 'editor' roles can edit
    return ['admin', 'editor'].includes(user.role);
}

/**
 * Initialize authentication UI components
 */
function initAuthUI() {
    // Ensure guest token is set if no regular token exists
    if (!localStorage.getItem('auth_token') && !localStorage.getItem('guest_token')) {
        localStorage.setItem('guest_token', 'guest_session_' + Date.now());
    }
    
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
    const userRoleElements = document.querySelectorAll('[data-user-role]');
    const userAvatarElements = document.querySelectorAll('[data-user-avatar]');
    
    const isLoggedIn = isAuthenticated();
    const isGuest = isGuestUser();
    const user = getCurrentUser();
    const canEdit = hasEditPermission();
    
    // Set edit permission attribute on body for CSS selectors
    document.body.setAttribute('data-edit-permission', canEdit.toString());
    document.body.setAttribute('data-guest-mode', isGuest.toString());
    
    // Show elements that require authentication - hide for guests
    authElements.forEach(el => {
        // Special handling for user profile dropdown
        if (el.classList.contains('auth-user-container') || el.closest('.auth-user-container')) {
            // Only show profile dropdown for regular users, not for guests
            el.style.display = isLoggedIn && !isGuest ? '' : 'none';
        } else {
            // For other auth-required elements, show for both guests and regular users
            el.style.display = isLoggedIn ? '' : 'none';
        }
    });
    
    // Show elements that should be hidden when authenticated
    noAuthElements.forEach(el => {
        el.style.display = isLoggedIn && !isGuest ? 'none' : '';
    });
    
    // Handle register-now container - show for guests but hide for regular users
    const registerContainer = document.getElementById('register-now-container');
    if (registerContainer) {
        registerContainer.style.display = isGuest ? '' : 'none';
    }
    
    // Handle guest mode indicator - only show for guest users
    const guestIndicator = document.getElementById('guest-mode-indicator');
    if (guestIndicator) {
        guestIndicator.style.display = isGuest ? '' : 'none';
    }
    
    // Update user info elements
    if (user) {
        // Update username displays
        userNameElements.forEach(el => {
            el.textContent = user.username;
        });
        
        // Update user role displays
        userRoleElements.forEach(el => {
            el.textContent = user.role.charAt(0).toUpperCase() + user.role.slice(1);
        });
        
        // Update user avatar elements with initials
        userAvatarElements.forEach(el => {
            // Create initials from username (first letter)
            let initials = user.username.charAt(0).toUpperCase();
            if (user.firstName && user.lastName) {
                initials = user.firstName.charAt(0).toUpperCase() + user.lastName.charAt(0).toUpperCase();
            }
            el.textContent = initials;
            
            // Add special styling for guest avatar
            if (isGuest) {
                el.classList.add('guest-avatar');
            } else {
                el.classList.remove('guest-avatar');
            }
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
    
    // Handle elements that require edit permission
    const editElements = document.querySelectorAll('.edit-only');
    editElements.forEach(el => {
        el.style.display = canEdit ? '' : 'none';
    });
}

// Check authentication on page load
document.addEventListener('DOMContentLoaded', function() {
    // Ensure guest token is set if no regular token exists for automatic guest mode
    if (!localStorage.getItem('auth_token') && !localStorage.getItem('guest_token')) {
        localStorage.setItem('guest_token', 'guest_session_' + Date.now());
        console.log("Auto-creating guest session");
    }
    
    // Initialize auth UI
    initAuthUI();
    
    // Redirect to login if strict authentication is required but user is not logged in
    const requiresStrictAuth = document.body.hasAttribute('data-requires-strict-auth');
    
    // Only strictly protected pages redirect to login
    if (requiresStrictAuth && !isAuthenticated()) {
        window.location.href = '/login';
    } else {
        console.log("Guest mode active or user authenticated");
    }
});