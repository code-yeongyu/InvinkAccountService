package errors

// Messages contains the error message for each error code
var Messages = map[int]string{
	FormErrorCode:      "Form error.",
	ParameterErrorCode: "You're not allowed to use this parameter.",

	EmailExistsCode:      "Email exists.",
	EmailFormatErrorCode: "Email format error.",

	UsernameExistsCode:      "Username exists.",
	UsernameFormatErrorCode: "Username format error.",

	PasswordTooShortCode:        "Password is too short.",
	PasswordVulnerableErrorCode: "Password is vulnerable.",

	PublicKeyErrorCode: "Not a proper public key.",

	AuthenticationFailureCode:    "Authentication Failed.",
	EmptyAuthorizationHeaderCode: "No authorization header.",
	WrongTokenTypeCode:           "Only Bearer token is available.",
}
