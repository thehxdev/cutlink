package models

import (
	"net/url"

	"github.com/thehxdev/cutlink/rand"
	"golang.org/x/crypto/bcrypt"
)

func encodeUrl(target string) string {
	// return base64.URLEncoding.EncodeToString([]byte(url.QueryEscape(target)))
	return url.QueryEscape(target)
}

func decodeUrl(target string) (string, error) {
	// t, err := base64.URLEncoding.DecodeString(target)
	// if err != nil {
	//     // https://stackoverflow.com/a/40673073/19005972
	// 	return string(t[:]), nil
	// }
	// return "", err
	return url.QueryUnescape(target)
}

func (c *Conn) GetUrl(hash string) (*Url, error) {
	stmt := `SELECT id, target, hash, pass_hash, clicked, created FROM urls WHERE hash = ?`
	url := &Url{}

	err := c.DB.QueryRow(stmt, hash).Scan(&url.ID, &url.Target, &url.Hash, &url.PassHash, &url.Clicked, &url.Created)
	if err != nil {
		return nil, err
	}

	url.Target, err = decodeUrl(url.Target)
	return url, err
}

func (c *Conn) GetAllUrls(id int) ([]*Url, error) {
	stmt := `SELECT id, target, hash, pass_hash, clicked, created FROM urls WHERE user_id = ? ORDER BY id DESC`
	urls := []*Url{}

	rows, err := c.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		url := &Url{}
		err := rows.Scan(&url.ID, &url.Target, &url.Hash, &url.PassHash, &url.Clicked, &url.Created)
		if err != nil {
			return nil, err
		}

		url.Target, err = decodeUrl(url.Target)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (c *Conn) CreateUrl(uid int, target, password string) (int, string, error) {
	var passHash []byte = nil
	var err error
	target = encodeUrl(target)

	stmt := `INSERT INTO urls (target, hash, pass_hash, user_id) VALUES (?, ?, ?, ?)`

	hashLen := rand.GenRandNum(5, 7)
	tHash := rand.GenRandString(hashLen)
	if password != "" {
		passHash, err = bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return 0, "", err
		}
	}

	res, err := c.DB.Exec(stmt, target, tHash, string(passHash), uid)
	if err != nil {
		return 0, "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, "", err
	}

	return int(id), tHash, nil
}

func (c *Conn) IncrementClicked(hash string) error {
	stmt := `UPDATE urls SET clicked = clicked + 1 WHERE hash = ?`

	_, err := c.DB.Exec(stmt, hash)
	return err
}

func (c *Conn) DeleteUrl(id int, hash string) error {
	stmt := `DELETE FROM urls WHERE hash = ? AND user_id = ?`

	_, err := c.DB.Exec(stmt, hash, id)
	return err
}
