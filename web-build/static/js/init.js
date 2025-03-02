// This script initializes the basic functionality for static HTML pages
// This is needed because the Caddy server serves HTML directly without template processing

document.addEventListener('DOMContentLoaded', function() {
    console.log('Initialization script loaded for static pages');
    
    // Set page title for static pages
    document.title = document.title || 'Unixify - UNIX Team Management System';
    
    // Update auth elements if needed
    const requiresAuth = document.body.getAttribute('data-requires-auth');
    if (requiresAuth === 'true') {
        console.log('This page requires authentication in dynamic mode');
    }
    
    // Initialize theme
    initializeTheme();
    
    // Debug info
    console.log('Static page initialization complete');
});

// Theme initialization function
function initializeTheme() {
    // Initialize manually for static pages
    const savedTheme = localStorage.getItem('theme');
    
    // Set initial theme
    if (savedTheme) {
        document.documentElement.setAttribute('data-theme', savedTheme);
        console.log('Loaded saved theme:', savedTheme);
    } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.setAttribute('data-theme', 'dark');
        localStorage.setItem('theme', 'dark');
        console.log('Using system preference: dark mode');
    } else {
        document.documentElement.setAttribute('data-theme', 'light');
        localStorage.setItem('theme', 'light');
        console.log('Using system preference or default: light mode');
    }
    
    // Find theme toggle buttons
    const themeToggles = document.querySelectorAll('[id^="theme-toggle"], [id^="js-theme-toggle"]');
    
    // Add click handlers
    themeToggles.forEach(button => {
        button.addEventListener('click', function() {
            const currentTheme = document.documentElement.getAttribute('data-theme');
            const newTheme = currentTheme === 'light' ? 'dark' : 'light';
            
            // Set the theme
            document.documentElement.setAttribute('data-theme', newTheme);
            localStorage.setItem('theme', newTheme);
            
            // Update button UI
            updateThemeButton(button, newTheme);
            
            console.log('Theme toggled to:', newTheme);
        });
        
        // Set initial button state based on current theme
        updateThemeButton(button, document.documentElement.getAttribute('data-theme'));
    });
    
    // If no theme buttons found, create one
    if (themeToggles.length === 0) {
        createThemeButton();
    }
    
    // Add listener for system theme changes
    if (window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
            if (!localStorage.getItem('theme')) {
                const newTheme = e.matches ? 'dark' : 'light';
                document.documentElement.setAttribute('data-theme', newTheme);
                localStorage.setItem('theme', newTheme);
                console.log('System theme changed to:', newTheme);
                
                // Update all theme buttons
                document.querySelectorAll('[id^="theme-toggle"], [id^="js-theme-toggle"]').forEach(button => {
                    updateThemeButton(button, newTheme);
                });
            }
        });
    }
}

// Update theme button appearance
function updateThemeButton(button, theme) {
    if (theme === 'dark') {
        button.setAttribute('title', 'Switch to Light Mode');
        button.setAttribute('aria-label', 'Switch to Light Mode');
        button.classList.remove('btn-warning');
        button.classList.add('btn-light');
        button.innerHTML = '<i class="bi bi-sun-fill"></i> Light';
    } else {
        button.setAttribute('title', 'Switch to Dark Mode');
        button.setAttribute('aria-label', 'Switch to Dark Mode');
        button.classList.remove('btn-light');
        button.classList.add('btn-warning');
        button.innerHTML = '<i class="bi bi-moon-stars-fill"></i> Theme';
    }
}

// Create a theme button if none exists
function createThemeButton() {
    const themeBtn = document.createElement('button');
    themeBtn.id = 'js-theme-toggle';
    themeBtn.className = 'btn btn-warning';
    themeBtn.innerHTML = '<i class="bi bi-moon-stars-fill"></i> Theme';
    themeBtn.title = 'Toggle dark/light mode';
    
    // Style the button for top-right positioning
    themeBtn.style.position = 'fixed';
    themeBtn.style.top = '15px';
    themeBtn.style.right = '15px';
    themeBtn.style.zIndex = '9999';
    
    // Add click event to toggle theme
    themeBtn.addEventListener('click', function() {
        const currentTheme = document.documentElement.getAttribute('data-theme') || 'light';
        const newTheme = currentTheme === 'light' ? 'dark' : 'light';
        
        // Set the theme
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
        
        // Update button appearance
        updateThemeButton(themeBtn, newTheme);
        
        console.log('Theme toggled to:', newTheme);
    });
    
    // Add the button to the body
    document.body.appendChild(themeBtn);
    console.log('Created dynamic theme button');
    
    // Set initial state
    updateThemeButton(themeBtn, document.documentElement.getAttribute('data-theme') || 'light');
}

// Ensure footer text is correct
document.addEventListener('DOMContentLoaded', function() {
    const footers = document.querySelectorAll('footer p.text-muted');
    footers.forEach(footer => {
        if (!footer.textContent.includes('Team Management')) {
            footer.innerHTML = 'Unixify &copy; 2025 - A UNIX Team Management System';
            console.log('Footer text updated');
        }
    });
});