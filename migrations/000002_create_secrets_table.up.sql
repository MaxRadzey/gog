CREATE TABLE IF NOT EXISTS secrets (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    secret_type VARCHAR(50) NOT NULL,
    data        BYTEA NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_secrets_user_id ON secrets (user_id);
CREATE INDEX idx_secrets_user_id_is_deleted ON secrets (user_id, is_deleted);
