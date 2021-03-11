package handlers

import (
	"time"

	"github.com/assignments-fixed-ssunni12/servers/gateway/models/users"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!
type SessionState struct {
  SessionBegin time.Time `json:"sessionBegin"`
  User users.User `json:"user"`
}
