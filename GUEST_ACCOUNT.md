# Guest Account Implementation

This document explains the implementation of the guest account feature in the Unixify application.

## Overview

The application now automatically logs in users as a "guest" account when they access the system without authentication. This provides a unified experience where:

1. All users are technically "authenticated" but with different permission levels
2. Guest users are clearly identified in the UI
3. A consistent navigation experience is maintained for all users

## How It Works

### 1. Automatic Guest Login

When a user accesses the application without a valid authentication token, the system:

1. Automatically creates a guest token in localStorage
2. Sets up a guest user profile with username "guest"
3. Displays the guest account in the navigation bar
4. Applies read-only permissions to all UI elements

### 2. Guest User Interface

The guest user experience includes:

- A yellow "Guest Account (Read-Only)" indicator in the navbar
- A special yellow dashed avatar with "G" (for Guest)
- A user dropdown menu showing the guest username and role
- Disabled edit buttons throughout the interface
- A "Register Now" banner encouraging account creation

### 3. Authentication Flow

The authentication system manages three states:

1. **Not Authenticated**: No token of any kind (redirects to login)
2. **Guest User**: Has a guest token (read-only access)
3. **Authenticated User**: Has a regular auth token (permissions based on role)

### 4. Technical Implementation

The guest account is implemented through:

- A `guest_token` in localStorage that identifies guest sessions
- A `isGuestUser()` function that differentiates between guest and regular users
- Special CSS styling for guest UI elements
- Modified permission checking to recognize and handle guest accounts
- Updated templates to display guest-specific UI elements

## Testing the Guest Account

### Using the Simplified Server

1. Run the simplified server:
```bash
cd /home/pfrederi/code/github.com/home/unixify/feature/UNO-861-acc-man
PORT=8083 go run cmd/simplified/main.go
```

2. Access the application at http://localhost:8083
3. You'll be automatically logged in as the guest user
4. To test regular authentication, use the login page with:
   - Username: admin
   - Password: admin

### Using the Login API

The application provides a mock login API that accepts:

- Guest login: username "guest" with any password
- Admin login: username "admin" with password "admin"

Example:
```javascript
// Guest login
fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: 'guest', password: 'anything' })
})
```

## Benefits

The guest account approach provides several benefits:

1. **Unified Code Path**: The code can treat all users as authenticated, simplifying logic
2. **Clear Visual Indicators**: Users always know their current access level
3. **Smoother UX**: No jarring transitions between authenticated and non-authenticated states
4. **Easy Registration Path**: Clear path for users to upgrade from guest to registered user
5. **Permission Management**: Centralized permission system that works for all user types