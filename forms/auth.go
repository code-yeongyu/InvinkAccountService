package forms

// Authentication is a form for authentication
type Authentication struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}
