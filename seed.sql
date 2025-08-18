CREATE SCHEMA IF NOT EXISTS eagle;
SET SCHEMA 'eagle';

DROP TABLE IF EXISTS user_accounts;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS user_verification_tokens;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS user_status;

CREATE TYPE user_status AS ENUM ('awaiting_verification', 'email_verified', 'active', 'suspended');
CREATE TYPE account_type AS ENUM ('personal', 'business');


CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       phone_number VARCHAR(20),
                       password_hash TEXT, -- nullable until verification
                       status user_status NOT NULL DEFAULT 'suspended',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/* TODO CREATE USER_AUDIT TABLE */

CREATE TABLE addresses (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           line1 VARCHAR(255) NOT NULL,
                           line2 VARCHAR(255),
                           line3 VARCHAR(255),
                           town VARCHAR(100) NOT NULL,
                           county VARCHAR(100),
                           postcode VARCHAR(20) NOT NULL,
                           created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_addresses_user_id ON addresses(user_id);

/* TODO CREATE USER_ADDRESS_AUDIT TABLE */

/* TODO CREATE NECESSARY VIEW OF USER PLUS ADDRESS */

CREATE TABLE user_verification_tokens (
                                          token UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          expires_at TIMESTAMPTZ NOT NULL,
                                          used_at TIMESTAMPTZ,
                                          UNIQUE (user_id)
);



INSERT INTO users (name, email, phone_number)
VALUES
    ('Alice Smith', 'alice@example.com', '+447123456789'),
    ('Bob Johnson', 'bob@example.com', '+447234567890'),
    ('Carol Davis', 'carol@example.com', '+447345678901');


INSERT INTO addresses (user_id, line1, line2, town, county, postcode)
SELECT id, '123 Main Street', NULL, 'London', 'Greater London', 'SW1A 1AA'
FROM users WHERE email = 'alice@example.com';

INSERT INTO addresses (user_id, line1, line2, town, county, postcode)
SELECT id, '456 Oak Avenue', 'Flat 2B', 'Manchester', 'Greater Manchester', 'M1 2AB'
FROM users WHERE email = 'bob@example.com';

INSERT INTO addresses (user_id, line1, line2, line3, town, county, postcode)
SELECT id, '789 Pine Road', NULL, NULL, 'Bristol', 'City of Bristol', 'BS1 3CD'
FROM users WHERE email = 'carol@example.com';

CREATE TABLE accounts (
                         account_number     CHAR(8) PRIMARY KEY,  -- fixed 8-digit account number
                         sort_code          CHAR(8) NOT NULL,     -- e.g. "10-10-10"
                         name               VARCHAR(100) NOT NULL,
                         account_type       account_type NOT NULL,
                         balance            NUMERIC(15,2) NOT NULL DEFAULT 0.00, -- allows for large values, 2 decimal places
                         currency           CHAR(3) NOT NULL,     -- ISO currency code like GBP, USD
                         created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/* TODO CREATE ACCOUNT HISTORY TABLES */

CREATE TABLE user_accounts (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               account_number CHAR(8) NOT NULL REFERENCES accounts(account_number) ON DELETE CASCADE,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               UNIQUE (user_id, account_number) -- prevents duplicate user/account pairs
);


