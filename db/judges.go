package db

import (
	"database/sql"
	"errors"
	"examblanc/models"
	"fmt"
)

func GetJudges(page *int) ([]models.Judge, error) {
	judges := []models.Judge{}

	var rows *sql.Rows
	var err error

	if page != nil {
		rows, err = Instance.Query("select id,name,phone,experience_years from judges limit ?,10", (*page - 1) * 10)
	} else {
		rows, err = Instance.Query("select id,name,phone,experience_years from judges")
	}

	if err != nil {
		return nil, fmt.Errorf("GetJudges(db): %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var judge models.Judge

		err := rows.Scan(&judge.Id, &judge.Name, &judge.Phone, &judge.ExperienceYears)
		if err != nil {
			return nil, fmt.Errorf("GetJudges(db): %v", err)
		}

		judges = append(judges, judge)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetJudges(db): %v", err)
	}

	return judges, nil
}

func GetJudgeWithName(name string) (*models.Judge, error) {
	var judge models.Judge

	row := Instance.QueryRow("select id,name,phone,experience_years from judges where name=?", name)
	err := row.Scan(&judge.Id, &judge.Name, &judge.Phone, &judge.ExperienceYears)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetJudges(db): %v", err)
	}

	return &judge, nil
}

func PostJudge(judge models.JudgeWithPassword) error {
	_, err := Instance.Exec("insert into judges (id,name,password,phone,experience_years) values (?,?,?,?,?)",
		judge.Id, judge.Name, judge.Password, judge.Phone, judge.ExperienceYears)

	if err != nil {
		return fmt.Errorf("PostJudge(db): %v", err)
	}

	return nil
}

func GetJudge(id string) (*models.Judge, error) {
	var judge models.Judge

	row := Instance.QueryRow("select id,name,phone,experience_years from judges where id=?", id)
	err := row.Scan(&judge.Id, &judge.Name, &judge.Phone, &judge.ExperienceYears)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetJudge(db): %v", err)
	}

	return &judge, nil
}

func DeleteJudge(id string) error {

	res, err := Instance.Exec("delete from judges where id=?", id)
	if err != nil {
		return fmt.Errorf("DeleteJudge(db): %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteJudge(db): %v", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
