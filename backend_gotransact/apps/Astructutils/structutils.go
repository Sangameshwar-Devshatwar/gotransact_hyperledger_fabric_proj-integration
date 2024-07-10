package astructutils

type SignupUser struct {
	FirstName   string `json:"firstname" binding:"required,alpha"`
	LastName    string `json:"lastname" binding:"required,alpha"`
	Email       string `json:"email" binding:"required,email"`
	CompanyName string `json:"companyname" binding:"required"`
	Password    string `json:"password" binding:"required" validate:"password_complexity"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required" validate:"password_complexity"`
}
