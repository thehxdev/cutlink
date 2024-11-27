package models

import (
	"database/sql"
	"time"
)

const DBschema = `
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    pass_hash VARCHAR(60) NOT NULL UNIQUE,
    isAdmin BOOLEAN NOT NULL CHECK (isAdmin IN (0,1))
);

CREATE TABLE urls (
    id INTEGER PRIMARY KEY,
    target VARCHAR(2048) NOT NULL,
    hash VARCHAR(10) UNIQUE NOT NULL,
    pass_hash VARCHAR(60),
    clicked INTEGER DEFAULT 0,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);`

type Url struct {
	ID       int
	Target   string
	Hash     string
	PassHash string
	Clicked  int
	Created  *time.Time
	UserID   int
}

type User struct {
	ID       int
	UUID     string
	PassHash string
	IsAdmin  bool
}

type Conn struct {
	DB *sql.DB
}

func (c *Conn) MigrateDB() error {
	_, err := c.DB.Exec(DBschema)
	return err
}
