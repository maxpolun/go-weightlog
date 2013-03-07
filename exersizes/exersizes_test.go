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

func findExersize(e *Exersize, slice []*Exersize) bool {
	for i := range slice {
		if slice[i].Name == e.Name {
			return true
		}
	}
	return false
}

func TestGetAll(t *testing.T) {
	db := util.GetTestDb()
	txn, _ := db.Begin()
	defer txn.Rollback()

	es := []*Exersize{
		&Exersize{Name: "squats"},
		&Exersize{Name: "press"},
		&Exersize{Name: "deadlift"},
		&Exersize{Name: "bench"},
		&Exersize{Name: "clean"}}
	for _, e := range es {
		err := e.Save(txn)
		if err != nil {
			t.Errorf("error saving exersize %v: %v", e, err)
		}
	}
	all, err := GetAll(txn)
	if err != nil {
		t.Errorf("error getting all exersizes: %v", err)
	}
	if len(all) != len(es) {
		t.Errorf("differnt number of exersizes returned than saved: got %v, expected %v", len(all), len(es))
	}
	for _, e := range all {
		if !findExersize(e, es) {
			t.Errorf("unable to find %v in %v.", e, es)
		}
	}

}
