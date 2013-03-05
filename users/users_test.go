package users

import (
	"../util"
	"bytes"
	"database/sql"
	"github.com/jameskeane/bcrypt"
	"testing"
)

var db *sql.DB

func init() {
	db = util.GetTestDb()
}

func TestGetByEmailExists(t *testing.T) {
	t.Parallel()
	txn, _ := db.Begin()
	defer txn.Rollback()
	txn.Exec("INSERT INTO users (email, pw_hash) VALUES ('test@test.com','randomjunk');")
	user, err := GetByEmail(txn, "test@test.com")
	if err != nil {
		t.Errorf("error while using GetByEmail: %v", err)
	}
	if user.Email != "test@test.com" {
		t.Errorf("expected 'test@test.com' when using GetByEmail, got %v instead", user.Email)
	}

}
func TestGetByEmailDoesNotExist(t *testing.T) {
	t.Parallel()
	txn, _ := db.Begin()
	user, err := GetByEmail(txn, "test@test.com")
	if err != sql.ErrNoRows {
		t.Errorf("expected no rows error, got %v instead", err)
	}
	if user != nil {
		t.Errorf("expected to get nil with nonexistant user, got %v instead", user)
	}
}

func TestVerifyCorrect(t *testing.T) {
	t.Parallel()
	hash, _ := bcrypt.HashBytes([]byte("password"))
	u := &User{Id: 1, Email: "test@test.com", PwHash: hash}
	if !u.Verify("password") {
		t.Errorf("Expected user password and hash to match, got user.PwHash = '%v', password = 'password'", u.PwHash)
	}
}
func TestVerifyIncorrect(t *testing.T) {
	t.Parallel()
	hash, _ := bcrypt.HashBytes([]byte("password"))
	u := &User{Id: 1, Email: "test@test.com", PwHash: hash}
	if u.Verify("wrongPassword") {
		t.Errorf("Expected user password and hash to NOT match, got user.PwHash = '%v', password = 'wrongPassword'", u.PwHash)
	}
}

func TestSave(t *testing.T) {
	t.Parallel()
	txn, _ := db.Begin()
	defer txn.Rollback()
	u := &User{
		Email:  "test@test.com",
		PwHash: []byte("randomjunk")}
	u.Save(txn)

	row := txn.QueryRow("SELECT email, pw_hash FROM users WHERE email = $1;", "test@test.com")
	email := ""
	hash := []byte{}
	row.Scan(&email, &hash)
	if email != u.Email {
		t.Errorf("expected %v in the database after creating a user, got %v instead.",
			u.Email,
			email)
	}
	if bytes.Compare(hash, u.PwHash) == 0 {
		t.Errorf("expected %v in the database after creating a user, got %v instead.",
			u.PwHash,
			hash)
	}

}
func TestSaveAlreadyExists(t *testing.T) {
	t.Parallel()
	txn, _ := db.Begin()
	defer txn.Rollback()
	u := &User{
		Email:  "test@test.com",
		PwHash: []byte("randomjunk")}
	err := u.Save(txn)
	if err != nil {
		t.Errorf("error when inserting user: %v", err)
	}
	originalId := u.Id
	u.Email = "test12345@test.com"
	err = u.Save(txn)
	if err != nil {
		t.Errorf("expected no errors when inserting user: %v", err)
	}

	row := txn.QueryRow("SELECT email, pw_hash FROM users WHERE email = $1;", u.Email)
	email := ""
	hash := []byte{}
	row.Scan(&email, &hash)
	if email != u.Email {
		t.Errorf("expected %v in the database after creating a user, got %v instead.",
			u.Email,
			email)
	}
	if originalId != u.Id {
		t.Errorf("user id changed when updating user. Original: %v, new: %v", originalId, u.Id)
	}
}
func TestNew(t *testing.T) {
	email := "test@test.com"
	password := "password"
	u, err := New(email, password)
	if err != nil {
		t.Errorf("error constructing new User: %v", err)
	}
	if !bcrypt.MatchBytes([]byte(password), u.PwHash) {
		t.Errorf("expected hash to match, got %v", u.PwHash)
	}
}
