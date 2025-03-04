CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255),
    name VARCHAR(255),
    description text,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);