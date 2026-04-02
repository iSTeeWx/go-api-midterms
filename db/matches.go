package db

import (
	"database/sql"
	"errors"
	"examblanc/models"
	"fmt"
	"time"
)

func GetCountMatchesByJudgeBetweenDates(id string, start time.Time, end time.Time) (int, error) {

	count := 0

	err := Instance.QueryRow("select count(id) from matches where judge_id=? and date>=? and date<=?",
		id, start, end).
		Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetCountMatchesByJudgeBetweenDates(db): %v", err)
	}

	return count, nil
}

func GetCountMatchesBetweenAthletesInDay(athlete1 string, athlete2 string, day time.Time) (int, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())

	count := 0

	fmt.Println(start)
	fmt.Println(end)

	err := Instance.QueryRow("select count(id) from matches where (athlete1_id=? and athlete2_id=?) or (athlete2_id=? and athlete1_id=?) and date>=? and date<=?",
		athlete1, athlete2, athlete1, athlete2, start, end).
		Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetCountMatchesBetweenAthletesInDay(db): %v", err)
	}

	return count, nil
}

func PostMatch(match models.Match) error {

	_, err := Instance.Exec(`insert into matches
													(id,athlete1_id,athlete2_id,judge_id,date,score1,score2)
													value (?,?,?,?,FROM_UNIXTIME(?),?,?)`,
		*match.Id,
		*match.Athlete1Id,
		*match.Athlete2Id,
		*match.JudgeId,
		*match.Date,
		*match.Score1,
		*match.Score2)

	if err != nil {
		return fmt.Errorf("PostMatch(db): %v", err)
	}

	return nil
}

func GetMatchesOfJudge(id string) ([]models.Match, error) {
	matches := []models.Match{}

	rows, err := Instance.Query("select id,athlete1_id,athlete2_id,judge_id,unix_timestamp(date) as date,score1,score2 from matches where judge_id=?", id)
	if err != nil {
		return nil, fmt.Errorf("GetMatchesOfJudge(db): %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var match models.Match

		err := rows.Scan(&match.Id, &match.Athlete1Id, &match.Athlete2Id, &match.JudgeId, &match.Date, &match.Score1, &match.Score2)

		if err != nil {
			return nil, fmt.Errorf("GetMatchesOfJudge(db): %v", err)
		}

		matches = append(matches, match)
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetMatchesOfJudge(db): %v", err)
	}

	return matches, nil
}

func DeleteMatch(id string) error {

	res, err := Instance.Exec("delete from matches where id=?", id)
	if err != nil {
		return fmt.Errorf("DeleteMatch(db): %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteMatch(db): %v", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
