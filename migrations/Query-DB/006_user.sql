CREATE TABLE IF NOT EXISTS user_projection (
    id INT PRIMARY KEY,
    username VARCHAR(100),
    email VARCHAR(255),
    is_admin BOOLEAN,
    created_at TIMESTAMP
);