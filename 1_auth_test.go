package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"invink/account-service/forms"
)

func createUser(email string, username string, password string, nickname string, bio string) {
	PUBLICKEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhTGv0frCyyhs3Xs5LyHE
4NXcM5lMqGJGNqCBo6zzjgv5BtZE5/bUHmJ8moUwTLLehtQt+wLq51wyJLe36142
3QNGO+5TCrKNWrOAxKhTRLwlHSjiXC/RgxbFYeD0EXGi54AwQRs27VFgzPRP7q4O
MtrXIinzqhhtJTorpP8t4n9FVXrpDmJnTbF5ct/3L+hCyeWmgAsrML3rHqJ+zfw1
DGogIrljdcLPzdlIcH9QjQJaWnfL7usl546aU0gkKjlUcB5+HUPNPkN3z9LEouHi
Kt8yVspTqyhnMnTNQnmGG7TuVCnWPXWaBaI/Aozgilj3+BIo9SiUIqKfc0FPeV61
LQIDAQAB
-----END PUBLIC KEY-----`

	form := &forms.Registration{
		Email:     email,
		Username:  username,
		Password:  password,
		PublicKey: PUBLICKEY,
	}
	if nickname != "" {
		form.Nickname = nickname
	}
	if bio != "" {
		form.Bio = bio
	}
	formJSON, _ := json.Marshal(form)
	performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
}

// test util

func TestInitiateForAuthentication(t *testing.T) {
	setupDB()
	ROUTER = setupServer()
	createUser(ExampleEmail, ExampleUsername, ExamplePassword, "", "")
}

func TestProperEmailAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       ExampleEmail,
		Password: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.NotEqual(t, "", response["token"]) // check if token is empty
}
func TestProperUsernameAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       ExampleUsername,
		Password: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.NotEqual(t, "", response["token"]) // check if token is empty
}

// test successful cases

func TestProperEmailWrongPasswordAuthRequest(t *testing.T) {
	form := &forms.Authentication{
		ID:       ExampleEmail,
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}
func TestWrongUsernameWrongPasswordAuthRequest(t *testing.T) {
	form := &forms.Authentication{
		ID:       "nothing",
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}
func TestWrongUsernameAuthRequest(t *testing.T) {
	form := &forms.Authentication{
		ID:       "wrong_user",
		Password: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}
func TestWrongInfoAuthRequest(t *testing.T) {
	form := &forms.Authentication{
		ID:       "nothing",
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

// test failure cases

func TestCleanUpForAuthentication(t *testing.T) {
	restoreEnvironment()
}
