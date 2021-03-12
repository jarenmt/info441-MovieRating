package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/assignments-fixed-ssunni12/servers/gateway/models/users"
	"github.com/assignments-fixed-ssunni12/servers/gateway/sessions"
	"github.com/go-redis/redis"
)

func TestContext(t *testing.T) {
	//Fail
	var first *redis.Client
	second := sessions.NewRedisStore(first, time.Second) //sessions.NewMemStore(time.Hour, time.Hour)
	third := &users.MySQLStore{}                         // &users.MyMockStore{}
	context := NewContext("", second, third)
	if context != nil {
		t.Error("Expected Context constructor to fail but it didn't return nil")
	}

	//Success
	a := redis.NewClient(&redis.Options{
		Addr: "172.17.0.2:6379",
	})
	b := sessions.NewRedisStore(a, time.Second) //sessions.NewMemStore(time.Hour, time.Hour)
	contextSuccess := NewContext("test", b, third)
	if contextSuccess == nil {
		t.Error("Expected Context constructor to work but it didn't")
	}
}

func TestContextPOSTUserHandler(t *testing.T) {
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "172.17.0.2:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	context := &Context{
		Key:           "testkey",
		SessionsStore: sessions.NewRedisStore(client, time.Hour), //sessions.NewMemStore(time.Hour, time.Hour),
		UsersStore:    &users.MySQLStore{},                       //users.NewMockUserDB(),
	}
	newUserValid := users.NewUser{
		Email:        "test@example.com",
		Password:     "password1",
		PasswordConf: "password1",
		UserName:     "TestUser1",
		FirstName:    "Adam",
		LastName:     "Smith",
	}
	newUserInvalid := users.NewUser{
		Email:        "test@example.com",
		Password:     "password2",
		PasswordConf: "Wrong",
		UserName:     "TestUser2",
		FirstName:    "Adam",
		LastName:     "Error",
	}
	newUserFailToInsert := users.NewUser{
		Email:        "test@example.com",
		Password:     "password3",
		PasswordConf: "password3",
		UserName:     "TestUser3",
		FirstName:    "Adam",
		LastName:     "Error",
	}
	NewUserFailToInsert2 := users.NewUser{
		Email:        "test@example.com",
		Password:     "password4",
		PasswordConf: "password4",
		UserName:     "TestUser4",
		FirstName:    "Error",
		LastName:     "Error",
	}
	validUser, err := newUserValid.ToUser()
	if err != nil {
		log.Printf("Error initializing test users")
	}
	validUser.ID = int64(1)

	cases := []struct {
		name                string
		method              string
		idPath              string
		newUser             *users.NewUser
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			"Valid POST Request",
			http.MethodPost,
			"1",
			&newUserValid,
			http.StatusCreated,
			false,
			"application/json",
			validUser,
		},
		{
			"Invalid POST Request - Content Type",
			http.MethodPost,
			"2",
			&newUserValid,
			http.StatusUnsupportedMediaType,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid POST Request - Invalid NewUser",
			http.MethodPost,
			"2",
			&newUserInvalid,
			http.StatusInternalServerError,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid Provided User - Insert Fail",
			http.MethodPost,
			"2",
			&newUserFailToInsert,
			http.StatusInternalServerError,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid Provided User - Insert Fail (User ID)",
			http.MethodPost,
			"2",
			&NewUserFailToInsert2,
			http.StatusInternalServerError,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
	}

	for _, c := range cases {
		log.Printf("\n\nRunning case name: %s", c.name)
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.Encode(c.newUser)
		request := httptest.NewRequest(c.method, "/v1/users/"+c.idPath, buffer)
		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			request.Header.Add("Content-Type", "application/json")
		}
		recorder := httptest.NewRecorder()
		context.UsersHandler(recorder, request)
		response := recorder.Result()

		responseContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != responseContentType {
			t.Errorf("Case %s: incorrect return type, expected %s but received %s",
				c.name, c.expectedContentType, responseContentType)
		}

		responseStatusCode := response.StatusCode
		if c.expectedStatusCode != responseStatusCode {
			t.Errorf("Case %s: incorrect status code, expected %d but received %d",
				c.name, c.expectedStatusCode, responseStatusCode)
		}

		user := &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedError && err == nil {
			t.Errorf("Case %s: expected error, but received none", c.name)
		}

		if !c.expectedError && c.expectedReturn.Email != user.Email && string(c.expectedReturn.PassHash) != string(user.PassHash) &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("Case %s: incorrect return, expected %v but received %v",
				c.name, c.expectedReturn, user)
		}
	}

}

