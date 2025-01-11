CREATE TABLE IF NOT EXISTS "users"
(
    "id" SERIAL PRIMARY KEY,
    "name" TEXT,
    "login" TEXT UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL
);

CREATE INDEX idx_users_login ON users(login);
