package controllers

// EmptyResponse is for the empty response type
type EmptyResponse struct {
}

// EmailExistsResponse Example
type EmailExistsResponse struct {
	Code    int
	Message string `json:"msg"`
}

// AuthenticatedResponse Example
type AuthenticatedResponse struct {
	Token string
}
