package models

type Athlete struct {
	Id      *string `json:"id"`
	Name    *string `json:"name"`
	Country *string `json:"country"`
	Age     *int    `json:"age"`
}
