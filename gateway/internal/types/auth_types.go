package types

// CreateUserRequest representa el cuerpo JSON para crear un usuario
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest representa el cuerpo JSON para el login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
