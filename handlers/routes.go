package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ah *AuthHandler, th *TaskHandler) {
	e.GET("/", ah.flagsMiddleware(ah.homeHandler))
	e.GET("/login", ah.flagsMiddleware(ah.loginHandler))
	e.POST("/login", ah.flagsMiddleware(ah.loginHandler))
	e.GET("/register", ah.flagsMiddleware(ah.registerHandler))
	e.POST("/register", ah.flagsMiddleware(ah.registerHandler))

	protectedGroup := e.Group("/todo", ah.authMiddleware)
	/* ↓ Protected Routes ↓ */
	protectedGroup.GET("/list", th.todoListHandler)
	protectedGroup.GET("/create", th.createTodoHandler)
	protectedGroup.POST("/create", th.createTodoHandler)
	protectedGroup.GET("/edit/:id", th.updateTodoHandler)
	protectedGroup.POST("/edit/:id", th.updateTodoHandler)
	protectedGroup.DELETE("/delete/:id", th.deleteTodoHandler)
	protectedGroup.POST("/logout", th.logoutHandler)

	/* ↓ Fallback Page ↓ */
	e.GET("/*", RouteNotFoundHandler)
}
