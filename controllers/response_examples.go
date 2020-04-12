package controllers

// EmptyResponse is for the empty response type
type EmptyResponse struct {
}

// TypicalErrorResponse Example
type TypicalErrorResponse struct {
	Code    int
	Message string `json:"msg"`
	Detail  string
}

// AuthenticatedResponse Example
type AuthenticatedResponse struct {
	Token string
}

// PublicProfileResponse Example
type PublicProfileResponse struct {
	Username     string `example:"testuser"`
	FollowingCnt int    `json:"following_cnt" example:"0"`
	FollowerCnt  int    `json:"follower_cnt" example:"0"`
	Nickname     string `example:"TheGreatMengmota"`
	Bio          string `example:"Hi, I'm the great mengmota!"`
	PictureURL   string `example:"some kind of aws s3 url"`
}

// MyProfileResponse Example
type MyProfileResponse struct {
	Bio               string   `example:"Think Different"`
	Nickname          string   `example:"Steve Jobs"`
	PictureURL        string   `json:"picture_url" example:"some kind of aws s3 url"`
	FollowerUsername  []string `json:"follower_username" example:"[\"mark_zuckerberg\", \"the_great_edison\", \"tim_cook\"]"`
	FollowingUsername []string `json:"following_username" example:"[\"mark_zuckerberg\", \"bob_dylan\", \"michael_jackson\"]"`
	MyKeys            string   `json:"my_keys" example:"{\"encrypted_master_key\": <a aes-256 master key encrypted with user's password>,\"encrypted_private_key\": <a rsa-2048 private key encrypted with master key>,\"encrypted_contents_key\": <a aes-256 contents key encrypted with master key>,\"encrypted_following_key\": [{\"mark_zuckerberg\":<a aes-256 mark_zuckerberg's contents key encrypted with master key>,\"bob_dylan\":<a aes-256 bob_dylan's contents key encrypted with master key>,\"michael_jackson\":<a aes-256 michael's contents key encrypted with master key>}]}"`
	PublicKey         string   `json:"public_key" example:"-----BEGIN PUBLIC KEY-----MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhTGv0frCyyhs3Xs5LyHE4NXcM5lMqGJGNqCBo6zzjgv5BtZE5/bUHmJ8moUwTLLehtQt+wLq51wyJLe361423QNGO+5TCrKNWrOAxKhTRLwlHSjiXC/RgxbFYeD0EXGi54AwQRs27VFgzPRP7q4OMtrXIinzqhhtJTorpP8t4n9FVXrpDmJnTbF5ct/3L+hCyeWmgAsrML3rHqJ+zfw1DGogIrljdcLPzdlIcH9QjQJaWnfL7usl546aU0gkKjlUcB5+HUPNPkN3z9LEouHiKt8yVspTqyhnMnTNQnmGG7TuVCnWPXWaBaI/Aozgilj3+BIo9SiUIqKfc0FPeV61LQIDAQAB-----END PUBLIC KEY-----"`
	Username          string   `json:"username" example:"jobs_the_future"`
}
