package errors

const (
	EmailEmpty       = 100
	EmailExists      = 101
	EmailFormatError = 102

	UsernameEmpty       = 200
	UsernameExists      = 201
	UsernameFormatError = 202

	PasswordEmpty           = 300
	PasswordTooShort        = 301
	PasswordVulnerableError = 302

	PublicKeyEmpty = 400
	PublicKeyError = 402
)
