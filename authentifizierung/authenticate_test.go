/*
@author: 2447899 8689159 3000685
*/
package authentifizierung

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// tests pon create user
func TestCreateUser(t *testing.T) {
	user := "admin"
	password := "admin"
	assert.Equal(t, 0, len(users.users))
	_ = CreateUser(&user, &password)
	assert.Equal(t, 1, len(users.users))
}

// tests that on authentication it returns a true user
func TestAuthenticateUser_True(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	wahr, _ := AuthenticateUser(&user, &password)
	assert.Equal(t, true, wahr)
}

// checks that if it fails that it returns a wrong user
func TestAuthenticateUser_False(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	passwordWrong := "user"
	wahr, cookie := AuthenticateUser(&user, &passwordWrong)
	assert.Equal(t, false, wahr)
	assert.NotEqual(t, nil, cookie)
}

// test that the cookie which is given back is the right one
func TestAuthenticateUser_CookieMngmt(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	// nutzer erstellen
	_ = CreateUser(&user, &password)
	// den coookie zurückholen aufgebaut wie:
	_, cookie := AuthenticateUser(&user, &password)
	// auslesen des cookies und des uernames
	username := cookie[:strings.Index(cookie, "|")]
	cookie = cookie[strings.Index(cookie, "|")+1:]
	// schauen, dass der neue hash richtig generiert ist
	assert.Equal(t, user, username)
	assert.Equal(t, 100, len(cookie))
	assert.NotEmpty(t, cookies.cookies[cookie])
	// check that the end time is roughly in 15 minutes right before it
	assert.Equal(t, true, time.Now().Add(14*time.Minute).Before(cookies.cookies[cookie].endTime))
}

// tests that cookie check is true on right input
func TestCheckCookie_True(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	//cookie is manually created
	allowed, authenticatedCookie := AuthenticateUser(&user, &password)
	//checked whether cookie is equal
	assert.Equal(t, true, allowed)
	isAllowed, username := CheckCookie(&authenticatedCookie)
	assert.Equal(t, true, isAllowed)
	assert.Equal(t, "admin", username)
}

func TestCheckCookie_minutesRunOut(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	//cookie is manually created
	_, authenticatedCookie := AuthenticateUser(&user, &password)
	//checked whether cookie is equal
	authenticatedCookie = authenticatedCookie[strings.Index(authenticatedCookie, "|")+1:]
	// set the timer before so it is automatically turned down
	cookies.cookies[authenticatedCookie] = authentication{
		user:    cookies.cookies[authenticatedCookie].user,
		endTime: time.Now().Add(-time.Minute * 1),
	}
	//should be discarded
	allowed, user := CheckCookie(&authenticatedCookie)
	assert.Equal(t, false, allowed)
}

// checks that authenticate user and check cookie work with each other
func TestCheckCookieAndAuthenticateUser(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	_, cookie := AuthenticateUser(&user, &password)
	isAllowed, _ := CheckCookie(&cookie)
	assert.Equal(t, true, isAllowed)
}

// tests if user already exists whether it returns an error
func TestCreateUser_AlreadyExists(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	err := CreateUser(&user, &password)
	assert.NotEqual(t, nil, err)
}

// tests if multiple users are created at the same time whether it saves them all or whether there are overlapping
func TestCreateUser_MultipleUserAtOnce(t *testing.T) {
	users.users = make(map[string]string)
	//create sync group for making done
	var wg sync.WaitGroup
	// function that checks whether the function has run through successfully
	addUser := func(name *string, passw *string) {
		err := CreateUser(name, passw)
		assert.Equal(t, nil, err)
		wg.Done()
	}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		test := "user" + strconv.Itoa(i)
		go addUser(&test, &test)
	}
	wg.Wait()
	//check at the end that the length of user is 20 and the information is secured
	assert.Equal(t, 20, len(users.users))
}

// TestCreateUser_NoSpecialChractersInUsername
// tests that if there are special characters in username that it automatically
// returns an error
func TestCreateUser_NoSpecialChractersInUsername(t *testing.T) {
	users.users = make(map[string]string)
	username := "test|"
	password := "test"
	err := CreateUser(&username, &password)
	assert.NotEqual(t, nil, err, "Create User works but shoudnt")
	assert.Equal(t, 0, len(users.users))
}

