-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- Extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create accounts table with audit fields
CREATE TABLE accounts (
    id VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    last_login TIMESTAMP,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    created_by VARCHAR(36) NOT NULL,
    updated_by VARCHAR(36) NOT NULL,
    deleted_by VARCHAR(36)
);

-- Create index on accounts for soft delete queries and username searches
CREATE INDEX idx_accounts_deleted_at ON accounts(deleted_at NULLS FIRST);
CREATE INDEX idx_accounts_username ON accounts(username);
CREATE INDEX idx_accounts_email ON accounts(email);

-- Create groups table
CREATE TABLE groups (
    id VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    created_by VARCHAR(36) NOT NULL,
    updated_by VARCHAR(36) NOT NULL,
    deleted_by VARCHAR(36)
);

-- Create index on groups for soft delete queries and name searches
CREATE INDEX idx_groups_deleted_at ON groups(deleted_at NULLS FIRST);
CREATE INDEX idx_groups_name ON groups(name);

-- Create account_groups table for many-to-many relationship
CREATE TABLE account_groups (
    account_id VARCHAR(36) NOT NULL REFERENCES accounts(id),
    group_id VARCHAR(36) NOT NULL REFERENCES groups(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by VARCHAR(36) NOT NULL,
    updated_by VARCHAR(36) NOT NULL,
    PRIMARY KEY (account_id, group_id)
);

-- Create indexes for account_groups
CREATE INDEX idx_account_groups_account_id ON account_groups(account_id);
CREATE INDEX idx_account_groups_group_id ON account_groups(group_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS account_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS accounts;