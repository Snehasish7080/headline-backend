package models

type Opinion struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	User        User   `json:"user"`
}
