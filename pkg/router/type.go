package router

type UserResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type SignUpRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
