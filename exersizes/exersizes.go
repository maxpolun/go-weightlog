package exersizes

import (
	"../util"
)

type Exersize struct {
	Id   int64
	Name string
}

func GetByName(name string, db util.DB) (*Exersize, error) {
	row := db.QueryRow("SELECT id, name FROM exersizes WHERE name=$1;", name)
	e := new(Exersize)
	err := row.Scan(&e.Id, &e.Name)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Exersize) Save(db util.DB) error {
	_, err := db.Exec("INSERT INTO exersizes (name) VALUES ($1);", e.Name)
	if err != nil {
		return err
	}
	e2, err := GetByName(e.Name, db)
	if err != nil {
		return err
	}
	e.Id = e2.Id
	return nil
}

func GetAll(db util.DB) ([]*Exersize, error) {
	rows, err := db.Query("SELECT id, name FROM exersizes;")
	if err != nil {
		return nil, err
	}
	es := make([]*Exersize, 0, 10)
	for rows.Next() {
		e := &Exersize{}
		err := rows.Scan(&e.Id, &e.Name)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}
