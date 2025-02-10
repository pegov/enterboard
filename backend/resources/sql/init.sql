CREATE TABLE IF NOT EXISTS boards (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    board_id INTEGER NOT NULL REFERENCES boards(id),
    thread_id INTEGER NOT NULL REFERENCES posts(id),
    title TEXT,
    message TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS threads (
    board_id INTEGER REFERENCES boards(id),
    post_id INTEGER REFERENCES posts(id),
    PRIMARY KEY (board_id, post_id)
);

