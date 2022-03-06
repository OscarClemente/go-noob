CREATE TABLE IF NOT EXISTS friends(
    id SERIAL PRIMARY KEY,
    username INTEGER,
    friend INTEGER
);