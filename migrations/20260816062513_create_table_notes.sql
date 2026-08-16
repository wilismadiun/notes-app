-- +goose Up
CREATE TABLE notes (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    create_at TIMESTAMPTZ NOT NULL,
    update_at TIMESTAMPTZ NOT NULL,
    owner VARCHAR(255) NOT NULL,

    CONSTRAINT fk_notes_owner
        FOREIGN KEY (owner)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE notes;
