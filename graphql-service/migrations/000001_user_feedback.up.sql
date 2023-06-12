CREATE TABLE IF NOT EXISTS user_feedback (
    id TEXT NOT NULL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    job_title TEXT,
    feedback BLOB,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
