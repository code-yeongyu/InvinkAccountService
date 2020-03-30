package models

// UserModel is a model for each user
type UserModel struct {
	ID         uint        `gorm:"primary_key"`
	Username   string      `gorm:"column:username;unique;not null"`
	Email      string      `gorm:"column:email;unique_index;not null"`
	Password   string      `gorm:"column:password;not null"`
	Nickname   string      `gorm:"column:nickname"`
	Bio        string      `gorm:"column:bio"`
	PictureURL string      `gorm:"column:picture_url"`
	Following  []UserModel `gorm:"column:following;default:[]"`
	FollowedBy []UserModel `gorm:"column:followed_by;default:[]"`
	PublicKey  string      `gorm:"column:public_key;not null"`
	MyKeys     string      `sql:"json" gorm:"column:my_keys;default:'{}'"`
}
