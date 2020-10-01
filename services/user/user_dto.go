package user

// LoginUser is a data transfer object for user login
type LoginUser struct {
	Email    string `json:"email" example:"dinopuguh@mycap.com"`
	Password string `json:"password" example:"s3cr3tp45sw0rd"`
}

// RegisterUser is a data transfer object for create user
type RegisterUser struct {
	Name     string `json:"name" example:"Dino Puguh"`
	Username string `json:"username" example:"dinopuguh"`
	Email    string `json:"email" example:"dinopuguh@mycap.com"`
	Password string `json:"password" example:"s3cr3tp45sw0rd"`
}
