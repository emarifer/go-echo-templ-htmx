package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emarifer/go-echo-templ-htmx/services"
	"github.com/emarifer/go-echo-templ-htmx/views/todo_views"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/********** Handlers for Todo Views **********/

type TaskService interface {
	CreateTodo(t services.Todo) (services.Todo, error)
	GetAllTodos(createdBy int) ([]services.Todo, error)
	GetTodoById(t services.Todo) (services.Todo, error)
	UpdateTodo(t services.Todo) (services.Todo, error)
	DeleteTodo(t services.Todo) error
}

func NewTaskHandler(ts TaskService) *TaskHandler {

	return &TaskHandler{
		TodoServices: ts,
	}
}

type TaskHandler struct {
	TodoServices TaskService
}

func (th *TaskHandler) createTodoHandler(c echo.Context) error {
	isError = false

	if c.Request().Method == "POST" {
		todo := services.Todo{
			CreatedBy:   c.Get(user_id_key).(int),
			Title:       strings.Trim(c.FormValue("title"), " "),
			Description: strings.Trim(c.FormValue("description"), " "),
		}

		_, err := th.TodoServices.CreateTodo(todo)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Task created successfully!!")

		return c.Redirect(http.StatusSeeOther, "/todo/list")
	}

	return renderView(c, todo_views.TodoIndex(
		"| Create Todo",
		c.Get(username_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		todo_views.CreateTodo(),
	))
}

func (th *TaskHandler) todoListHandler(c echo.Context) error {
	isError = false
	userId := c.Get(user_id_key).(int)

	todos, err := th.TodoServices.GetAllTodos(userId)
	if err != nil {
		return err
	}

	titlePage := fmt.Sprintf(
		"| %s's Task List",
		cases.Title(language.English).String(c.Get(username_key).(string)),
	)

	return renderView(c, todo_views.TodoIndex(
		titlePage,
		c.Get(username_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		todo_views.TodoList(titlePage, todos),
	))
}

func (th *TaskHandler) updateTodoHandler(c echo.Context) error {
	isError = false

	idParams, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	t := services.Todo{
		ID:        idParams,
		CreatedBy: c.Get(user_id_key).(int),
	}

	todo, err := th.TodoServices.GetTodoById(t)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {

			return echo.NewHTTPError(
				echo.ErrNotFound.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	if c.Request().Method == "POST" {
		var status bool
		if c.FormValue("status") == "on" {
			status = true
		} else {
			status = false
		}

		todo := services.Todo{
			Title:       strings.Trim(c.FormValue("title"), " "),
			Description: strings.Trim(c.FormValue("description"), " "),
			Status:      status,
			CreatedBy:   c.Get(user_id_key).(int),
			ID:          idParams,
		}

		_, err := th.TodoServices.UpdateTodo(todo)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Task successfully updated!!")

		return c.Redirect(http.StatusSeeOther, "/todo/list")
	}

	return renderView(c, todo_views.TodoIndex(
		fmt.Sprintf("| Edit Todo #%d", todo.ID),
		c.Get(username_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"), // ↓ getting time zone from context ↓
		todo_views.UpdateTodo(todo, c.Get(tzone_key).(string)),
	))
}

func (th *TaskHandler) deleteTodoHandler(c echo.Context) error {
	idParams, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	t := services.Todo{
		CreatedBy: c.Get(user_id_key).(int),
		ID:        idParams,
	}

	err = th.TodoServices.DeleteTodo(t)
	if err != nil {
		if strings.Contains(err.Error(), "an affected row was expected") {

			return echo.NewHTTPError(
				echo.ErrNotFound.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	setFlashmessages(c, "success", "Task successfully deleted!!")

	return c.Redirect(http.StatusSeeOther, "/todo/list")
}

func (th *TaskHandler) logoutHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	// Revoke users authentication
	sess.Values = map[interface{}]interface{}{
		auth_key:     false,
		user_id_key:  "",
		username_key: "",
		tzone_key:    "",
	}
	sess.Save(c.Request(), c.Response())

	setFlashmessages(c, "success", "You have successfully logged out!!")

	fromProtected = false

	return c.Redirect(http.StatusSeeOther, "/login")
}
