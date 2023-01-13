package web

type User struct {
	PasswordHash  string          `json:"password_hash"`
	SessionTokens map[string]bool `json:"session_tokens"`
	Emails        []string        `json:"emails"`
	Orgs          map[string]bool `json:"orgs"`
}

// PrimaryEmail returns the primary email address of the user.
// It panic if the user has no email addresses.
func (u *User) PrimaryEmail() string {
	return u.Emails[0]
}