// TestCreateUser_SpecialCharacktersInPasswordAllowed
// it tests that it is allowed to put in special characters in the password
func TestCreateUser_SpecialCharacktersInPasswordAllowed(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test|$"
	err := CreateUser(&username, &password)
	assert.Equal(t, nil, err, "should work")
	assert.Equal(t, 1, len(users.users))
}

// should work correctly
func TestChangeUser(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test"
	_ = CreateUser(&username, &password)
	newPassword := "bla"
	//change password
	newCookie, err := ChangeUser(&username, &password, &newPassword)
	assert.Equal(t, err, nil)
	//check  that cookie is allowed --> the right cookie should be given back
	isCorrect, _ := AuthenticateUser(&username, &newPassword)
	assert.Equal(t, true, isCorrect)

	isCorrect, _ = CheckCookie(&newCookie)
	assert.Equal(t, true, isCorrect)
}

// should return an error if the wrong user is given
func TestChangeUser_WrongUser(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test"
	_ = CreateUser(&username, &password)
	newPassword := "bla"
	_, err := ChangeUser(&newPassword, &password, &newPassword)
	assert.NotEqual(t, nil, err)
}

// should return an error if it is the wrong old password
func TestChangeUser_WrongOldPassword(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test"
	_ = CreateUser(&username, &password)
	newPassword := "bla"
	_, err := ChangeUser(&username, &newPassword, &password)
	assert.NotEqual(t, nil, err)
}

// TestSaveUserData
// tests that it correcelty saves if there are no complications
func TestSaveUserData(t *testing.T) {
	users.users = make(map[string]string)
	wg := sync.WaitGroup{}
	wg.Add(10)
	// make multiple usernames to paralize it
	for i := 0; i < 10; i++ {
		j := i
		go func() {
			username := "test" + strconv.Itoa(j)
			password := "test" + strconv.Itoa(j)
			_ = CreateUser(&username, &password)
			wg.Done()
		}()
	}
	wg.Wait()
	path, err := filepath.Abs("../")
	// the saving the user should not return any errors
	err = SaveUserData(&path)
	assert.Equal(t, nil, err)
}

// TestSaveUserData_WriteToRightFunction
// checks that all users from the previous test have been saved correctly
func TestSaveUserData_WriteToRightFunction(t *testing.T) {
	//create user to load in
	user_test := []UserJSON{}
	//open files and check that it checks out
	path, err := filepath.Abs("../data/user")
	assert.Equal(t, err, nil)

	path = filepath.Join(path, "user-data.json")

	file, err := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			assert.Fail(t, "coudn't close file")
		}
	}(file)
	assert.Equal(t, nil, err)

	bytes, err := io.ReadAll(file)
	assert.Equal(t, nil, err)
	//derefering it and out it in the right json file
	err = json.Unmarshal(bytes, &user_test)
	assert.Equal(t, nil, err)
	//it should be of length 10
	assert.Equal(t, 10, len(user_test))
	//all 10 elements shoudnt be empty
	for _, element := range user_test {
		assert.NotEmpty(t, element)
	}

}

// TestLoadUserData
// tests the right functionality for the loada user data (uses data from 2 tests before)
func TestLoadUserData(t *testing.T) {
	//set the path so it doesn't fuck up because os.getwd doesnt return the dir the file is in
	path, _ := filepath.Abs("../")
	user := "admin"
	//load user in
	err := LoadUserData(&user, &user, &path)
	path = filepath.Join(path, "data", "user", "user-data.json")
	//loaduserdata should run without error and load 10 users created 2 tests before
	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(users.users))
	//the key and the element shoudnt be empty
	for key, element := range users.users {
		assert.NotEmpty(t, key)
		assert.NotEmpty(t, element)
	}
	// preparing for the next input
	err = os.Remove(path)
	assert.Equal(t, nil, err)
}

// TestLoadUserData_FileNotExists
// tests behaviour if the file in the given directory doesn't exist
func TestLoadUserData_FileNotExists(t *testing.T) {
	users.users = make(map[string]string)
	path, _ := filepath.Abs("../data/user")
	path = filepath.Join(path, "user-data.json")
	_ = os.Remove(path)
	user := "admin"
	path, _ = filepath.Abs("../")
	err := LoadUserData(&user, &user, &path)
	//it should create one standard user and the error should be null because it is an expected behaviour
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(users.users))
}

