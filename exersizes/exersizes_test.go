package exersizes

import (
	"../util"
	"testing"
)

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

	txn.Exec("TRUNCATE exersizes CASCADE;")

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
