package models

import (
	"time"
    "database/sql"

    "golang.org/x/crypto/bcrypt"
)


type Url struct {
    ID          int
    Target      string
    Hash        string
    Clicked     int
    Created     *time.Time
    UserID      int
}


type User struct {
    ID             int
    UUID_hash    string
    PassHash    string
}


type Conn struct {
    DB *sql.DB
}


func (c *Conn) GetUrl(hash string) (*Url, error) {
    stmt := `SELECT id, target, hash, clicked, created FROM urls WHERE hash = ?`
    url := &Url{}

    err := c.DB.QueryRow(stmt, hash).Scan(&url.ID, &url.Target, &url.Hash, &url.Clicked, &url.Created)
    if err != nil {
        return nil, err
    }

    return url, nil
}


func (c *Conn) GetAllUrls(id int) ([]*Url, error) {
    stmt := `SELECT id, target, hash, clicked, created FROM urls WHERE user_id = ? ORDER BY id DESC`
    urls := []*Url{}

    rows, err := c.DB.Query(stmt, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        url := &Url{}
        err := rows.Scan(&url.ID, &url.Target, &url.Hash, &url.Clicked, &url.Created)
        if err != nil {
            return nil, err
        }

        urls = append(urls, url)
    }

    return urls, nil
}


func (c *Conn) CreateUrl(uid int, target string) (int, string, error) {
    stmt := `INSERT INTO urls (target, hash, user_id) VALUES (?, ?, ?)`

    hashLen := genRandNum(5, 7)
    tHash := genHash(target, hashLen)

    res, err := c.DB.Exec(stmt, target, tHash, uid)
    if err != nil {
        return 0, "", err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, "",err
    }

    return int(id), tHash, nil
}


func (c *Conn) IncrementClicked(hash string) error {
    stmt := `UPDATE urls SET clicked = clicked + 1 WHERE hash = ?`

    _, err := c.DB.Exec(stmt, hash)
    if err != nil {
        return err
    }

    return nil
}


func (c *Conn) DeleteUrl(id int, hash string) error {
    stmt := `DELETE FROM urls WHERE hash = ? AND user_id = ?`

    _, err := c.DB.Exec(stmt, hash, id)
    if err != nil {
        return err
    }
    return nil
}


/*
func (c *Conn) TableIsEmpty(table string) (int, error) {
    var isEmpty int

    stmt := `SELECT CASE WHEN EXISTS(SELECT 1 FROM ?) THEN 0 ELSE 1 END AS IsEmpty`
    res := c.DB.QueryRowx(stmt, table)

    err := res.Scan(&isEmpty)
    if err != nil {
        return -1, err
    }

    return isEmpty, nil
}
*/


func (c *Conn) CreateUser(uuid, password string) error {
    stmt := `INSERT INTO users (uuid, pass_hash) VALUES (?, ?)`

    pass_hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    _, err = c.DB.Exec(stmt, uuid, string(pass_hash))
    if err != nil {
        return err
    }

    return nil
}


func (c *Conn) AuthenticatUser(uuid, pass string) (int, error) {
    var id int
    var tmpPass string
    stmt := `SELECT id, pass_hash FROM users WHERE uuid = ?`

    err := c.DB.QueryRow(stmt, uuid).Scan(&id, &tmpPass)
    if err != nil {
        return 0, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(tmpPass), []byte(pass))
    if err != nil {
        return 0, err
    }

    return id, nil
}


func (c *Conn) DeleteUser(id int) error {
    stmt1 := `DELETE FROM users WHERE id = ?`
    stmt2 := `DELETE FROM urls WHERE user_id = ?`

    _, err := c.DB.Exec(stmt1, id)
    if err != nil {
        return err
    }

    _, err = c.DB.Exec(stmt2, id)
    if err != nil {
        return err
    }

    return nil
}
