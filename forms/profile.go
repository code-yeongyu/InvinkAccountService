package forms

// Profile is a form for Profile
type Profile struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	Nickname        string `json:"nickname"`
	PictureURL      string `json:"picture_url"`
	Bio             string `json:"bio"`
	MyKeys          string `json:"my_keys"`
	CurrentPassword string `json:"current_password" binding:"required"`
}
