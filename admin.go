package web

import (
	"net/http"
	"net/url"
)

// Admin is a web administrator.
type Admin struct {
	PhoneNumber string

	TwillioAccountSID  string
	TwillioAuthToken   string
	TwillioPhoneNumber string
}

// Notify sends a text message to the admin.
func (a *Admin) Notify(message string) error {
	u := url.URL{
		Scheme: "https",
		Host:   "api.twilio.com",
		Path:   "/2010-04-01/Accounts/" + a.TwillioAccountSID + "/Messages.json",
		RawQuery: url.Values{
			"From": {a.TwillioPhoneNumber},
			"To":   {a.PhoneNumber},
			"Body": {message},
		}.Encode(),
	}
	res, err := http.Post(u.String(), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
