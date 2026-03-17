package models

type Judge struct {
	Id              *string `json:"id"`
	Name            *string `json:"name"`
	Phone           *string `json:"phone"`
	ExperienceYears *int    `json:"experience_years"`
}

type JudgeWithPassword struct {
	Id              *string `json:"id"`
	Name            *string `json:"name"`
	Password        *string `json:"password"`
	Phone           *string `json:"phone"`
	ExperienceYears *int    `json:"experience_years"`
}
