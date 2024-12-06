package site_model

type Site struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SecretKey string `json:"secret_key"`
}
