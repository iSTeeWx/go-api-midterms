package models

type Judge struct {
	Id              *string `json:"id"`
	Name            *string `json:"name"`
	Password        *string `json:"password,omitempty"`
	Phone           *string `json:"phone"`
	ExperienceYears *int    `json:"experience_years"`
}
