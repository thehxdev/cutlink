package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)


type Url struct {
    ID          int         `json:"id"`
    Target      string      `json:"target"`
    Hash        string      `json:"hash"`
    Clicked     int         `json:"clicked"`
    Created     *time.Time  `json:"created"`
}


type Urls struct {
    DB *sqlx.DB
}


func (u *Urls) Get(hash string) (*Url, error) {
    stmt := `SELECT id, target, hash, clicked, created FROM urls WHERE hash = ?`
    url := &Url{}

    err := u.DB.QueryRowx(stmt, hash).Scan(&url.ID, &url.Target, &url.Hash, &url.Clicked, &url.Created)
    if err != nil {
        return nil, err
    }

    return url, nil
}


func (u *Urls) GetAll() ([]*Url, error) {
    stmt := `SELECT id, target, hash, clicked, created FROM urls ORDER BY id DESC`
    urls := []*Url{}

    rows, err := u.DB.Queryx(stmt)
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


func (u *Urls) Create(target string) (int, string, error) {
    stmt := `INSERT INTO urls (target, hash) VALUES (?, ?)`

    hashLen := genRandNum(5, 7)
    // hashLen := 5
    tHash := genHash(target, hashLen)

    res, err := u.DB.Exec(stmt, target, tHash)
    if err != nil {
        return 0, "", err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, "",err
    }

    return int(id), tHash, nil
}


func (u *Urls) IncrementClicked(hash string) error {
    stmt := `UPDATE urls SET clicked = clicked + 1 WHERE hash = ?`

    _, err := u.DB.Exec(stmt, hash)
    if err != nil {
        return err
    }

    return nil
}


func (u *Urls) Delete(hash string) error {
    stmt := `DELETE FROM urls WHERE hash = ?`

    _, err := u.DB.Exec(stmt, hash)
    if err != nil {
        return err
    }
    return nil
}


func (u *Urls) TableIsEmpty() (int, error) {
    var isEmpty int

    stmt := `SELECT CASE WHEN EXISTS(SELECT 1 FROM urls) THEN 0 ELSE 1 END AS IsEmpty`
    res := u.DB.QueryRowx(stmt)

    err := res.Scan(&isEmpty)
    if err != nil {
        return -1, err
    }

    return isEmpty, nil
}
