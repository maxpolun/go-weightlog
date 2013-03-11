package users

import (
	"../util"
	"errors"
	"fmt"
	"github.com/jameskeane/bcrypt"
)

type User struct {
	Id     int64
	Email  string
	PwHash []byte
	saved  bool
}

var ErrUserDoesNotExist error = errors.New("User does not exist")

func GetByEmail(db util.DB, email string) (u *User, err error) {
	row := db.QueryRow("SELECT id, email, pw_hash FROM users WHERE email = $1;", email)

	u = &User{}
	err = row.Scan(&u.Id, &u.Email, &u.PwHash)
	if err != nil {
		return nil, err
	}
	u.saved = true
	return u, nil
}
func GetById(db util.DB, id int64) (u *User, err error) {
	row := db.QueryRow("SELECT id, email, pw_hash FROM users WHERE id = $1;", id)

	u = &User{}
	err = row.Scan(&u.Id, &u.Email, &u.PwHash)
	if err != nil {
		return nil, err
	}
	u.saved = true
	return u, nil
}
func (u *User) String() string {
	return fmt.Sprintf("User{id: %v, email:%v, pw_hash:%v, saved: %v}", u.Id,
		u.Email, u.PwHash, u.saved)
}
func New(email, password string) (*User, error) {
	u := &User{}
	u.Email = email
	hash, err := bcrypt.HashBytes([]byte(password))
	if err != nil {
		return nil, err
	}
	u.PwHash = hash
	return u, nil
}

func (u *User) Verify(password string) bool {
	return bcrypt.MatchBytes([]byte(password), u.PwHash)
}

func (u *User) saveNew(db util.DB) error {
	_, err := db.Exec("INSERT INTO users (email, pw_hash) VALUES ($1, $2);", u.Email, u.PwHash)
	if err != nil {
		return err
	}
	u2, err := GetByEmail(db, u.Email)
	if err != nil {
		return err
	}
	u.saved = true
	u.Id = u2.Id
	return nil
}

func (u *User) update(db util.DB) error {
	_, err := db.Exec("UPDATE users SET email=$1, pw_hash=$2 WHERE id=$3", u.Email, u.PwHash, u.Id)
	return err
}

func (u *User) Save(db util.DB) error {
	if u.saved {
		return u.update(db)
	} else {
		return u.saveNew(db)
	}
	return nil
}
