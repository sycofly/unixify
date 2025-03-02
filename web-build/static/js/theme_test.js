// This is a test file to see if changes are being reflected in the container
console.log("Theme test file loaded - UPDATED IN CONTAINERFILE BUILD");

// Create a theme button in JavaScript that will appear regardless of template issues
document.addEventListener('DOMContentLoaded', function() {
    // Update footer text as a test
    const footers = document.querySelectorAll('footer p.text-muted');
    footers.forEach(footer => {
        footer.innerHTML = 'Unixify &copy; 2025 - A UNIX Team Management System';
        console.log('Footer updated via JavaScript');
    });
    
    // Create the button element
    const themeBtn = document.createElement('button');
    themeBtn.id = 'js-theme-toggle';
    themeBtn.className = 'btn btn-danger';
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
        
        // Save theme preference
        localStorage.setItem('theme', newTheme);
        
        // Update button appearance
        if (newTheme === 'dark') {
            themeBtn.innerHTML = '<i class="bi bi-sun-fill"></i> Light';
            themeBtn.classList.remove('btn-danger');
            themeBtn.classList.add('btn-light');
        } else {
            themeBtn.innerHTML = '<i class="bi bi-moon-stars-fill"></i> Theme';
            themeBtn.classList.remove('btn-light');
            themeBtn.classList.add('btn-danger');
        }
        
        console.log('Theme toggled to:', newTheme);
    });
    
    // Add the button to the body
    document.body.appendChild(themeBtn);
    
    // Apply saved theme if any
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
        document.documentElement.setAttribute('data-theme', savedTheme);
        
        // Update button to match theme
        if (savedTheme === 'dark') {
            themeBtn.innerHTML = '<i class="bi bi-sun-fill"></i> Light';
            themeBtn.classList.remove('btn-danger');
            themeBtn.classList.add('btn-light');
        }
    }
});