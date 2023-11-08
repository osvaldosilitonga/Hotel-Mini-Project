package dto

// --------- REQUEST BODY -----------
type RegisterBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// --------- RESPONSE BODY -----------
type RegisterResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type LoginResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}
