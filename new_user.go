package web

func NewUser(password string) *User {
	return &User{
		PasswordHash:  passwordHash(password),
		SessionTokens: map[string]bool{},
		Emails:        []string{},
	}
}
