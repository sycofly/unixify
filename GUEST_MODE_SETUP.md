# Guest Mode and Theme Updates Setup

This document explains how to run and test the new authentication improvements, guest mode functionality, and theme enhancements in the Unixify application.

## Running the Application

There are two ways to run the application:

### 1. Using the Simplified Server (Recommended for Testing UI Only)

The simplified server provides a quick way to test the UI changes without dealing with database connections or authentication complexities.

```bash
# Run on port 8083
cd /home/pfrederi/code/github.com/home/unixify/feature/UNO-861-acc-man
PORT=8083 go run cmd/simplified/main.go
```

Then access the application at http://localhost:8083

### 2. Using the Full Application

If you want to run the full application with all the backend functionality:

1. First, apply the model fixes to ensure the application compiles:
   - Added account and group type definitions
   - Updated models with missing fields
   - Added compatibility helpers

2. Then run the application with:
```bash
cd /home/pfrederi/code/github.com/home/unixify/feature/UNO-861-acc-man
SERVER_PORT=8083 go run cmd/unixify/main.go
```

## Features Implemented

1. **Guest Mode**
   - Anyone can now access the UI without being authenticated
   - Guest users see a "Read-Only Mode" indicator in the navbar
   - Edit buttons and form actions are disabled for guests
   - A "Register Now" banner encourages guests to create an account

2. **Improved Authentication**
   - Enhanced user avatar with username and role display
   - Added dropdown menu for user settings and logout
   - Implemented proper permission-based UI elements
   - Created a multi-step registration workflow

3. **Theme Enhancements**
   - Theme toggle is now integrated into the navbar
   - Improved dark mode colors for better contrast and readability
   - Added new CSS variables for more comprehensive theming
   - Fixed text visibility issues in dark mode

4. **Other UI Improvements**
   - Updated footer text to "UNIX Team Management System"
   - Improved footer styling
   - Better mobile responsiveness
   - Consistent styling across all pages

## Authentication Flow

For testing, a simplified auth flow is provided:

1. **Guest Mode (Default):**
   - No authentication required
   - Read-only access to all pages
   - Edit buttons disabled
   - "Register Now" banner visible

2. **Registration Flow:**
   - Multi-step process: Account Details → Email Verification → Admin Approval
   - Demo simulation of email verification
   - Demo simulation of admin approval

3. **Login Access:**
   - Demo credentials available (see login page)
   - After login, full edit permissions are granted

## Implementation Details

The main changes were made to:

1. **CSS and Theme System**
   - Added new CSS variables for complete theming
   - Improved styling for authentication elements
   - Added guest mode indicators

2. **JavaScript Authentication**
   - Enhanced auth.js to support guest mode
   - Implemented proper permission checks
   - Added new handler for guest vs. authenticated states

3. **HTML Templates**
   - Updated navigation with auth UI elements
   - Added guest mode indicators
   - Integrated theme toggle in navbar
   - Added permission-based display rules

4. **Backend Support**
   - Created a GuestMiddleware to allow anonymous access
   - Updated route handling to support guest mode
   - Maintained API protection for write operations

## Known Limitations

- The demo version doesn't have fully functional authentication with backend integration
- Local storage is used for auth state (lost on browser refresh)
- Registration form submissions are simulated in JS

## Next Steps

After testing the UI improvements, these changes will be properly integrated with the backend authentication system to provide a complete end-to-end solution.