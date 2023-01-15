package web

import "golang.org/x/crypto/bcrypt"

// passwordHashMatches returns true if the given password hash matches the given password.
func passwordHashMatches(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
