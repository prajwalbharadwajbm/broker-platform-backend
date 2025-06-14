package dtos

// this struct is used to fetch user data from request body
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
