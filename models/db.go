package models

import (
    "fmt"
    "time"
    "crypto/sha256"

    "github.com/jmoiron/sqlx"
)


type Url struct {
    ID          int
    Target      string
    Hash        string
    Clicked     int
    Created     *time.Time
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

    sha := sha256.New()
    sha.Write([]byte(target))
    hashed := fmt.Sprintf("%x", sha.Sum(nil))

    res, err := u.DB.Exec(stmt, target, hashed[0:10])
    if err != nil {
        return 0, "", err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, "",err
    }

    return int(id), hashed[0:10], nil
}


func (u *Urls) IncrementClicked(hash string) error {
    stmt := `UPDATE urls SET clicked = clicked + 1 WHERE hash = ?`

    _, err := u.DB.Exec(stmt, hash)
    if err != nil {
        return err
    }

    // count, err := res.RowsAffected()
    // _, err = res.RowsAffected()
    // if err != nil {
    //     return err
    // }

    return nil
}


func (u *Urls) Delete(hash string) error {
    stmt := `DELETE FROM urls WHERE hash = ?`

    _, err := u.DB.Exec(stmt, hash)
    if err != nil {
        return err
    }
    // _, err = res.RowsAffected()
    return nil
}
