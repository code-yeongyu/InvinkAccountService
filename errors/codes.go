package errors

const (
	FormErrorCode = 1

	EmailExistsCode      = 101
	EmailFormatErrorCode = 102

	UsernameExistsCode      = 201
	UsernameFormatErrorCode = 202

	PasswordTooShortCode        = 301
	PasswordVulnerableErrorCode = 302

	PublicKeyErrorCode = 402

	AuthenticationFailureCode    = 500
	EmptyAuthorizationHeaderCode = 501
	WrongTokenTypeCode           = 502
)
