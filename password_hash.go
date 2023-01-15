package web

import "golang.org/x/crypto/bcrypt"

// See https://security.stackexchange.com/questions/211/how-to-securely-hash-passwords/31846#31846
func passwordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
