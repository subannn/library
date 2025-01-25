CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE,
    lyrics TEXT,
    link TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);
