package shared_dto

type UserDTO struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Name        string   `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	Role        []string `json:"role"`
	Site        string   `json:"site"`
}

type SiteDTO struct {
	ID        string `json:"id"`
	SecretKey string `json:"secret_key"`
}

type LoginRequestDTO struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterRequestDTO struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required,min=6"`
	Email       string `form:"email" json:"email" binding:"email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber"`
}
