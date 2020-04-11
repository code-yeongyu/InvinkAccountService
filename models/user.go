package models

import (
	"invink/account-service/errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User is a model for each user
type User struct {
	ID         uint64 `gorm:"primary_key"`
	Username   string `gorm:"column:username;unique;not null"`
	Email      string `gorm:"column:email;unique_index;not null"`
	Password   string `gorm:"column:password;not null"`
	Nickname   string `gorm:"column:nickname"`
	Bio        string `gorm:"column:bio"`
	PictureURL string `gorm:"column:picture_url"`
	Following  []User `gorm:"many2many:following;foreignkey:username;default:[]"`
	Follower   []User `gorm:"many2many:follower;foreignkey:username;default:[]"`
	PublicKey  string `gorm:"column:public_key;not null"`
	MyKeys     string `sql:"json" gorm:"column:my_keys;default:'{}'"`
}

func (u *User) isProperUsername(s string) bool {
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
		case unicode.IsLower(c):
		case unicode.IsUpper(c):
		case c == '.' || c == '-' || c == '_':
		default:
			return false
		}
	}
	return true
}
func (u *User) isSecurePassword(s string) bool {
	var number, lower, upper, special bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}
	return number && lower && upper && special
}

func NewUser() *User {
	return &User{}
}
func (u *User) GenerateHashedPassword(password string) (hashedPassword []byte) {
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), 12)
	return
}
func (u *User) ValidateUsername(db *gorm.DB, username string) int {
	if err := db.Where("username = ?", username).First(&User{}).Error; err == nil {
		return errors.UsernameExistsCode
	} // validating username duplicates
	if !u.isProperUsername(username) {
		return errors.UsernameFormatErrorCode
	} // validating if username is in proper format
	return -1
}
func (u *User) ValidateEmail(db *gorm.DB, email string) int {
	if err := db.Where("email = ?", email).First(&User{}).Error; err == nil {
		return errors.EmailExistsCode
	} // validating email duplicates
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // email format
	if !re.MatchString(email) {
		return errors.EmailFormatErrorCode
	} // validating if email is in proper format
	return -1
}
func (u *User) ValidatePassword(password string) int {
	if len(password) < 8 {
		return errors.PasswordTooShortCode
	} // validating password length
	if !u.isSecurePassword(password) {
		return errors.PasswordVulnerableCode
	}
	return -1
}
func (u *User) IsPublicKeyInFormat(publicKey string) bool {
	return strings.Contains(publicKey, "PUBLIC KEY")
}
func (u *User) IsPasswordCorrect(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
func (u *User) SetUserByEmailOrID(db *gorm.DB, emailOrId string) error {
	return db.Where("email = ? OR username = ?", emailOrId, emailOrId).First(&u).Error
}
func (u *User) SetUserByID(db *gorm.DB, ID uint64) error {
	return db.Model(u).Where("ID = ?", ID).First(&u).Error
}
func (u *User) SetUserByUsername(db *gorm.DB, username string) error {
	return db.Model(u).Where("username = ?", username).First(&u).Error
}
func (u *User) ToMyProfileMap() (myProfileMap map[string]interface{}) {
	myProfileMap = structs.Map(u)

	myProfileMap["my_keys"] = myProfileMap["MyKeys"]
	delete(myProfileMap, "MyKeys")
	myProfileMap["picture_url"] = myProfileMap["PictureURL"]
	delete(myProfileMap, "PictureURL")
	myProfileMap["public_key"] = myProfileMap["PublicKey"]
	delete(myProfileMap, "PublicKey")
	// camelCase to snake_case

	delete(myProfileMap, "ID")
	delete(myProfileMap, "Password")
	delete(myProfileMap, "Follower")
	delete(myProfileMap, "Following")
	// remove secure informations

	for k, v := range myProfileMap {
		if unicode.IsUpper(rune(k[0])) {
			if v != "" {
				myProfileMap[strings.ToLower(k)] = v
			}
			delete(myProfileMap, k)
		} else {
			if v == "" {
				delete(myProfileMap, k)
			}
		}
	}
	// Upper-string key to Lower-string key

	followingUsername := make([]string, len(u.Following))
	followerUsername := make([]string, len(u.Follower))
	for i, v := range u.Following {
		followingUsername[i] = v.Username
	}
	for i, v := range u.Follower {
		followerUsername[i] = v.Username
	}
	myProfileMap["following_username"] = followingUsername
	myProfileMap["follower_username"] = followerUsername
	// Add following and follwed usernames instead of raw following/followed model
	return
}
func (u *User) ToPublicProfileMap() (publicProfileMap map[string]interface{}) {
	publicProfileMap = map[string]interface{}{
		"username":      u.Username,
		"following_cnt": len(u.Following),
		"follower_cnt":  len(u.Follower),
	}
	if u.Nickname != "" {
		publicProfileMap["nickname"] = u.Nickname
	}
	if u.Bio != "" {
		publicProfileMap["bio"] = u.Bio
	}
	if u.PictureURL != "" {
		publicProfileMap["picture_url"] = u.PictureURL
	}
	// Add nickname, bio, and picture_url only if they exists
	return
}
