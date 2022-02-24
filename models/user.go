package models

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Friends []User `json:"friends"`
}

type UserList struct {
	Users []User `json:"users"`
}
