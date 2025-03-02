// Theme switcher functionality

/**
 * Sets the theme for the application
 * @param {string} theme - The theme to set (light or dark)
 */
function setTheme(theme) {
    // Set the theme attribute on the document
    document.documentElement.setAttribute('data-theme', theme);
    
    // Save the theme preference to localStorage
    localStorage.setItem('theme', theme);
    
    // Update button states
    updateThemeToggle(theme);
}

/**
 * Get the current theme preference
 * @returns {string} The current theme (light or dark)
 */
function getTheme() {
    // Check if a theme is saved in localStorage
    const savedTheme = localStorage.getItem('theme');
    
    if (savedTheme) {
        return savedTheme;
    }
    
    // Check if user has system preference
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'dark';
    }
    
    // Default to light theme
    return 'light';
}

/**
 * Update the theme toggle buttons based on the current theme
 * @param {string} theme - The current theme (light or dark)
 */
function updateThemeToggle(theme) {
    // Find the theme toggle button
    const toggleBtn = document.getElementById('theme-toggle');
    
    if (toggleBtn) {
        if (theme === 'dark') {
            toggleBtn.setAttribute('title', 'Switch to Light Mode');
            toggleBtn.setAttribute('aria-label', 'Switch to Light Mode');
            toggleBtn.classList.remove('btn-warning');
            toggleBtn.classList.add('btn-light');
            toggleBtn.innerHTML = '<i class="bi bi-sun-fill"></i> Light';
        } else {
            toggleBtn.setAttribute('title', 'Switch to Dark Mode');
            toggleBtn.setAttribute('aria-label', 'Switch to Dark Mode');
            toggleBtn.classList.remove('btn-light');
            toggleBtn.classList.add('btn-warning');
            toggleBtn.innerHTML = '<i class="bi bi-moon-stars-fill"></i> Theme';
        }
    }
    
    // Log theme state for debugging
    console.log("Theme toggled. Current theme:", theme);
    console.log("Theme button found:", toggleBtn !== null);
}

/**
 * Toggle between light and dark themes
 */
function toggleTheme() {
    const currentTheme = getTheme();
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    setTheme(newTheme);
}

/**
 * Initialize the theme system
 */
function initTheme() {
    // Apply the saved theme or default
    const theme = getTheme();
    setTheme(theme);
    
    // Add event listener to all theme toggle buttons
    const themeToggles = document.querySelectorAll('[id^="theme-toggle"]');
    themeToggles.forEach(toggleBtn => {
        toggleBtn.addEventListener('click', toggleTheme);
    });
    
    // Listen for system theme changes
    if (window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
            const newTheme = e.matches ? 'dark' : 'light';
            setTheme(newTheme);
        });
    }
}

// Initialize theme when DOM is loaded
document.addEventListener('DOMContentLoaded', initTheme);