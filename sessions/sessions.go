package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

const SESSION_BYTES = 64

type Session struct {
	Id        []byte
	UserId    int64
	CreatedAt time.Time
}

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func New(uid int64) (s *Session, err error) {
	bytes := make([]byte, SESSION_BYTES)
	n, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	if n != len(bytes) {
		return nil, errors.New("Not enough random bytes generated")
	}
	return &Session{
		Id:     encodeBase64(bytes),
		UserId: uid}, nil
}
func encodeBase64(in []byte) []byte {
	out_bytes := base64.StdEncoding.EncodedLen(SESSION_BYTES)
	out := make([]byte, out_bytes)
	base64.StdEncoding.Encode(out, in)
	return out
}
func (s *Session) Save(db DB) error {
	_, err := db.Exec("INSERT INTO sessions (id, user_id) VALUES ($1, $2);", s.Id, s.UserId)
	return err
}

func (s *Session) String() string {
	return fmt.Sprintf("Session{id:'%s', user_id:%v, created_at:%v}", s.Id, s.UserId, s.CreatedAt)
}
