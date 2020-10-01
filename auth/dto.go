package auth

// Login is a data transfer object for user's login
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
