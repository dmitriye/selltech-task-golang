-- +goose Up
-- +goose StatementBegin
-- auto-generated definition
CREATE TABLE IF NOT EXISTS sdn_entries
(
    id                    uuid default gen_random_uuid() not null primary key,
    created_at            timestamptz default current_timestamp,
    updated_at            timestamptz default current_timestamp,
    uid                   integer                        not null,
    -- published_at          timestamptz,
    first_name            varchar(255),
    last_name             varchar(255)
);
CREATE UNIQUE INDEX idx_sdn_entries_uid ON sdn_entries (uid);
CREATE INDEX idx_sdn_entries_first_name on sdn_entries (first_name);
CREATE INDEX idx_sdn_entries_last_name on sdn_entries (last_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sdn_entries;
-- +goose StatementEnd


