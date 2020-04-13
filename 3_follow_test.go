package main

import (
	"fmt"
	"testing"
)

const User1Username = "user1"
const User1Email = "user1@example.com"
const User2Username = "user2"
const User2Email = "user2@example.com"

func TestInitiateForFollow(t *testing.T) {
	setupDB()
	createUser(User1Email, User1Email, ExamplePassword, "", "")
	createUser(User2Email, User2Email, ExamplePassword, "", "")
	AUTHHEADER = []map[string]string{
		{"Authorization": fmt.Sprintf("Bearer %s", getToken(User1Email, ExamplePassword))},
		{"Authorization": fmt.Sprintf("Bearer %s", getToken(User2Email, ExamplePassword))},
	}
}

func TestProperFollowTestRequest(t *testing.T) {
	// follow specific user
	// check whether it's in the following list from the requested user's
	// check whether it's in the follower list from the followed user's
}

func TestFollowAlreadyFollowingUserRequest(t *testing.T) {
	// follow specific user
	// check 400 already following user
}

func TestUnfollowUserRequest(t *testing.T) {
	// unfollow specific user
	// check whether it's
}

func TestUnfollowNotFollowingUserRequest(t *testing.T) {
	// follow specific user
	// check 400 not a following user
}

func TestBlockUserRequest(t *testing.T) {
	// block specific user
	// check whether the blocked user still exists on the follower list
	// check whether it's accessible for the blocked user (should return 404 if the user is blocked)
	// check whether it's able to follow for the blocked user (should return 404)
	// check wheter it's accessible for the requested user (should return 400 blocked user)
	// check blocked user list
}

func TestBlockAlreadyBlockedUserRequest(t *testing.T) {
	// block specific user
	// check whether the blocked user still exists on the follower list
	// check whether it's accessible for the blocked user (should return 404 if the user is blocked)
	// check whether it's able to follow for the blocked user (should return 404)
	// check wheter it's accessible for the requested user (should return 400 blocked user)
	// check blocked user list
}

func TestUnBlockUserRequest(t *testing.T) {
	// block specific user
	// check whether it's accessible for the blocked user (should return 404 if the user is blocked)
	// check whether it's able to follow
}

func TestUnBlockNotBlockedUserRequest(t *testing.T) {
	// block specific user
	// check whether it's accessible for the blocked user (should return 404 if the user is blocked)
	// check whether it's able to follow
}

func TestCleanupForFollow(t *testing.T) {
	restoreEnvironment()
}
