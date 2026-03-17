package models

type Match struct {
	Id         *string `json:"id"`
	Athlete1Id *string `json:"athlete1_id"`
	Athlete2Id *string `json:"athlete2_id"`
	JudgeId    *string `json:"judge_id"`
	Date       *int64  `json:"date"`
	Score1     *int    `json:"score1"`
	Score2     *int    `json:"score2"`
}
