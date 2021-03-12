package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/assignments-fixed-ssunni12/servers/gateway/models/users"
	"github.com/assignments-fixed-ssunni12/servers/gateway/sessions"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.
func (cont *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Server directed towards UsersHandler")
	if r.Method == "POST" {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			log.Println("Header content type confirmed to be json")
			decoder := json.NewDecoder(r.Body)
			var newUser users.NewUser
			err := decoder.Decode(&newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = newUser.Validate()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u, err := newUser.ToUser()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			userResponse, err := cont.UsersStore.Insert(u)
			if err != nil {
				log.Printf("Insert failed, error was: %v", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if userResponse.ID == 0 {
				log.Printf("ID not assigned, error was: %v", err.Error())
				http.Error(w, "Database didn't assign ID to user", http.StatusBadRequest)
				return
			}
			newSessionState := SessionState{time.Now(), *userResponse}
			_, err = sessions.BeginSession(cont.Key, cont.SessionsStore, newSessionState, w)
			if err != nil {
				log.Println("Failed to create session")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			encoder := json.NewEncoder(w)
			err = encoder.Encode(userResponse)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		} else {
			http.Error(w, "JSON is only supported media type", http.StatusUnsupportedMediaType)
			return
		}
	} else {
		http.Error(w, "POST is the only method supported by this handler", http.StatusMethodNotAllowed)
		return
	}
}

func (cont *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Server directed towards SpecificUsersHandler")
	currentSession := &SessionState{}
	_, err := sessions.GetState(r, cont.Key, cont.SessionsStore, currentSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if r.Method == "GET" {
		providedUser, err := getUserID(r, currentSession.User.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		returnedUser, err := cont.UsersStore.GetByID(providedUser)
		if err != nil {
			http.Error(w, "User was not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(returnedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == "PATCH" {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			providedUser, err := getUserID(r, currentSession.User.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			if currentSession.User.ID == providedUser {
				decoder := json.NewDecoder(r.Body)
				var providedUpdates users.Updates
				err := decoder.Decode(&providedUpdates)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				updatedUser, err := cont.UsersStore.Update(providedUser, &providedUpdates)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				preUpdateUser := currentSession.User
				cont.UsersStore.Delete(preUpdateUser.ID) //added this
				cont.UsersStore.Insert(updatedUser)      //added this
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				encoder := json.NewEncoder(w)
				encoder.Encode(updatedUser)
				return
			}
			http.Error(w, "Can't update a user that isn't yourself", http.StatusForbidden)
			return
		}
		http.Error(w, "JSON is  only supported media type", http.StatusUnsupportedMediaType)
		return
	}
	http.Error(w, "GET and PATCH are the only methods supported by this handler", http.StatusMethodNotAllowed)
	return
}

func (cont *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Server directed towards SessionsHandler")
	if r.Method == "POST" {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			log.Println("Attempted Sign In")
			var providedCreds users.Credentials
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&providedCreds)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			providedUser, err := cont.UsersStore.GetByEmail(providedCreds.Email)
			if err != nil {
				time.Sleep(2 * time.Second)
				http.Error(w, "Invalid Credentials Provided", http.StatusUnauthorized)
				return
			}
			err = providedUser.Authenticate(providedCreds.Password)
			if err != nil {
				http.Error(w, "Invalid Credentials Provided", http.StatusUnauthorized)
				return
			}
			newSessionState := SessionState{time.Now(), *providedUser}
			_, err = sessions.BeginSession(cont.Key, cont.SessionsStore, newSessionState, w)
			if err != nil {
				http.Error(w, "Invalid Credentials Provided", http.StatusUnauthorized)
				return
			}
			_, err = cont.UsersStore.InsertSignIn(providedUser.ID, r.RemoteAddr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			encoder := json.NewEncoder(w)
			err = encoder.Encode(providedUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		http.Error(w, "JSON is only supported media type", http.StatusUnsupportedMediaType)
		return
	}
	http.Error(w, "POST is the only method supported by this handler", http.StatusMethodNotAllowed)
	return
}

func (cont *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Server directed towards SpecificSessionHandler")
	currentSession := &SessionState{}
	_, err := sessions.GetState(r, cont.Key, cont.SessionsStore, currentSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if r.Method == "DELETE" {
		finalSegment := path.Base(r.URL.Path)
		if finalSegment == "mine" {
			_, err := sessions.EndSession(r, cont.Key, cont.SessionsStore)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Write([]byte("Signed Out"))
			return
		}
		http.Error(w, "Cannot end the session of another user", http.StatusForbidden)
		return
	}
	http.Error(w, "DELETE is the only method supported by this handler", http.StatusMethodNotAllowed)
	return
}

func getUserID(r *http.Request, authenticatedUser int64) (int64, error) {
	requestedUser := path.Base(r.URL.Path)
	if requestedUser == "me" {
		return authenticatedUser, nil
	} else {
		finalRequest, err := strconv.Atoi(requestedUser)
		if err != nil {
			return 0, err
		}
		resultUser := int64(finalRequest)
		return resultUser, nil
	}
}
