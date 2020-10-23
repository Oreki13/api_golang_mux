DO $$ BEGIN
  CREATE EXTENSION pgcrypto;
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    status TEXT NOT NULL,
    role_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);

CREATE TABLE roles (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  roles TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (id, roles) VALUES ('00000000-0000-0000-0000-000000000001', 'ADMIN');
INSERT INTO users (id, username, email, role_id, status) VALUES ('00000000-0000-0000-0000-000000000002', 'bambang', 'bambang@getnada.com', '00000000-0000-0000-0000-000000000001', 'ACTIVE');