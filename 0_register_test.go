package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"invink/account-service/errors"
	"invink/account-service/forms"
	"invink/account-service/models"
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var ROUTER *gin.Engine

var DBNAMEORIGIN string

const ExampleEmail = "test@example.com"
const ExampleUsername = "testuser"
const ExamplePassword = "A-maz1ng*pass"
const ExampleNickname = "AmazingMengmota"
const ExampleBio = "Hi, I'm the great Mengmota"
const PUBLICKEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhTGv0frCyyhs3Xs5LyHE
4NXcM5lMqGJGNqCBo6zzjgv5BtZE5/bUHmJ8moUwTLLehtQt+wLq51wyJLe36142
3QNGO+5TCrKNWrOAxKhTRLwlHSjiXC/RgxbFYeD0EXGi54AwQRs27VFgzPRP7q4O
MtrXIinzqhhtJTorpP8t4n9FVXrpDmJnTbF5ct/3L+hCyeWmgAsrML3rHqJ+zfw1
DGogIrljdcLPzdlIcH9QjQJaWnfL7usl546aU0gkKjlUcB5+HUPNPkN3z9LEouHi
Kt8yVspTqyhnMnTNQnmGG7TuVCnWPXWaBaI/Aozgilj3+BIo9SiUIqKfc0FPeV61
LQIDAQAB
-----END PUBLIC KEY-----`

func setupDB() (db *gorm.DB) {
	DBNAMEORIGIN = os.Getenv("ACCOUNT_DB_DBNAME")
	os.Setenv("ACCOUNT_DB_DBNAME", "testing_db")
	db = models.Setup()
	db.DropTable(&models.User{})
	db.DropTable("follower")
	db.DropTable("following")
	return
}

func restoreEnvironment() {
	os.Setenv("ACCOUNT_DB_DBNAME", DBNAMEORIGIN)
}

// test util

func TestInitiateForRegistration(t *testing.T) {
	setupDB()
	ROUTER = setupServer()
}

func TestProperRegistrationRequest(t *testing.T) {
	form := &forms.Registration{
		Email:     ExampleEmail,
		Username:  ExampleUsername,
		Password:  ExamplePassword,
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusCreated, w.Code) // check http status code
}

// test proper request
func TestEmailDuplicateRegistrationRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Registration{
		Email:     ExampleEmail,
		Username:  "testuser1",
		Password:  ExamplePassword,
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                        // check response form error
	assert.Equal(t, errors.EmailExistsCode, int(response["error"].(float64))) // check error
}
func TestImproperEmailRegistrationRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Registration{
		Email:     "test@example",
		Username:  "test1",
		Password:  ExamplePassword,
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                             // check response form error
	assert.Equal(t, errors.EmailFormatErrorCode, int(response["error"].(float64))) // check error
}
func TestEmailEmptyRegistrationRequest(t *testing.T) {
	form := &forms.Registration{
		Username:  "test",
		Password:  "12345678",
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

// test email error
func TestUsernameDuplicateRegistrationRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  ExampleUsername,
		Password:  ExamplePassword,
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                           // check response form error
	assert.Equal(t, errors.UsernameExistsCode, int(response["error"].(float64))) // check error
}
func TestImproperUsernameRegistrationRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test/user",
		Password:  ExamplePassword,
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                                // check response form error
	assert.Equal(t, errors.UsernameFormatErrorCode, int(response["error"].(float64))) // check error
}
func TestUsernameEmptyRegistrationRequest(t *testing.T) {
	form := &forms.Registration{
		Email:     "test1@example.com",
		Password:  ExamplePassword,
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

// test username error
func TestTooShortPasswordRegistrationRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		Password:  "a",
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                             // check response form error
	assert.Equal(t, errors.PasswordTooShortCode, int(response["error"].(float64))) // check error
}
func TestVulnerablePasswordRegistrationRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		Password:  "12345678",
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                               // check response form error
	assert.Equal(t, errors.PasswordVulnerableCode, int(response["error"].(float64))) // check error
}
func TestPasswordEmptyRegistrationRequest(t *testing.T) {
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		PublicKey: PUBLICKEY,
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

// test password error
func TestImproperPublicKeyRegistrationRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		Password:  ExamplePassword,
		PublicKey: "error key",
		Nickname:  ExampleNickname,
		Bio:       ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                           // check response form error
	assert.Equal(t, errors.PublicKeyErrorCode, int(response["error"].(float64))) // check error
}
func TestPublicKeyEmptyRegistrationRequest(t *testing.T) {
	form := &forms.Registration{
		Email:    "test1@example.com",
		Username: "test",
		Password: ExamplePassword,
		Nickname: ExampleNickname,
		Bio:      ExampleBio,
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register/",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
}

// test public key error
func TestCleanUpForRegistration(t *testing.T) {
	restoreEnvironment()
}
