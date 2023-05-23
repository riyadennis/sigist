CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    job_title TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
