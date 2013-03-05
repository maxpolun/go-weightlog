package sessions

import (
	"../users"
	"../util"
	"bytes"
	"database/sql"
	"testing"
)

var db *sql.DB

func init() {
	db = util.GetTestDb()
}

func TestNewSessionWorks(t *testing.T) {
	t.Parallel()
	session1, err := New(0)
	if err != nil {
		t.Errorf("got an error while attempting to create a new session: %v", err)
	}
	if len(session1.Id) == 0 {
		t.Errorf("blank session id when creating a new session")
	}
	//t.Errorf("created session {%s, %v}", encodeBase64(session1.Id), session1.UserId)

}
func TestNewSessionUniqueIds(t *testing.T) {
	t.Parallel()
	session1, err := New(0)
	if err != nil {
		t.Errorf("got an error while attempting to create a new session: %v", err)
	}
	session2, err := New(1)
	if err != nil {
		t.Errorf("got an error while attempting to create a new session: %v", err)
	}
	if bytes.Compare(session1.Id, session2.Id) == 0 {
		t.Errorf("random session ids created were not unique! %v == %v", session1.Id, session2.Id)
	}

}

func TestSaveWorks(t *testing.T) {
	t.Parallel()
	txn, err := db.Begin()
	defer txn.Rollback()
	if err != nil {
		t.Errorf("got error while creating transaction: %v", err)
	}
	u, err := users.New("test@test.com", "testpassword")
	err = u.Save(txn)
	if err != nil {
		t.Errorf("got error while saving user: %v", err)
	}
	t.Logf("saved user %d", u.Id)
	s, _ := New(u.Id)
	err = s.Save(txn)
	if err != nil {
		t.Errorf("got error while saving session: %v", err)
	}
	row := txn.QueryRow("SELECT id, user_id FROM sessions WHERE id = $1;", s.Id)
	id := make([]byte, SESSION_BYTES)
	user_id := int64(0)
	err = row.Scan(&id, &user_id)
	if err != nil {
		t.Errorf("got error when scanning session row: %v", err)
	}
	if bytes.Compare(id, s.Id) != 0 || user_id != s.UserId {
		t.Errorf("returned session differs from saved one! saved = %v returned = {%s, %v}", s, id, user_id)
	}

}
