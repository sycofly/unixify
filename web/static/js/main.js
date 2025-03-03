// Main JavaScript file for all pages

// Function to display error messages
function showError(message) {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'alert alert-danger alert-dismissible fade show';
    errorDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;
    
    // Insert at the top of the container
    const container = document.querySelector('.container');
    if (container) {
        container.insertBefore(errorDiv, container.firstChild);
        
        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            errorDiv.remove();
        }, 5000);
    } else {
        console.error('Container element not found for error message:', message);
    }
}

// Function to display success messages
function showSuccess(message) {
    const successDiv = document.createElement('div');
    successDiv.className = 'alert alert-success alert-dismissible fade show';
    successDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;
    
    // Insert at the top of the container
    const container = document.querySelector('.container');
    if (container) {
        container.insertBefore(successDiv, container.firstChild);
        
        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            successDiv.remove();
        }, 5000);
    } else {
        console.error('Container element not found for success message:', message);
    }
}

// Helper function for making API requests
async function apiRequest(url, method = 'GET', data = null) {
    console.log(`Making ${method} request to ${url}`, data);
    
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json'
        }
    };
    
    if (data) {
        options.body = JSON.stringify(data);
        console.log('Request body:', options.body);
    }
    
    try {
        console.log('Fetch options:', options);
        const response = await fetch(url, options);
        console.log('Response status:', response.status);
        
        const contentType = response.headers.get('content-type');
        console.log('Response content type:', contentType);
        
        let result;
        
        if (contentType && contentType.includes('application/json')) {
            result = await response.json();
            console.log('Response JSON:', result);
        } else {
            const text = await response.text();
            console.log('Response text:', text);
            result = { message: text };
        }
        
        if (!response.ok) {
            console.error('Response not OK:', response.status, result);
            throw new Error(result.error || `Error: ${response.status} ${response.statusText}`);
        }
        
        return result;
    } catch (error) {
        console.error('API request error:', error);
        showError(error.message || 'Unknown error occurred');
        throw error;
    }
}

// Format date string using the International Standard format (ISO 8601)
function formatDate(dateString) {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    // Format as YYYY-MM-DD HH:MM:SS
    return date.getFullYear() + '-' + 
           String(date.getMonth() + 1).padStart(2, '0') + '-' + 
           String(date.getDate()).padStart(2, '0') + ' ' + 
           String(date.getHours()).padStart(2, '0') + ':' + 
           String(date.getMinutes()).padStart(2, '0') + ':' + 
           String(date.getSeconds()).padStart(2, '0');
}

// No longer needed - moved to auth.js to ensure it runs before any auth checks