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

func integerExists(array []int, target int) bool {
	for i := range array {
		if array[i] == target {
			return true
		}
	}
	return false
}

var DB *gorm.DB
var ROUTER *gin.Engine
var PUBLICKEY string
var DBNAMEORIGIN string

func TestInitiate(t *testing.T) {
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
}

func TestProperReuqest(t *testing.T) {
	var userModel models.UserModel
	form := &forms.Registration{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code) // check http status code
	err := DB.Where("username = ?", form.Username).First(userModel).Error
	assert.NotNil(t, err) // check record existing
}

func TestEmailDuplicate(t *testing.T) {
	var response map[string][]int
	form := &forms.Registration{
		Email:     "test@example.com",
		Username:  "testuser1",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                           // check response form error
	assert.Equal(t, response["error"][0], errors.EmailEmptyCode) // check error
}

func TestUsernameDuplicate(t *testing.T) {
	var response map[string][]int
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "testuser",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                              // check response form error
	assert.Equal(t, response["error"][0], errors.UsernameEmptyCode) // check error
}

func TestEmailUsernameDuplicate(t *testing.T) {
	var response map[string][]int
	form := &forms.Registration{
		Email:     "test@example.com",
		Username:  "testuser1",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                          // check response form error
	assert.True(t, integerExists(response["error"], errors.EmailExistsCode))    // check email error
	assert.True(t, integerExists(response["error"], errors.UsernameExistsCode)) // check username error
}

func TestImproperEmail(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test@example",
		Username:  "test1",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                 // check response form error
	assert.Equal(t, response["error"][0], errors.EmailFormatErrorCode) // check error
}

func TestImproperUsername(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test/user",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                    // check response form error
	assert.Equal(t, response["error"][0], errors.UsernameFormatErrorCode) // check error
}

func TestTooShortPassword(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test/user",
		Password:  "a",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                 // check response form error
	assert.Equal(t, response["error"][0], errors.PasswordTooShortCode) // check error
}

func TestVulnerablePassword(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		Password:  "12345678",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                                        // check response form error
	assert.Equal(t, response["error"][0], errors.PasswordVulnerableErrorCode) // check error
}

func TestImproperPublicKey(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		Password:  "12345678",
		PublicKey: "",
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                               // check response form error
	assert.Equal(t, response["error"][0], errors.PublicKeyErrorCode) // check error
}

func TestEmailEmpty(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Username:  "test",
		Password:  "12345678",
		PublicKey: "",
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                           // check response form error
	assert.Equal(t, response["error"][0], errors.EmailEmptyCode) // check error
}

func TestUsernameEmpty(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test1@example.com",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                              // check response form error
	assert.Equal(t, response["error"][0], errors.UsernameEmptyCode) // check error
}

func TestPasswordEmpty(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                              // check response form error
	assert.Equal(t, response["error"][0], errors.PasswordEmptyCode) // check error
}

func TestPublicKeyEmpty(t *testing.T) {
	var response map[string][]string
	form := &forms.Registration{
		Email:     "test1@example.com",
		Username:  "test",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)

	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)                                               // check response form error
	assert.Equal(t, response["error"][0], errors.PublicKeyEmptyCode) // check error
}

func TestCleanup(t *testing.T) {
	os.Setenv("ACCOUNT_DB_DBNAME", DBNAMEORIGIN)
	os.Setenv("ACCOUNT_DB_DBNAME", "testing_db")
	DB.Where("1=1").Delete(models.UserModel{})
}
