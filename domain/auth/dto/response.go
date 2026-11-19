package dto

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"` // JWT yang akan digunakan client
	TokenType   string `json:"token_type"`   // jenis token, yaitu Bearer
	ExpiresIn   int64  `json:"expires_in"`   // masa berlaku token dalam detik
}
