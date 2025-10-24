package model

// UserRequest represents a user submission
type UserRequest struct {
	UserID  int  `json:"userId"`
	Correct bool `json:"correct"`
}
