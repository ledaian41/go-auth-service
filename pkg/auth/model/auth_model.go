package auth_model

type RegisterAccount struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `form:"email" json:"email" binding:"email"`
}
