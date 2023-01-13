package web

import (
	"encoding/json"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// AuthDB is a service for authentication.
// It is intended to be set up as a service that is used by other services and not end users.
// It is intended to be used on the open internet.
type AuthDB struct {
	// Users are the users of the service.
	Users map[string]*User `json:"users"`
	// PasswordSalt is used to salt passwords.
	PasswordSalt []byte `json:"password_salt"`
	// RegistrationCodes are needed to create new users.
	// Codes can be created by sending a POST request to /register with a valid registrar key.
	RegistrationCodes map[string]bool `json:"registration_codes"`
	// RegistrarKey is a secret key that is needed to generate registration codes.
	// If your service is open to new signups, you can share this string publicly.
	RegistrarKey string `json:"registrar_key"`
	// locks are used to prevent race conditions.
	locks map[string]*sync.Mutex
}

// Invite creates a new registration code and returns it.
func (s *AuthDB) Invite() string {
	registrationCode := s.newID("registration_code")
	s.RegistrationCodes[registrationCode] = true
	return registrationCode
}

// Register creates a new user with the given password and returns the new user's ID.
// If it returns an error, it will be of type ErrInvalidRegistrationCode.
func (s *AuthDB) Register(registrationCode, password string) (string, error) {
	if !s.RegistrationCodes[registrationCode] {
		return "", ErrInvalidRegistrationCode
	}
	delete(s.RegistrationCodes, registrationCode)
	userID := s.newID("user")
	user := &User{
		PasswordHash:  s.passwordHash(password),
		SessionTokens: map[string]bool{},
	}
	s.Users[userID] = user
	return userID, nil
}

// Login returns a new session token for the given user ID and password.
// If it returns an error, it will be of type ErrInvalidPassword or ErrUserNotFound.
func (s *AuthDB) Login(userID, password string) (string, error) {
	user, ok := s.Users[userID]
	if !ok {
		return "", ErrUserNotFound
	}
	if !s.passwordHashMatches(password, user.PasswordHash) {
		return "", ErrInvalidPassword
	}
	sessionToken := s.newID("session_token")
	user.SessionTokens[sessionToken] = true
	return sessionToken, nil
}

// User returns the User with the given ID.
func (s *AuthDB) User(id string) (*User, bool) {
	u, ok := s.Users[id]
	return u, ok
}

// ServeHTTP serves the authentication server.
// There are 5 endpoints: /invite, /register, /login, /logout, and /user.
func (s *AuthDB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/invite":
		s.invite(w, r)
	case "/register":
		s.register(w, r)
	case "/login":
		s.login(w, r)
	case "/logout":
		s.logout(w, r)
	case "/user":
		s.user(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// invite handles an invite request.
// It creates a new registration code and returns it.
func (s *AuthDB) invite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RegistrarKey string `json:"registrar_key"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if req.RegistrarKey != s.RegistrarKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	registrationCode := s.Invite()
	json.NewEncoder(w).Encode(registrationCode)
}

// register handles a registration request.
// It creates a new user with the given password and returns the new user's ID.
func (s *AuthDB) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RegistrationCode string `json:"registration_code"`
		Password         string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	userID, err := s.Register(req.RegistrationCode, req.Password)
	json.NewEncoder(w).Encode(userID)
}

// login handles a login request.
// It returns a session token that can be used to validate a session.
// The session token is valid until the user logs out.
func (s *AuthDB) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequeset
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, ok := s.Users[req.UserID]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if user.PasswordHash != s.passwordHash(req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token := s.newID(req.UserID)
	user.SessionTokens[token] = true
	json.NewEncoder(w).Encode(&Session{
		UserID: req.UserID,
		Token:  token,
	})
}

// logout handles a logout request.
// It invalidates the given session token if it is valid.
func (s *AuthDB) logout(w http.ResponseWriter, r *http.Request) {
	var session Session
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, ok := s.Users[session.UserID]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	delete(user.SessionTokens, session.Token)
}

// validateSession handles a session validation request.
// Session validation requests come from apps to determine if a user token is valid.
func (server *AuthDB) validateSession(w http.ResponseWriter, r *http.Request) {
	var session Session
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, ok := server.Users[session.UserID]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok = user.SessionTokens[session.Token]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// user handles requests for user info.
func (s *AuthDB) user(w http.ResponseWriter, r *http.Request) {
	var session Session
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, ok := s.User(session.UserID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok = user.SessionTokens[session.Token]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// newID generates a new unique ID for the given scope.
func (s *AuthDB) newID(scope string) string {
	s.locks[scope].Lock()
	defer s.locks[scope].Unlock()
	return NewID()
}

// See https://security.stackexchange.com/questions/211/how-to-securely-hash-passwords/31846#31846
func (s *AuthDB) passwordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// generateRegistrationCode generates a new registration code.
func (s *AuthDB) generateRegistrationCode() string {
	code := s.newID("registration_code")
	s.RegistrationCodes[code] = true
	return code
}

// passwordHashMatches returns true if the given password hash matches the given password.
func (s *AuthDB) passwordHashMatches(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
