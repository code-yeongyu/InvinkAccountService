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

	"invink/account-service/errors"
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
	if w.Code == 200 {
		json.Unmarshal([]byte(w.Body.String()), &response)
		return response["token"]
	}
	return ""
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

// test getting profile

func TestProperUsernameProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Profile{
		Username:        "changer",
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code)
	// change username to changer
	AUTHHEADER[0] = map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", getToken("changer", ExamplePassword)),
	}
	// update the token since the username has changed

	w = performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "changer", response["username"].(string))
	// check whether the username has changed
} // here changes test1's username to changer
func TestDuplicateUsernameProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Profile{
		Username:        "test2",
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, errors.UsernameExistsCode, int(response["error"].(float64)))
}
func TestImProperUsernameProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Profile{
		Username:        "test user",
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, errors.UsernameFormatErrorCode, int(response["error"].(float64)))
}
func TestProperNicknameProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Profile{
		Nickname:        "the game changer",
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code)
	// change username to changer

	w = performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "the game changer", response["nickname"].(string))
}
func TestProperBioProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Profile{
		Bio:             "Think different.",
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code)
	// change username to changer

	w = performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Think different.", response["bio"].(string))
}
func TestProperNicknameEmailProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := map[string]interface{}{
		"nickname":         "JOBS",
		"email":            "email@example.com",
		"current_password": ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code)
	// change username to changer

	w = performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "JOBS", response["nickname"].(string))
	assert.Equal(t, ExampleEmail, response["email"].(string))
}
func TestProperNicknameBioProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := map[string]interface{}{
		"nickname":         "thegreatmengmota",
		"bio":              "This is bio.",
		"current_password": ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code)
	// change username to changer

	w = performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/",
		AUTHHEADER[0],
		nil,
	)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, "thegreatmengmota", response["nickname"].(string))
	assert.Equal(t, "This is bio.", response["bio"].(string))
}

func TestProperPasswordProfilePatchRequest(t *testing.T) {
	form := &forms.Profile{
		Password:        "changed" + ExamplePassword,
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[1],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code)
	// change password

	assert.NotNil(t, getToken("test2", "changed"+ExamplePassword))
	// authenticate with the changed password
} // here changes test2's password to "changed" + ExamplePassword
func TestTooShortPasswordProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Profile{
		Password:        "aA-0",
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, errors.PasswordTooShortCode, int(response["error"].(float64)))
}
func TestVulerablePasswordProfilePatchRequest(t *testing.T) {
	var response map[string]interface{}
	form := &forms.Profile{
		Password:        "12345678",
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"PATCH",
		"/profile/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, errors.PasswordVulnerableErrorCode, int(response["error"].(float64)))
}

/*
func TestProperPictureProfilePatchRequest(t *testing.T) {
	// try to change picture should, success
}

// test updating profile


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

func TestImProperProfileDeleteRequest(t *testing.T) {
	form := &forms.Profile{
		CurrentPassword: ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"DELETE",
		"/profile/",
		AUTHHEADER[1],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProperProfileDeleteRequest(t *testing.T) {
	form := &forms.Profile{
		CurrentPassword: "changed" + ExamplePassword,
	}
	formJSON, _ := json.Marshal(form)
	w := performRequestWithHeader(
		ROUTER,
		"DELETE",
		"/profile/",
		AUTHHEADER[1],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequestWithHeader(
		ROUTER,
		"GET",
		"/profile/test2/",
		AUTHHEADER[0],
		strings.NewReader(string(formJSON)),
	)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// test delete profile

func TestCleanupForProfile(t *testing.T) {
	db := models.Setup()
	cleanUp(db)
}
