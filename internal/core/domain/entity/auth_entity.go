package entity

type LoginRequest struct {
	Email    string
	Password string
}

type AccessToken struct {
	AccessToken string
	ExpiredAt int64
}