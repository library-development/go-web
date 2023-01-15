package web

func NewAuthDB(adminPassword, registrarKey string) *AuthDB {
	authDB := &AuthDB{
		Users: map[string]*User{
			"admin": NewUser(adminPassword),
		},
		RegistrationCodes: map[string]bool{},
		RegistrarKey:      registrarKey,
		AdminID:           "admin",
	}
	return authDB
}
