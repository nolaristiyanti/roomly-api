package entity

type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
}
