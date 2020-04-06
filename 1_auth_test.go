package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"invink/account-service/forms"
	"invink/account-service/models"
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
	performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
}

func TestInitiateForAuthentication(t *testing.T) {
	DBNAMEORIGIN = os.Getenv("ACCOUNT_DB_DBNAME")
	os.Setenv("ACCOUNT_DB_DBNAME", "testing_db")
	ROUTER = setupServer()
	PUBLICKEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhTGv0frCyyhs3Xs5LyHE
4NXcM5lMqGJGNqCBo6zzjgv5BtZE5/bUHmJ8moUwTLLehtQt+wLq51wyJLe36142
3QNGO+5TCrKNWrOAxKhTRLwlHSjiXC/RgxbFYeD0EXGi54AwQRs27VFgzPRP7q4O
MtrXIinzqhhtJTorpP8t4n9FVXrpDmJnTbF5ct/3L+hCyeWmgAsrML3rHqJ+zfw1
DGogIrljdcLPzdlIcH9QjQJaWnfL7usl546aU0gkKjlUcB5+HUPNPkN3z9LEouHi
Kt8yVspTqyhnMnTNQnmGG7TuVCnWPXWaBaI/Aozgilj3+BIo9SiUIqKfc0FPeV61
LQIDAQAB
-----END PUBLIC KEY-----`

	createUser(ExampleEmail, ExampleUsername, ExamplePassword, "", "")
}

func TestProperEmailAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       ExampleEmail,
		Password: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth",
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
	w := performRequest(ROUTER, "POST", "/auth",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.NotEqual(t, "", response["token"]) // check if token is empty
}

func TestWrongUsernameAuthRequest(t *testing.T) {
	form := &forms.Authentication{
		ID:       "wrong_user",
		Password: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

func TestProperUsernameWrongPasswordAuthRequest(t *testing.T) {
	form := &forms.Authentication{
		ID:       "nothing",
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

func TestProperEmailWrongPasswordAuthRequest(t *testing.T) {
	form := &forms.Authentication{
		ID:       ExampleEmail,
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/auth",
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
	w := performRequest(ROUTER, "POST", "/auth",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

func TestCleanupForAuthentication(t *testing.T) {
	db := models.Setup()
	db.DropTable(&models.User{})
	db.DropTable("followed_by")
	db.DropTable("following")
	os.Setenv("ACCOUNT_DB_DBNAME", DBNAMEORIGIN)
	os.Setenv("ACCOUNT_DB_DBNAME", "testing_db")
}
