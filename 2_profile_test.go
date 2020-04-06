package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"invink/account-service/forms"
	"invink/account-service/models"
)

var AUTHHEADER []map[string]string

func performRequestWithHeader(r http.Handler, method string, path string, header map[string]string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	for key, value := range header {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func getToken(ID string, password string) string {
	authForm := &forms.Authentication{
		ID:       ID,
		Password: password,
	}
	var response map[string]string
	formJSON, _ := json.Marshal(authForm)
	w := performRequest(ROUTER, "POST", "/auth",
		strings.NewReader(string(formJSON)),
	)
	json.Unmarshal([]byte(w.Body.String()), &response)
	return response["token"]
}

func TestInitiateForProfile(t *testing.T) {
	DBNAMEORIGIN = os.Getenv("ACCOUNT_DB_DBNAME")
	os.Setenv("ACCOUNT_DB_DBNAME", "testing_db")
	ROUTER = setupServer()
	createUser(ExampleEmail, "test1", ExamplePassword, ExampleNickname, "")
	createUser("test2@example.com", "test2", ExamplePassword, "", ExampleBio)
	/* register */

	AUTHHEADER = append(AUTHHEADER, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", getToken("test1", ExamplePassword)),
	})
	// get token for test1
	AUTHHEADER = append(AUTHHEADER, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", getToken("test2", ExamplePassword)),
	})
	// get token for test2

	/* get token */
}

// init test

func TestMyEmptyBioProfileRequest(t *testing.T) {
	var response map[string]interface{}
	w := performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/test1/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, "test1", response["username"].(string))
	assert.Equal(t, ExampleEmail, response["email"].(string))
	assert.Equal(t, ExampleNickname, response["nickname"].(string))
	assert.Nil(t, response["bio"])
	assert.Nil(t, response["picture_url"])
	assert.Equal(t, []interface{}{}, response["following_username"])
	assert.Equal(t, []interface{}{}, response["follower_username"])
	assert.NotNil(t, response["public_key"])
	assert.Equal(t, "{}", response["my_keys"])
}
func TestMyEmptyNicknameProfileRequest(t *testing.T) {
	var response map[string]interface{}
	w := performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/",
		AUTHHEADER[1],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, "test2", response["username"].(string))
	assert.Equal(t, "test2@example.com", response["email"].(string))
	assert.Nil(t, response["nickname"])
	assert.Equal(t, ExampleBio, response["bio"].(string))
	assert.Nil(t, response["picture_url"])
	assert.Equal(t, []interface{}{}, response["following_username"])
	assert.Equal(t, []interface{}{}, response["follower_username"])
	assert.NotNil(t, response["public_key"])
	assert.Equal(t, "{}", response["my_keys"])
}
func TestOtherUserRequest(t *testing.T) {
	var response map[string]interface{}
	w := performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/test2/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, "test2", response["username"].(string))
	assert.Nil(t, response["nickname"])
	assert.Equal(t, ExampleBio, response["bio"].(string))
	assert.Nil(t, response["picture_url"])
	assert.Equal(t, 0, int(response["following_cnt"].(float64)))
	assert.Equal(t, 0, int(response["follower_cnt"].(float64)))
	// should not be empty

	assert.Nil(t, response["email"])
	assert.Nil(t, response["my_keys"])
	// should nil
}
func Test404ProfileRequest(t *testing.T) {
	w := performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/nothing/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusNotFound, w.Code)
	// should nil
}
func TestNoAuthorizationRequest(t *testing.T) {
	w := performRequest(
		ROUTER,
		"GET",
		"/profile/test1/",
		nil,
	)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// test get profile

func TestMyProfileUsernamePatchRequest(t *testing.T) {
	// try to change username, should fail
}
func TestMyProfileEmailPatchRequest(t *testing.T) {
	// try to change email, should fail
}
func TestMyProfilePasswordPatchRequest(t *testing.T) {
	// try to change pass, should success
}
func TestMyProfileNicknamePatchRequest(t *testing.T) {
	// try to change nickname, should success
}
func TestMyProfileBioPatchRequest(t *testing.T) {
	// try to change bio should, success
}
func TestMyProfilePicturePatchRequest(t *testing.T) {
	// try to change picture should, success
}
func TestMyProfileNicknameEmailPatchRequest(t *testing.T) {
	// try to change nickname and Email, only nickname should be changed
}
func TestMyProfileNicknameBioPatchRequest(t *testing.T) {
	// try to change nickname and bio, both should be changed
}
func TestOtherUserNicknameBioPatchRequest(t *testing.T) {
	// try to change nickname and bio, both should not be changed with 403
}

/*
func TestDeleteNicknameRequest(t *testing.T) {
	w := performRequestWithHeader(
		ROUTER,
		"DELETE",
		"/profile/test1/nickname",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestDeleteBioRequest(t *testing.T) {
	w := performRequestWithHeader(
		ROUTER,
		"DELETE",
		"/profile/test1/bio",
		AUTHHEADER[1],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestEmptyNicknameCheck(t *testing.T) {
	var response map[string]interface{}
	w := performRequestWithHeader(
		ROUTER,
		"POST",
		"/profile/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, response["nickname"])
}
func TestEmptyBioCheck(t *testing.T) {
	var response map[string]interface{}
	w := performRequestWithHeader(
		ROUTER,
		"POST",
		"/profile/",
		AUTHHEADER[1],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, response["bio"])
}
*/
// uncomment the upper codes when they're required

// test update profile

func TestOtherProfileDeleteRequest(t *testing.T) {
	// delete other profile, should fail
}
func TestMyProfileDeleteRequest(t *testing.T) {
	// delete my profile, should success
}

// test delete profile

func TestCleanupForProfile(t *testing.T) {
	db := models.Setup()
	cleanUp(db)
}
