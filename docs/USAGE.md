# Unixify Usage Guide

This document provides instructions on how to use the Unixify UNIX account management system.

## Table of Contents

1. [Introduction](#introduction)
2. [Account Types](#account-types)
3. [Web Interface](#web-interface)
4. [API Endpoints](#api-endpoints)
5. [UID/GID Ranges](#uidgid-ranges)
6. [Database Schema](#database-schema)

## Introduction

Unixify is a comprehensive system for managing UNIX account UIDs and Group GIDs. It provides a web interface and API for creating, updating, and managing accounts and groups across different categories.

## Account Types

Unixify supports four types of accounts and groups:

1. **People**: Regular user accounts and groups for human users
2. **System**: System accounts and groups for system processes
3. **Database**: Database accounts and groups for database access
4. **Service**: Service accounts and groups for application services

## Web Interface

The web interface is divided into four sections corresponding to the account types:

### Navigation

- **People**: Manage regular user accounts and groups
- **System**: Manage system accounts and groups
- **Database**: Manage database accounts and groups
- **Service**: Manage service accounts and groups

### Managing Accounts

1. Go to the appropriate section (People, System, Database, or Service)
2. Click the "Accounts" tab
3. Use the "New Account" button to create a new account
4. For existing accounts, use the action buttons:
   - **Edit**: Modify account details
   - **Delete**: Soft delete the account
   - **Groups**: Manage the account's group memberships

### Managing Groups

1. Go to the appropriate section
2. Click the "Groups" tab
3. Use the "New Group" button to create a new group
4. For existing groups, use the action buttons:
   - **Edit**: Modify group details
   - **Delete**: Soft delete the group
   - **Members**: Manage the group's members

### Search

Use the search bar in the navigation to find accounts or groups by:
- UID
- Username
- GID
- Groupname

## API Endpoints

Unixify provides a RESTful API for programmatic access:

### Account Endpoints

- `GET /api/accounts`: Get all accounts (optional query param: `type`)
- `GET /api/accounts/:id`: Get account by ID
- `POST /api/accounts`: Create a new account
- `PUT /api/accounts/:id`: Update an account
- `DELETE /api/accounts/:id`: Delete an account
- `GET /api/accounts/uid/:uid`: Get account by UID
- `GET /api/accounts/username/:username`: Get account by username
- `GET /api/accounts/:id/groups`: Get groups for an account

### Group Endpoints

- `GET /api/groups`: Get all groups (optional query param: `type`)
- `GET /api/groups/:id`: Get group by ID
- `POST /api/groups`: Create a new group
- `PUT /api/groups/:id`: Update a group
- `DELETE /api/groups/:id`: Delete a group
- `GET /api/groups/gid/:gid`: Get group by GID
- `GET /api/groups/groupname/:groupname`: Get group by groupname
- `GET /api/groups/:id/accounts`: Get accounts in a group

### Membership Endpoints

- `POST /api/memberships`: Assign account to group
  ```json
  {
    "account_id": 1,
    "group_id": 2
  }
  ```
- `DELETE /api/memberships`: Remove account from group
  ```json
  {
    "account_id": 1,
    "group_id": 2
  }
  ```

### Search Endpoints

- `GET /api/search/accounts?q=query`: Search accounts
- `GET /api/search/groups?q=query`: Search groups

### Audit Endpoints

- `GET /api/audit`: Get audit entries
- `GET /api/audit/:id`: Get specific audit entry

## UID/GID Ranges

The system enforces specific UID/GID ranges for different account types:

| Type     | UID/GID Range |
|----------|---------------|
| People   | 1000-60000    |
| System   | 1-999         |
| Service  | 60001-65535   |
| Database | 70000-79999   |

## Database Schema

The database consists of the following tables:

1. **accounts**: Stores user accounts with UIDs
   - id (PK)
   - uid (unique)
   - username (unique)
   - type (people, system, database, service)
   - primary_group_id (FK to groups)
   - created_at, updated_at, deleted_at

2. **groups**: Stores groups with GIDs
   - id (PK)
   - gid (unique)
   - groupname (unique)
   - description
   - type (people, system, database, service)
   - created_at, updated_at, deleted_at

3. **account_groups**: Many-to-many relationship between accounts and groups
   - account_id (PK, FK to accounts)
   - group_id (PK, FK to groups)
   - created_at, updated_at

4. **audit_entries**: Audit log for all actions
   - id (PK)
   - action
   - entity_id
   - entity_type
   - details
   - user_id
   - username
   - ip_address
   - timestamp