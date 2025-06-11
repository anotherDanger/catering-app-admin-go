package web

type AdminResponse struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}
