package db

import (
	"database/sql"
	"errors"
	"examblanc/models"
	"fmt"
)

func GetAthletes() ([]models.Athlete, error) {
	var athletes []models.Athlete

	rows, err := Instance.Query("select * from athletes")
	if err != nil {
		return nil, fmt.Errorf("GetAthletes(db): %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var athlete models.Athlete

		err := rows.Scan(&athlete.Id, &athlete.Name, &athlete.Country, &athlete.Age)
		if err != nil {
			return nil, fmt.Errorf("GetAthletes(db): %v", err)
		}

		athletes = append(athletes, athlete)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetAthletes(db): %v", err)
	}

	return athletes, nil
}

func AddAthlete(athlete models.Athlete) error {
	_, err := Instance.Exec("insert into athletes (id,name,country,age) values (?,?,?,?)",
		athlete.Id, athlete.Name, athlete.Country, athlete.Age)

	if err != nil {
		return fmt.Errorf("AddAthlete(db): %v", err)
	}

	return nil
}

func GetAthlete(id string) (*models.Athlete, error) {
	var athlete models.Athlete

	row := Instance.QueryRow("select * from athletes where id=?", id)
	err := row.Scan(&athlete.Id, &athlete.Name, &athlete.Country, &athlete.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetAthleteById(db): %v", err)
	}

	return &athlete, nil
}

func PutAthlete(id string, athlete models.Athlete) error {
	_, err := Instance.Exec("update athletes set name=?,country=?,age=? where id=?",
		athlete.Name, athlete.Country, athlete.Age, athlete.Id)

	if err != nil {
		return fmt.Errorf("PutAthlete(db): %v", err)
	}

	return nil
}

func DeleteAthlete(id string) error {

	res, err := Instance.Exec("delete from athletes where id=?", id)
	if err != nil {
		return fmt.Errorf("GetAthleteById(db): %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("GetAthleteById(db): %v", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
