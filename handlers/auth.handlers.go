package handlers

import (
	// "fmt"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/emarifer/go-echo-templ-htmx/services"
	"github.com/emarifer/go-echo-templ-htmx/views/auth_views"
	"golang.org/x/crypto/bcrypt"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	auth_sessions_key string = "authenticate-sessions"
	auth_key          string = "authenticated"
	user_id_key       string = "user_id"
	username_key      string = "username"
	tzone_key         string = "time_zone"
)

/********** Handlers for Auth Views **********/

type AuthService interface {
	CreateUser(u services.User) error
	CheckEmail(email string) (services.User, error)
	// GetUserById(id int) (services.User, error)
}

func NewAuthHandler(us AuthService) *AuthHandler {

	return &AuthHandler{
		UserServices: us,
	}
}

type AuthHandler struct {
	UserServices AuthService
}

func (ah *AuthHandler) homeHandler(c echo.Context) error {
	homeView := auth_views.Home(fromProtected)
	isError = false

	return renderView(c, auth_views.HomeIndex(
		"| Home",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		homeView,
	))
}

func (ah *AuthHandler) registerHandler(c echo.Context) error {
	registerView := auth_views.Register(fromProtected)
	isError = false

	if c.Request().Method == "POST" {
		user := services.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
			Username: c.FormValue("username"),
		}

		err := ah.UserServices.CreateUser(user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				err = errors.New("the email is already in use")
				setFlashmessages(c, "error", fmt.Sprintf(
					"something went wrong: %s",
					err,
				))

				return c.Redirect(http.StatusSeeOther, "/register")
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "You have successfully registered!!")

		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return renderView(c, auth_views.RegisterIndex(
		"| Register",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		registerView,
	))
}

func (ah *AuthHandler) loginHandler(c echo.Context) error {
	loginView := auth_views.Login(fromProtected)
	isError = false

	if c.Request().Method == "POST" {
		// obtaining the time zone from the POST request of the login form
		tzone := ""
		if len(c.Request().Header["X-Timezone"]) != 0 {
			tzone = c.Request().Header["X-Timezone"][0]
		}

		// Authentication goes here
		user, err := ah.UserServices.CheckEmail(c.FormValue("email"))
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				setFlashmessages(c, "error", "There is no user with that email")

				return c.Redirect(http.StatusSeeOther, "/login")
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)
		if err != nil {
			// In production you have to give the user a generic message
			setFlashmessages(c, "error", "Incorrect password")

			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Get Session and setting Cookies
		sess, _ := session.Get(auth_sessions_key, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600, // in seconds
			HttpOnly: true,
		}

		// Set user as authenticated, their username,
		// their ID and the client's time zone
		sess.Values = map[interface{}]interface{}{
			auth_key:     true,
			user_id_key:  user.ID,
			username_key: user.Username,
			tzone_key:    tzone,
		}
		sess.Save(c.Request(), c.Response())

		setFlashmessages(c, "success", "You have successfully logged in!!")

		return c.Redirect(http.StatusSeeOther, "/todo/list")
	}

	return renderView(c, auth_views.LoginIndex(
		"| Login",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		loginView,
	))
}

func (ah *AuthHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			// fmt.Println(ok, auth)
			fromProtected = false

			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Please provide valid credentials")
		}

		if userId, ok := sess.Values[user_id_key].(int); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if username, ok := sess.Values[username_key].(string); ok && len(username) != 0 {
			c.Set(username_key, username) // set the username in the context
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		fromProtected = true

		return next(c)
	}
}

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
