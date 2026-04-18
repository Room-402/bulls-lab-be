CREATE TABLE IF NOT EXISTS watchlists (
    id          SERIAL PRIMARY KEY,
    active      BOOLEAN NOT NULL DEFAULT true,
    name        VARCHAR(255) NOT NULL UNIQUE,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_watchlists_user_id ON watchlists(user_id);