func TestContext_GETSpecificUserHandler(t *testing.T) {
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "172.17.0.2:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	context := &Context{
		Key:           "testkey",
		SessionsStore: sessions.NewRedisStore(client, time.Hour), //sessions.NewMemStore(time.Hour, time.Hour),
		UsersStore:    &users.MySQLStore{},                       //users.NewMockUserDB(),
	}

	newUserValid := users.NewUser{
		Email:        "test@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "TestUser1",
		FirstName:    "Adam",
		LastName:     "Smith",
	}
	validUser, err := newUserValid.ToUser()
	if err != nil {
		log.Printf("Error initializing test users")
	}
	validUser.ID = int64(1)
	context.UsersStore.Insert(validUser)

	cases := []struct {
		name                string
		method              string
		idPath              string
		newUser             *users.NewUser
		useValidCredentials bool
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			"Valid GET Request",
			http.MethodGet,
			"1",
			&newUserValid,
			true,
			http.StatusOK,
			false,
			"application/json",
			validUser,
		},
		{
			"Valid GET Request - Me",
			http.MethodGet,
			"me",
			&newUserValid,
			true,
			http.StatusOK,
			false,
			"application/json",
			validUser,
		},
		{
			"Invalid GET Request - No Credentials (No Session)",
			http.MethodGet,
			"1",
			&newUserValid,
			false,
			http.StatusUnauthorized,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid GET Request - Wrong Method",
			http.MethodPost,
			"1",
			&newUserValid,
			true,
			http.StatusMethodNotAllowed,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid GET Request - Requested User Not Found",
			http.MethodGet,
			"2",
			&newUserValid,
			true,
			http.StatusNotFound,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid GET Request - No Provided ID",
			http.MethodGet,
			"",
			&newUserValid,
			true,
			http.StatusNotFound,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
	}

	for _, c := range cases {
		log.Printf("\n\nRunning case name: %s", c.name)
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.Encode(c.newUser)
		request := httptest.NewRequest(c.method, "/v1/users/"+c.idPath, buffer)

		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			request.Header.Add("Content-Type", "application/json")
		}

		recorder := httptest.NewRecorder()
		if c.name != "Invalid GET Request - No Credentials (No Session)" {
			newSessionState := SessionState{time.Now(), *c.expectedReturn}
			sid, err := sessions.BeginSession(context.Key, context.SessionsStore, &newSessionState, recorder)
			if err != nil {
				log.Fatal(err.Error())
			}
			request.Header.Add("Authorization", "Bearer "+sid.String())
		}
		context.SpecificUserHandler(recorder, request)
		response := recorder.Result()
		responseContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != responseContentType {
			t.Errorf("Case %s: incorrect return type, expected %s but received %s",
				c.name, c.expectedContentType, responseContentType)
		}

		responseStatusCode := response.StatusCode
		if c.expectedStatusCode != responseStatusCode {
			t.Errorf("Case %s: incorrect status code, expected %d but received %d",
				c.name, c.expectedStatusCode, responseStatusCode)
		}

		user := &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedError && err == nil {
			t.Errorf("Case %s: expected error, but received none", c.name)
		}

		if !c.expectedError && c.expectedReturn.Email != user.Email && string(c.expectedReturn.PassHash) != string(user.PassHash) &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("Case %s: incorrect return, expected %v but received %v",
				c.name, c.expectedReturn, user)
		}
	}
}

