CREATE TABLE IF NOT EXISTS user_feedback (
    id VARCHAR(200) NOT NULL PRIMARY KEY,
    first_name VARCHAR(200),
    last_name VARCHAR(200),
    email VARCHAR(200),
    job_title VARCHAR(200),
    feedback VARCHAR(1200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
