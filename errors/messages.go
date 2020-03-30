package errors

const Messages = map[int]String{
	EmailEmptyCode:       "The email field is empty.",
	EmailExistsCode:      "Email exists.",
	EmailFormatErrorCode: "Email format error.",

	UsernameEmptyCode:       "Username field is empty.",
	UsernameExistsCode:      "Username exists.",
	UsernameFormatErrorCode: "Username format error.",

	PasswordEmptyCode:           "Password field is empty.",
	PasswordTooShortCode:        "Password is too short.",
	PasswordVulnerableErrorCode: "Password is vulnerable.",

	PublicKeyEmptyCode: "PublicKey field is empty.",
	PublicKeyErrorCode: "Not a proper public key.",
}