// TestLoadUserData_WrongFile
// tests the behaviour if the file has the same name but a different file
func TestLoadUserData_WrongFile(t *testing.T) {
	users.users = make(map[string]string)
	path, _ := filepath.Abs("../data/user")
	path = filepath.Join(path, "user-data.json")
	_ = os.Remove(path)
	user := "admin"
	file := "test"
	//preparing a wrong stated file
	err := os.WriteFile(path, []byte(file), 0644)
	assert.Equal(t, nil, err)
	path = "../"
	err = LoadUserData(&user, &user, &path)
	//it should still accept it and just add a new user
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(users.users))
	//remove the path to prevent errors on startup of application
	_ = os.Remove(path)
}

// TestCheckCookie_WrongCookieFormat
// tests that if check cookie gets a wrong string it returnss error and not allowed
func TestCheckCookie_WrongCookieFormat(t *testing.T) {
	// test first case if cookie is empty string
	wrongCookie := ""
	returnval, _ := CheckCookie(&wrongCookie)
	assert.Equal(t, false, returnval)
	// second case: cookie is not empty but has no |
	wrongCookie = "asdasd"
	returnval, _ = CheckCookie(&wrongCookie)
	assert.Equal(t, false, returnval)
	// third case: the pipe is there but at the end
	wrongCookie = "asdasd|"
	returnval, _ = CheckCookie(&wrongCookie)
	assert.Equal(t, false, returnval)

}

// TestCreateUser_PasswordOrUserNotEmpty
// checks that it returns error if one of the both are empty
func TestCreateUser_PasswordOrUserNotEmpty(t *testing.T) {
	user := "admin"
	empty := ""
	err := CreateUser(&user, &empty)
	assert.NotEqual(t, nil, err)
	err = CreateUser(&empty, &user)
	assert.NotEqual(t, nil, err)
}

// TestCheckCookie_UserNotFound
// checks that if the user is not inside it returns false
func TestCheckCookie_UserNotFound(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	CreateUser(&user, &user)
	testCookie := "user|aoishhdoüiashd"
	isRight, username := CheckCookie(&testCookie)
	assert.Equal(t, false, isRight)
	assert.Equal(t, "", username)
}

// TestLoadUserData_Wrong_Directory
// similiar to test: TestLoadUserData_WrongFile -->
// checks what happens if the open gives any other error
func TestLoadUserData_Wrong_Directory(t *testing.T) {
	users.users = make(map[string]string)
	path, _ := filepath.Abs("../data/user")
	path = filepath.Join(path, "user-data.json")
	_ = os.Remove(path)
	user := "admin"
	file := "test"
	//preparing a wrong stated file
	err := os.WriteFile(path, []byte(file), 0644)
	assert.Equal(t, nil, err)
	err = LoadUserData(&user, &user, &path)
	//it should still accept it and just add a new user
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(users.users))
	//remove the path to prevent errors on startup of application
	_ = os.Remove(path)
}

// TestDeleteOldCookies_noDeletion
// checks that it doesn't delete if they are both on time
func TestDeleteOldCookies_noDeletion(t *testing.T) {
	users.users = make(map[string]string)
	cookies.cookies = make(map[string]authentication)
	user := "admin"
	CreateUser(&user, &user)
	AuthenticateUser(&user, &user)
	user = "user"
	CreateUser(&user, &user)
	AuthenticateUser(&user, &user)
	assert.Equal(t, 2, len(cookies.cookies))
	DeleteOldCookies()
	assert.Equal(t, 2, len(cookies.cookies))

}

// TestDeleteOldCookies_DeletionBecauseToOld
// checks that it does delete those cookies who are too old
func TestDeleteOldCookies_DeletionBecauseToOld(t *testing.T) {
	// create two users
	users.users = make(map[string]string)
	cookies.cookies = make(map[string]authentication)
	user := "admin"
	CreateUser(&user, &user)
	AuthenticateUser(&user, &user)
	user = "user"
	CreateUser(&user, &user)
	// get cookie from one and subtract it so it is on the old cookie
	_, cookie := AuthenticateUser(&user, &user)
	cookie = cookie[strings.Index(cookie, "|")+1:]
	cookies.cookies[cookie] = authentication{
		user:    cookies.cookies[cookie].user,
		endTime: time.Now().Add(-1 * time.Minute),
	}
	// it should reduce by one through delete old cookies
	assert.Equal(t, 2, len(cookies.cookies))
	DeleteOldCookies()
	assert.Equal(t, 1, len(cookies.cookies))

}
