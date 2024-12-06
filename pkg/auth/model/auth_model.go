package auth_model

type LoginAccount struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterAccount struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required,min=6"`
	Email       string `form:"email" json:"email" binding:"email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber"`
}
