Please create a Go application that serves as a registry for UNIX account UIDs and Group GIDs. The application should:

1. Use PostgreSQL as the database backend
2. Provide a Web interface for managing UIDs and GIDs using Gin web framework
2. Web interface has four sections, People, System, Database and Service
2. Sections People, System, Database and Service will have an account and group operations
3. Support the following operations:
   - Add a new user with a UID
   - Add a new group with a GID
   - Assign users to groups
   - Remove users from groups
   - List all users
   - List all groups
   - List users in a specific group
   - Look up user by UID
   - Look up group by GID
   - soft delete users
   - Soft delete groups
   - search database by either UID GID Username Groupname
   - have a auditing table for all events
   - System accounts can belong to System Groups
   - Users accounts can belong to People Groups
   - Users accounts can belong to Database Groups
   - Database accounts belong to Database Groups
   - Database admins belong to People Groups
   - Service accounts belong to service Groups

4. Include proper validation to ensure:
   - UIDs and GIDs are unique
   - UIDs and GIDs are within valid UNIX ranges
   - Users can't be assigned to non-existent groups

5. Implement proper error handling, logging, and data persistence
6. Follow Go best practices for project structure and code organization
7. Include Podman files for easy deployment
8. Include database migration scripts
9. Provide documentation on how to build, deploy, and use the application

Please include example usage of the CLI commands and schema design for the PostgreSQL database.
