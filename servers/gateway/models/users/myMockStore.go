package users

// MockUserStore mocks users.Store to use in testing
type MyMockStore struct {
	currID      int64
	idMap       map[int64]*User
	emailMap    map[string]*User
	usernameMap map[string]*User
}

// NewMockUserDB create a mock UserStore according to the users.Store interface
func NewMockUserDB() *MyMockStore {
	return &MyMockStore{1, make(map[int64]*User), make(map[string]*User), make(map[string]*User)}
}

//GetByID returns the User with the given ID
func (mus *MyMockStore) GetByID(id int64) (*User, error) {
	if _, exists := mus.idMap[id]; exists {
		return mus.idMap[id], nil
	}
	return nil, ErrUserNotFound
}

//GetByEmail returns the User with the given email
func (mus *MyMockStore) GetByEmail(email string) (*User, error) {
	if _, exists := mus.emailMap[email]; exists {
		return mus.emailMap[email], nil
	}
	return nil, ErrUserNotFound
}

//GetByUserName returns the User with the given Username
func (mus *MyMockStore) GetByUserName(username string) (*User, error) {
	if _, exists := mus.usernameMap[username]; exists {
		return mus.usernameMap[username], nil
	}
	return nil, ErrUserNotFound
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (mus *MyMockStore) Insert(user *User) (*User, error) {
	id := mus.currID
	email := user.Email
	username := user.UserName
	// email already exists, return error
	if _, emailExists := mus.emailMap[email]; emailExists {
		return nil, ErrUserNotFound
	}
	// username already exists, return error
	if _, usernameExists := mus.usernameMap[username]; usernameExists {
		return nil, ErrUserNotFound
	}
	user.ID = id
	mus.idMap[id] = user
	mus.emailMap[email] = user
	mus.usernameMap[username] = user
	mus.currID++
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (mus *MyMockStore) Update(id int64, updates *Updates) (*User, error) {
	user, _ := mus.GetByID(id)
	err := user.ApplyUpdates(updates)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Delete deletes the user with the given ID
func (mus *MyMockStore) Delete(id int64) error {
	user, err := mus.GetByID(id)
	if err != nil {
		return err
	}
	email := user.Email
	username := user.UserName
	delete(mus.idMap, id)
	delete(mus.emailMap, email)
	delete(mus.usernameMap, username)
	return nil
}
