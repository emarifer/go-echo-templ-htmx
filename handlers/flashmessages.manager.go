package handlers

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// cookie name & flash messages key
// this should be a .env file
const (
	session_name              string = "fmessages"
	session_flashmessages_key string = "flashmessages-key"
)

func getCookieStore() *sessions.CookieStore {

	return sessions.NewCookieStore([]byte(session_flashmessages_key))
}

// Set adds a new message to the cookie store
func setFlashmessages(c echo.Context, kind, value string) {
	session, _ := getCookieStore().Get(c.Request(), session_name)

	session.AddFlash(value, kind)

	session.Save(c.Request(), c.Response())
}

// Get receives flash messages from cookie store
func getFlashmessages(c echo.Context, kind string) []string {
	session, _ := getCookieStore().Get(c.Request(), session_name)

	fm := session.Flashes(kind)

	// if there are some messages…
	if len(fm) > 0 {
		session.Save(c.Request(), c.Response())

		// we start an empty strings slice that we
		// then return with messages
		var flashes []string
		for _, fl := range fm {
			// we add the messages to the slice
			flashes = append(flashes, fl.(string))
		}

		return flashes
	}

	return nil
}
