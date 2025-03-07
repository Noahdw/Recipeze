CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255),
    name VARCHAR(255),
    description text,
    image_url VARCHAR(255),
    likes INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(128) NOT NULL,
    name VARCHAR(128),
    image_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);