func TestContext_PATCHSpecificUserHandler(t *testing.T) {
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "172.17.0.2:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	context := &Context{
		Key:           "testkey",
		SessionsStore: sessions.NewRedisStore(client, time.Hour), //sessions.NewMemStore(time.Hour, time.Hour),
		UsersStore:    &users.MySQLStore{},                       //users.NewMockUserDB(),
	}
	newUserValid := users.NewUser{
		Email:        "test@example.com",
		Password:     "password1",
		PasswordConf: "password1",
		UserName:     "TestUser1",
		FirstName:    "Adam",
		LastName:     "Smith",
	}

	validUser, err := newUserValid.ToUser()
	if err != nil {
		log.Printf("Error initializing test users")
	}
	validUser.ID = int64(1)
	context.UsersStore.Insert(validUser)
	validUpdates := users.Updates{
		FirstName: "Test",
		LastName:  "Smith",
	}
	invalidUpdates := users.Updates{
		FirstName: "Error",
		LastName:  "Smith",
	}

	cases := []struct {
		name                string
		method              string
		idPath              string
		updates             *users.Updates
		useValidCredentials bool
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			"Valid PATCH Request",
			http.MethodPatch,
			"1",
			&validUpdates,
			true,
			http.StatusOK,
			false,
			"application/json",
			validUser,
		},
		{
			"Invalid PATCH Request - Invalid Header",
			http.MethodPatch,
			"1",
			&invalidUpdates,
			true,
			http.StatusUnsupportedMediaType,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid Patch Request - Trying to update another user",
			http.MethodPatch,
			"1000",
			&invalidUpdates,
			true,
			http.StatusForbidden,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid Patch Request - No ID given",
			http.MethodPatch,
			"",
			&invalidUpdates,
			true,
			http.StatusNotFound,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid Patch Request - Invalid header",
			http.MethodPatch,
			"1",
			&invalidUpdates,
			true,
			http.StatusUnsupportedMediaType,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
	}

	for _, c := range cases {
		log.Printf("\n\nRunning case name: %s", c.name)
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.Encode(c.updates)
		request := httptest.NewRequest(c.method, "/v1/users/"+c.idPath, buffer)
		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			request.Header.Add("Content-Type", "application/json")
		}

		recorder := httptest.NewRecorder()
		if c.name != "Invalid Get Request - No Credentials (no session)" {
			newSessionState := SessionState{time.Now(), *c.expectedReturn}
			sid, err := sessions.BeginSession(context.Key, context.SessionsStore, &newSessionState, recorder)
			if err != nil {
				log.Fatal(err.Error())
			}
			request.Header.Add("Authorization", "Bearer "+sid.String())

		}
		context.SpecificUserHandler(recorder, request)
		response := recorder.Result()

		responseContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != responseContentType {
			t.Errorf("Case %s: incorrect return type, expected %s but received %s",
				c.name, c.expectedContentType, responseContentType)
		}

		responseStatusCode := response.StatusCode
		if c.expectedStatusCode != responseStatusCode {
			t.Errorf("Case %s: incorrect status code, expected %d but received %d",
				c.name, c.expectedStatusCode, responseStatusCode)
		}

		user := &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedError && err == nil {
			t.Errorf("Case %s: expected error, but received none", c.name)
		}

		if !c.expectedError && c.expectedReturn.Email != user.Email && string(c.expectedReturn.PassHash) != string(user.PassHash) &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("Case %s: incorrect return, expected %v but received %v",
				c.name, c.expectedReturn, user)
		}
	}
}

func TestContext_SessionsHandler(t *testing.T) {
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "172.17.0.2:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	context := &Context{
		Key:           "testkey",
		SessionsStore: sessions.NewRedisStore(client, time.Hour), //sessions.NewMemStore(time.Hour, time.Hour),
		UsersStore:    &users.MySQLStore{},                       //users.NewMockUserDB(),
	}
	validCredentials := users.Credentials{
		Email:    "test@example.com",
		Password: "password",
	}

	invalidCredentialsPassword := users.Credentials{
		Email:    "test@example.com",
		Password: "wrong",
	}
	invalidCredentialsEmail := users.Credentials{
		Email:    "wrong@example.com",
		Password: "password",
	}
	var badCredentials *users.Credentials

	validNewUser := users.NewUser{
		Email:        "test@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "TestUser1",
		FirstName:    "Adam",
		LastName:     "Smith",
	}
	validUser, err := validNewUser.ToUser()
	if err != nil {
		log.Printf("Error initializing test users")
	}
	validUser.ID = int64(1)
	context.UsersStore.Insert(validUser)

	cases := []struct {
		name                string
		method              string
		idPath              string
		credentials         *users.Credentials
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			"Valid POST Request",
			http.MethodPost,
			"1",
			&validCredentials,
			http.StatusOK,
			false,
			"application/json",
			validUser,
		},
		{
			"Valid POST Request - Wrong Credentials (Email)",
			http.MethodPost,
			"1",
			&invalidCredentialsEmail,
			http.StatusUnauthorized,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Valid POST Request - Wrong Credentials (Password)",
			http.MethodPost,
			"1",
			&invalidCredentialsPassword,
			http.StatusUnauthorized,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid POST Request - Wrong Method",
			http.MethodGet,
			"1",
			&validCredentials,
			http.StatusMethodNotAllowed,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid POST Request - Wrong Header",
			http.MethodPost,
			"1",
			&validCredentials,
			http.StatusUnsupportedMediaType,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
		{
			"Invalid POST Request - Bad Credentials",
			http.MethodPost,
			"1",
			badCredentials,
			http.StatusInternalServerError,
			true,
			"text/plain; charset=utf-8",
			validUser,
		},
	}

	for _, c := range cases {
		log.Printf("Running case name: %s", c.name)
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		if c.credentials != badCredentials {
			encoder.Encode(c.credentials)
		}
		request := httptest.NewRequest(c.method, "/v1/users/"+c.idPath, buffer)

		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			request.Header.Add("Content-Type", "application/json")
		}

		recorder := httptest.NewRecorder()
		context.SessionsHandler(recorder, request)
		response := recorder.Result()

		responseContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != responseContentType {
			t.Errorf("Case %s: incorrect return type, expected %s but received %s",
				c.name, c.expectedContentType, responseContentType)
		}

		responseStatusCode := response.StatusCode
		if c.expectedStatusCode != responseStatusCode {
			t.Errorf("Case %s: incorrect status code, expected %d but received %d",
				c.name, c.expectedStatusCode, responseStatusCode)
		}

		user := &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedError && err == nil {
			t.Errorf("Case %s: expected error, but received none", c.name)
		}

		if !c.expectedError && c.expectedReturn.Email != user.Email && string(c.expectedReturn.PassHash) != string(user.PassHash) &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("Case %s: incorrect return, expected %v but received %v",
				c.name, c.expectedReturn, user)
		}
	}
}

