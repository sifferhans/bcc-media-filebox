-- +goose Up
-- +goose StatementBegin
CREATE TABLE targets (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL UNIQUE,
    path       TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE groups (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL UNIQUE,
    kind        TEXT NOT NULL CHECK (kind IN ('builtin', 'custom')),
    description TEXT NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE group_members (
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    email    TEXT NOT NULL,
    PRIMARY KEY (group_id, email)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE grants (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    principal_kind  TEXT NOT NULL CHECK (principal_kind IN ('user', 'group')),
    principal_value TEXT NOT NULL,
    admin           INTEGER NOT NULL DEFAULT 0,
    all_targets     INTEGER NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (principal_kind, principal_value)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE grant_targets (
    grant_id  INTEGER NOT NULL REFERENCES grants(id) ON DELETE CASCADE,
    target_id INTEGER NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    PRIMARY KEY (grant_id, target_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO groups (name, kind, description) VALUES
  ('All BCC members',         'builtin', 'Everyone with a BCC Login account.'),
  ('All bcc.media employees', 'builtin', 'Staff in the bcc.media Azure AD tenant.');
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE uploads ADD COLUMN target_name TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE uploads DROP COLUMN target_name;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS grant_targets;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS grants;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS group_members;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS targets;
-- +goose StatementEnd
