package handlers

import (
	"github.com/assignments-fixed-ssunni12/servers/gateway/models/users"
	"github.com/assignments-fixed-ssunni12/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store
type Context struct {
	Key           string
	SessionsStore *sessions.RedisStore /* *sessions.MemStore */
	UsersStore    *users.MySQLStore    /* *users.MyMockStore   */
}

func NewContext(key string, session *sessions.RedisStore /* *sessions.MemStore */, user *users.MySQLStore /* *users.MyMockStore  */) *Context {
	if session == nil || user == nil || key == "" {
		return nil
	}
	context := Context{key, session, user}
	return &context
}
