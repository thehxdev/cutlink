DROP TABLE IF EXISTS urls;
DROP TABLE IF EXISTS users;


CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    pass_hash VARCHAR(60) NOT NULL UNIQUE,
    isAdmin BOOLEAN NOT NULL CHECK (isAdmin IN (0,1))
);


CREATE TABLE urls (
    id INTEGER PRIMARY KEY,
    target VARCHAR(1024) NOT NULL,
    hash VARCHAR(10) UNIQUE NOT NULL,
    pass_hash VARCHAR(60),
    clicked INTEGER DEFAULT 0,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
