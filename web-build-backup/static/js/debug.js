// Debug script to check if static files are loading correctly
window.addEventListener('DOMContentLoaded', () => {
  console.log('Debug script loaded successfully');
  
  // Create debugging element
  const debugDiv = document.createElement('div');
  debugDiv.style.position = 'fixed';
  debugDiv.style.top = '80px';
  debugDiv.style.right = '15px';
  debugDiv.style.zIndex = '9999';
  debugDiv.style.background = 'red';
  debugDiv.style.color = 'white';
  debugDiv.style.padding = '10px';
  debugDiv.style.borderRadius = '5px';
  debugDiv.textContent = 'Debug Message: Theme scripts should be loaded';
  document.body.appendChild(debugDiv);
  
  // Check if theme files are loaded
  const scriptElements = document.querySelectorAll('script');
  const styles = document.querySelectorAll('link[rel="stylesheet"]');
  
  // Get theme button status
  const themeButton = document.getElementById('theme-toggle');
  const jsThemeButton = document.getElementById('js-theme-toggle');
  
  // Update debug information
  setTimeout(() => {
    debugDiv.innerHTML = `
      <strong>Debug Info:</strong><br>
      theme.js loaded: ${Array.from(scriptElements).some(s => s.src.includes('theme.js'))}<br>
      theme_test.js loaded: ${Array.from(scriptElements).some(s => s.src.includes('theme_test.js'))}<br>
      main.css loaded: ${Array.from(styles).some(s => s.href.includes('main.css'))}<br>
      Theme button exists: ${themeButton !== null}<br>
      JS Theme button exists: ${jsThemeButton !== null}<br>
      Current theme: ${document.documentElement.getAttribute('data-theme') || 'not set'}<br>
      Footer text: ${document.querySelector('footer p.text-muted')?.textContent || 'not found'}
    `;
  }, 2000);
});