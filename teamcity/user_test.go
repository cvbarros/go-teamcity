package teamcity_test

import (
	"fmt"
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/require"
)

func TestUser_Create(t *testing.T) {
	newUser, _ := teamcity.NewUser("username", "First Middle Last", "test@test.test")
	client := setup()
	actual, err := client.Users.Create(newUser)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotEmpty(t, actual.Name)
	require.NotEmpty(t, actual.Username)

	cleanUpUser(t, client, actual.ID)

	require.NotEqual(t, newUser.ID, actual.ID)
	require.Equal(t, newUser.Name, actual.Name)
	require.Equal(t, newUser.Email, actual.Email)
}

func TestUser_Get(t *testing.T) {
	newUser, _ := teamcity.NewUser("usernameGet", "First Middle Last", "test@test.test")

	client := setup()
	actual, err := client.Users.Create(newUser)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotEmpty(t, actual.Name)
	require.NotEmpty(t, actual.Username)
	require.NotZero(t, actual.ID)

	userByUsername, err := client.Users.GetByUsername(newUser.Username)

	require.NoError(t, err)
	require.NotNil(t, userByUsername)
	require.NotEmpty(t, userByUsername.Name)
	require.NotEmpty(t, userByUsername.Username)
	require.NotZero(t, userByUsername.ID)

	userByName, err := client.Users.GetByName(newUser.Name)

	require.NoError(t, err)
	require.NotNil(t, userByName)
	require.NotEmpty(t, userByName.Name)
	require.NotEmpty(t, userByName.Username)
	require.NotZero(t, userByName.ID)

	userByID, err := client.Users.GetByID(actual.ID)

	require.NoError(t, err)
	require.NotNil(t, userByID)
	require.NotEmpty(t, userByID.Name)
	require.NotEmpty(t, userByID.Username)
	require.NotZero(t, userByID.ID)

	cleanUpUser(t, client, actual.ID)

	require.NotEqual(t, newUser.ID, actual.ID)
	require.Equal(t, newUser.Name, actual.Name)
	require.Equal(t, newUser.Email, actual.Email)

	require.Equal(t, actual, userByID)
	require.Equal(t, actual, userByUsername)
}

func TestUser_Delete(t *testing.T) {
	newUser, _ := teamcity.NewUser("usernameDel", "First Middle Last", "test@test.test")
	client := setup()
	actual, err := client.Users.Create(newUser)

	err = client.Users.DeleteByID(actual.ID)

	require.NoError(t, err)

	_, err = client.Users.GetByUsername(newUser.Username)

	// User is deleted, so expect error, and message to contain 404 (NOT FOUND)
	require.Error(t, err)
	require.Contains(t, err.Error(), "404")
}

func TestUser_List(t *testing.T) {
	users := []*teamcity.User{}
	client := setup()

	for i := 0; i < 5; i++ {
		user, err := teamcity.NewUser(
			fmt.Sprint("TESTUSERNAME", i),
			fmt.Sprint("Test User List ", i),
			fmt.Sprint("Test User Description List ", i),
		)
		require.NoError(t, err)
		users = append(users, user)
		_, err = client.Users.Create(user)
		require.NoError(t, err)
	}
	userList, err := client.Users.List()

	require.NoError(t, err)
	require.Equal(t, 5, userList.Count-1)

	for _, user := range userList.Items {
		if user.Username == "admin" {
			continue
		}
		actual, err := client.Users.GetByUsername(user.Username)
		require.NoError(t, err)

		cleanUpUser(t, client, actual.ID)

		_, err = client.Users.GetByName(user.Name)
		require.Error(t, err)
	}
	userList, err = client.Users.List()
	require.NoError(t, err)
	require.Equal(t, 1, userList.Count)
}

func TestUser_Group(t *testing.T) {
	client := setup()

	newGroup, _ := teamcity.NewGroup("TESTGROUPMEMBER", "Test Group Member", "")
	actualGroup, err := client.Groups.Create(newGroup)
	require.NoError(t, err)
	require.Zero(t, actualGroup.Users.Count)

	newUser, _ := teamcity.NewUser("testusermember", "Test User Member", "test@member.com")
	actualUser, err := client.Users.Create(newUser)
	require.NoError(t, err)

	actualGroup, err = client.Users.GroupAddByID(actualUser.ID, actualGroup.Key)
	require.NoError(t, err)
	require.NotNil(t, actualGroup)
	require.NotZero(t, actualGroup.Users.Count)

	actualGroup, err = client.Users.GroupDeleteByID(actualUser.ID, actualGroup.Key)

	cleanUpGroup(t, client, newGroup.Key)
	cleanUpUser(t, client, actualUser.ID)

	require.NoError(t, err)
	require.NotNil(t, actualGroup)
	require.Nil(t, actualGroup.Users)

}

func cleanUpUser(t *testing.T, client *teamcity.Client, id int) {
	client.Users.DeleteByID(id)
}
