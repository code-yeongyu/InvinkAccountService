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
	"github.com/stretchr/testify/assert"

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

	form := &forms.Registration{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "A-maz1ng*pass",
		PublicKey: PUBLICKEY,
		Nickname:  "AmazingMengmota",
		Bio:       "Hi, I'm the great Mengmota",
	}
	formJSON, _ := json.Marshal(form)
	performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
}

func TestProperEmailAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       "test@example.com",
		Password: "A-maz1ng*pass",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/register",
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
		ID:       "testuser",
		Password: "A-maz1ng*pass",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.NotEqual(t, "", response["token"]) // check if token is empty
}

func TestWrongUsernameAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       "wrong_user",
		Password: "A-maz1ng*pass",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "", response["token"]) // check if token empty
}

func TestProperUsernameWrongPasswordAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       "nothing",
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "", response["token"]) // check if token empty
}

func TestProperEmailWrongPasswordAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       "test@example.com",
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "", response["token"]) // check if token empty
}

func TestWrongInfoAuthRequest(t *testing.T) {
	var response map[string]string
	form := &forms.Authentication{
		ID:       "nothing",
		Password: "12345678",
	}
	formJSON, _ := json.Marshal(form)
	w := performRequest(ROUTER, "POST", "/register",
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code) // check http status code
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "", response["token"]) // check if token empty
}

func TestCleanup(t *testing.T) {
	db := models.Setup()
	db.DropTable(&models.User{})
	db.DropTable("followed_by")
	db.DropTable("following")
	os.Setenv("ACCOUNT_DB_DBNAME", DBNAMEORIGIN)
	os.Setenv("ACCOUNT_DB_DBNAME", "testing_db")
}
