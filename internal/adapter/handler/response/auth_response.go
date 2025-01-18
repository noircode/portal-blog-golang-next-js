package response

type SuccessAuthResponse struct {
	Meta
	AccessToken string `json:"access_token"`
	ExpiredAt int64 `json:"expired_at"`
}