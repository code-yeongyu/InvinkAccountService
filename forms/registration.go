package forms

// Registration is a form for registration
type Registration struct {
	Email       string `json:"email" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PublicKey   string `json:"public_key" binding:"required"`
	ContentsKey string `json:"contents_Key" binding:"required"`
	Nickname    string `json:"nickname"`
	Bio         string `json:"bio"`
}
