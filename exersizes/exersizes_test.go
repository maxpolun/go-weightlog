package exersizes

import (
	"../util"
	"testing"
)

func TestSaveWorks(t *testing.T) {
	db := util.GetTestDb()
	txn, _ := db.Begin()
	defer txn.Rollback()
	e := &Exersize{
		Name: "squats"}

	err := e.Save(txn)
	if err != nil {
		t.Errorf("got error saving exersize: %v", err)
	}
	row := txn.QueryRow("SELECT id, name FROM exersizes WHERE id=$1;", e.Id)
	e2 := new(Exersize)
	err = row.Scan(&e2.Id, &e2.Name)
	if err != nil {
		t.Errorf("got error reading exersize from database: %v", err)
	}
	if e.Id != e2.Id || e.Name != e2.Name {
		t.Errorf("returned exersize differs from expected: expected = %v, got %v", e, e2)
	}
}
