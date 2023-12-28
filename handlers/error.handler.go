package handlers

import (
	"fmt"
	"net/http"

	"github.com/emarifer/go-echo-templ-htmx/views/errors_pages"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	/* errorPage := fmt.Sprintf("views/%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	} */

	var errorPage func(fp bool) templ.Component

	switch code {
	case 401:
		errorPage = errors_pages.Error401
	case 404:
		errorPage = errors_pages.Error404
	case 500:
		errorPage = errors_pages.Error500
	}

	isError = true

	renderView(c, errors_pages.ErrorIndex(
		fmt.Sprintf("| Error (%d)", code),
		"",
		fromProtected,
		isError,
		errorPage(fromProtected),
	))
}