func TestContext_DELETESpecificSessionHandler(t *testing.T) {
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "172.17.0.2:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	context := &Context{
		Key:           "testkey",
		SessionsStore: sessions.NewRedisStore(client, time.Hour), //sessions.NewMemStore(time.Hour, time.Hour),
		UsersStore:    &users.MySQLStore{},                       //users.NewMockUserDB(),
	}
	validNewUser := users.NewUser{
		Email:        "test@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "TestUser1",
		FirstName:    "Adam",
		LastName:     "Smith",
	}
	validUser, err := validNewUser.ToUser()
	if err != nil {
		log.Printf("Error initializing test users")
	}
	validUser.ID = int64(1)
	context.UsersStore.Insert(validUser)

	cases := []struct {
		name                string
		method              string
		idPath              string
		useValidCredentials bool
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      string
	}{
		{
			"Valid DELETE Request",
			http.MethodDelete,
			"mine",
			true,
			http.StatusOK,
			false,
			"text/plain; charset=utf-8",
			"Signed Out",
		},
		{
			"Invalid DELETE Request - Invalid Method",
			http.MethodPatch,
			"1",
			true,
			http.StatusMethodNotAllowed,
			true,
			"text/plain; charset=utf-8",
			"Expected Error",
		},
		{
			"Invalid DELETE Request - Trying To Log Out Another User",
			http.MethodDelete,
			"2",
			true,
			http.StatusForbidden,
			true,
			"text/plain; charset=utf-8",
			"Expected Error",
		},
		{
			"Invalid DELETE Request - No Existing Session",
			http.MethodDelete,
			"2",
			false,
			http.StatusUnauthorized,
			true,
			"text/plain; charset=utf-8",
			"Expected Error",
		},
	}

	for _, c := range cases {
		log.Printf("\n\nRunning case name: %s", c.name)
		request := httptest.NewRequest(c.method, "/v1/sessions/"+c.idPath, nil)
		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			request.Header.Add("Content-Type", "application/json")
		}

		recorder := httptest.NewRecorder()
		if c.useValidCredentials {
			newSessionState := SessionState{time.Now(), *validUser}
			sid, err := sessions.BeginSession(context.Key, context.SessionsStore, &newSessionState, recorder)
			if err != nil {
				log.Fatal(err.Error())
			}
			request.Header.Add("Authorization", "Bearer "+sid.String())

		}
		context.SpecificSessionHandler(recorder, request)
		response := recorder.Result()

		responseContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != responseContentType {
			t.Errorf("Case %s: incorrect return type, expected %s but received %s",
				c.name, c.expectedContentType, responseContentType)
		}

		responseStatusCode := response.StatusCode
		if c.expectedStatusCode != responseStatusCode {
			t.Errorf("Case %s: incorrect status code, expected %d but received %d",
				c.name, c.expectedStatusCode, responseStatusCode)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Failed trying to read response data")
		}
		responseString := string(responseData)
		log.Printf("Returned string was: %s", responseString)

		if !c.expectedError && c.expectedReturn != "Signed Out" {
			t.Errorf("Case %s: incorrect return, expected %v but received %v",
				c.name, c.expectedReturn, responseString)
		}
	}
}
