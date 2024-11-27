package models

import (
	"golang.org/x/crypto/bcrypt"
)

func (c *Conn) CreateUser(uuid, password string) error {
	stmt := `INSERT INTO users (uuid, pass_hash, isAdmin) VALUES (?, ?, ?)`

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = c.DB.Exec(stmt, uuid, string(passHash), false)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) AuthenticatUser(uuid, pass string) (int, bool, error) {
	var id int
	var tmpPass string
	var isAdmin bool
	stmt := `SELECT id, pass_hash, isAdmin FROM users WHERE uuid = ?`

	err := c.DB.QueryRow(stmt, uuid).Scan(&id, &tmpPass, &isAdmin)
	if err != nil {
		return 0, false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(tmpPass), []byte(pass))
	if err != nil {
		return 0, false, err
	}

	return id, isAdmin, nil
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
