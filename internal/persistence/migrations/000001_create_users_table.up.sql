CREATE EXTENSION IF NOT EXISTS "citext" ;

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email CITEXT NOT NULL,
  -- password hash can be null for cases where users registered using oauth
  password_hash TEXT,
  email_verified_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT users_email_unique UNIQUE (email) 
);

-- index for active account email lookups
CREATE INDEX idx_users_email ON users (email);
