package forms

type CreateUserRequest struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DetailUserResponse struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
