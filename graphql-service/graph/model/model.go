package model

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	JobTitle  string `json:"jobTitle,omitempty"`
	Feedback  string `json:"feedback"`
	CreateAt  string `json:"createAt"`
}
