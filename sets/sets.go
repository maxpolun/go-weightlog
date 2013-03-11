package sets

import (
	"../util"
	"time"
)

// this is 
type Set struct {
	Id          int64
	CompletedAt time.Time
	Exersize    string
	Reps        int
	UserId      int64
	Weight      int
	Unit        string
	Notes       string
}

func New(exersize string, reps int, userid int64, weight int, unit string) *Set {
	// database time is in ms, set time to be prescise only to the second 
	now := time.Now().UTC()
	return &Set{
		CompletedAt: time.Unix(now.Unix(), 0),
		Exersize:    exersize,
		Reps:        reps,
		UserId:      userid,
		Weight:      weight,
		Unit:        unit,
		Notes:       ""}
}

func (s *Set) Save(db util.DB) error {
	if s.Id > 0 {
		_, err := db.Exec(`UPDATE sets SET 
			completed_at=$1, exersize = $2, reps=$3, user_id=$4, weight=$5, uit=$6, notes=$7
			WHERE id=$8;`,
			s.CompletedAt,
			s.Exersize,
			s.Reps,
			s.UserId,
			s.Weight,
			s.Unit,
			s.Notes,
			s.Id)
		return err
	}
	_, err := db.Exec(`INSERT INTO sets(completed_at, exersize, reps, user_id, weight, unit, notes) 
		VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		s.CompletedAt,
		s.Exersize,
		s.Reps,
		s.UserId,
		s.Weight,
		s.Unit,
		s.Notes)
	if err != nil {
		return err
	}
	id, err := util.GetLastId(db, "sets")
	s.Id = id
	return err
}

func GetByUserId(user_id int64, db util.DB) ([]*Set, error) {
	rows, err := db.Query(`SELECT id, completed_at, exersize, reps, user_id, weight, unit, notes FROM
		sets WHERE
		user_id = $1;`, user_id)
	if err != nil {
		return nil, err
	}
	sets := []*Set{}
	for rows.Next() {
		set := new(Set)
		err := rows.Scan(&set.Id, &set.CompletedAt, &set.Exersize, &set.Reps, &set.UserId, &set.Weight, &set.Unit, &set.Notes)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	return sets, nil
}
