CREATE TABLE IF NOT EXISTS roles (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    level INT NOT NULL DEFAULT 0
);

INSERT INTO roles (id, name, description, level) VALUES
(1, 'user', 'A user can create posts and comments', 1),
(2, 'moderator', 'A moderator can update other users posts', 2),
(3, 'admin', 'An admin can update and delete other users posts', 3);