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
	TypeID   uint   `json:"type_id" example:"1"` // (1: Free, 2: Premium, 3: Pro)
}

// UpdateUser is a data transfer object for update user
type UpdateUser struct {
	Name             string `json:"name" example:"Dino Puguh"`
	RemainingTime    int64  `json:"remaining_time" example:"1800"`
	ReachedTimeLimit bool   `json:"reached_time_limit" example:"false"`
}
