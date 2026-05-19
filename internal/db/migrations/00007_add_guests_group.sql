-- +goose Up
-- +goose StatementBegin
INSERT INTO groups (name, kind, description) VALUES
  ('All guests', 'builtin', 'Anyone signed in as a guest (no BCC/Azure AD account).');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM groups WHERE name = 'All guests' AND kind = 'builtin';
-- +goose StatementEnd
