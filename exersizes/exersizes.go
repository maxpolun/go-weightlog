package exersizes

import (
	"../util"
)

type Exersize struct {
	Name string `json:name`
}

func (e *Exersize) Save(db util.DB) error {
	_, err := db.Exec("INSERT INTO exersizes (name) VALUES ($1);", e.Name)
	if err != nil {
		return err
	}
	return nil
}

func GetAll(db util.DB) ([]*Exersize, error) {
	rows, err := db.Query("SELECT name FROM exersizes;")
	if err != nil {
		return nil, err
	}
	es := make([]*Exersize, 0, 10)
	for rows.Next() {
		e := &Exersize{}
		err := rows.Scan(&e.Name)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}
