package sets

import (
	"../exersizes"
	"../users"
	"../util"
	"testing"
)

func TestSaveWorks(t *testing.T) {
	db := util.GetTestDb()
	txn, _ := db.Begin()
	defer txn.Rollback()
	e, err := exersizes.GetByName("squat", db)
	if err != nil {
		t.Errorf("Error getting exersize:", err)
	}
	u, err := users.New("test@test.com", "password123")
	err = u.Save(txn)
	if err != nil {
		t.Errorf("Error saving user:", err)
	}

	set := New(e.Id, 5, u.Id, 250, "lb")
	err = set.Save(txn)
	if err != nil {
		t.Errorf("error saving set: %v", err)
	}
	row := txn.QueryRow("SELECT id, completed_at, exersize_id, reps, user_id, weight, unit, notes FROM sets WHERE id=$1;",
		set.Id)
	set2 := new(Set)
	err = row.Scan(&set2.Id, &set2.CompletedAt, &set2.ExersizeId, &set2.Reps, &set2.UserId, &set2.Weight, &set2.Unit, &set2.Notes)
	if err != nil {
		t.Errorf("error scanning set: %v", err)
	}
	if set.Id != set2.Id || set.ExersizeId != set2.ExersizeId || set.Reps != set2.Reps || set.Weight != set2.Weight || set.Unit != set2.Unit {
		t.Errorf("returned set does not equal expected set: expected %v, got %v", set, set2)
	}
}

func TestGetByUserId(t *testing.T) {
	db := util.GetTestDb()
	txn, _ := db.Begin()
	defer txn.Rollback()
	user1, _ := users.New("test1@example.com", "testpassword")
	user1.Save(txn)
	user2, _ := users.New("test2@example.com", "testpassword")
	user2.Save(txn)

	e, err := exersizes.GetByName("squat", txn)

	for i := 0; i < 10; i++ {
		s := New(e.Id, 5, user1.Id, 50, "lb")
		err := s.Save(txn)
		if err != nil {
			t.Errorf("got error saving set:%v", err)
		}
	}
	s := New(e.Id, 5, user2.Id, 50, "lb")
	err = s.Save(txn)
	if err != nil {
		t.Errorf("got error saving set:%v", err)
	}
	users, err := GetByUserId(user1.Id, txn)
	if err != nil {
		t.Errorf("got error getting sets:%v", err)
	}
	if len(users) != 10 {
		t.Errorf("Expected 10 users, got %v instead", len(users))
	}
}
