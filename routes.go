package main

import (
	"expence_management/api"

	"github.com/gin-gonic/gin"
)

func makeRouts(h *api.Handler, allowedOrigins []string) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware(allowedOrigins))

	router.Static("/assets", "./dist/assets")
	router.StaticFile("/", "./dist/index.html")

	router.StaticFile("/incomes", "./dist/index.html")
	router.StaticFile("/expenses", "./dist/index.html")
	router.StaticFile("/users", "./dist/index.html")
	router.StaticFile("/report", "./dist/index.html")
	router.StaticFile("/login", "./dist/index.html")
	router.StaticFile("/logout", "./dist/index.html")
	router.StaticFile("/signup", "./dist/index.html")
	router.StaticFile("/back-login.avif", "./dist/back-login.avif")
	router.StaticFile("/back.png", "./dist/back.png")
	router.StaticFile("/money.ico", "./dist/money.ico")
	router.StaticFile("/userinfo", "./dist/index.html")
	router.StaticFile("/upcoming", "./dist/index.html")

	router.POST("/api/v1/signup", h.Signup)
	router.POST("/api/v1/login", h.Login)
	router.GET("/api/v1/logout", api.CheckAuth, h.Logout)
	router.POST("/api/v1/change_password", api.CheckAuth, h.ChangePassword)
	router.GET("/api/v1/userinfo", api.CheckAuth, h.UserInfo)

	router.GET("/api/v1/users/list", api.CheckAuth, api.IsAdmin, h.ListUsers)
	router.POST("/api/v1/users/add", api.CheckAuth, api.IsAdmin, h.AddNewUser)
	router.POST("/api/v1/users/update", api.CheckAuth, api.IsAdmin, h.UpdateUser)
	router.DELETE("/api/v1/users/delete", api.CheckAuth, api.IsAdmin, h.DeleteUser)
	router.POST("/api/v1/users/change_password", api.CheckAuth, api.IsAdmin, h.ChangePasswordByAdmin)

	router.POST("/api/v1/incomes/add", api.CheckAuth, h.AddIncome)
	router.POST("/api/v1/incomes/update", api.CheckAuth, h.UpdateIncome)
	router.DELETE("/api/v1/incomes/delete", api.CheckAuth, h.DeleteIncome)
	router.GET("/api/v1/incomes/list", api.CheckAuth, h.Incomes)

	router.POST("/api/v1/expenses/add", api.CheckAuth, h.AddExpense)
	router.POST("/api/v1/expenses/update", api.CheckAuth, h.UpdateExpense)
	router.DELETE("/api/v1/expenses/delete", api.CheckAuth, h.DeleteExpense)
	router.GET("/api/v1/expenses/list", api.CheckAuth, h.Expenses)

	router.GET("/api/v1/incomes/report", api.CheckAuth, h.ReportIncomes)
	router.GET("/api/v1/expenses/report", api.CheckAuth, h.ReportExpenses)

	router.GET("/api/v1/incomes/export", api.CheckAuth, h.ExportIncomes)
	router.GET("/api/v1/expenses/export", api.CheckAuth, h.ExportExpenses)

	router.GET("/api/v1/incomes/summary", api.CheckAuth, h.IncomeSummary)
	router.GET("/api/v1/expenses/summary", api.CheckAuth, h.ExpenseSummary)

	router.GET("/api/v1/summary", api.CheckAuth, h.OverallSummary)

	router.POST("/api/v1/categories/add", api.CheckAuth, h.AddCategory)
	router.POST("/api/v1/categories/update", api.CheckAuth, h.UpdateCategory)
	router.DELETE("/api/v1/categories/delete", api.CheckAuth, h.DeleteCategory)
	router.GET("/api/v1/categories/list", api.CheckAuth, h.Categories)

	router.POST("/api/v1/upcoming_expenses/add", api.CheckAuth, h.AddUpcomingExpense)
	router.POST("/api/v1/upcoming_expenses/update", api.CheckAuth, h.UpdateUpcomingExpense)
	router.DELETE("/api/v1/upcoming_expenses/delete", api.CheckAuth, h.DeleteUpcomingExpense)
	router.GET("/api/v1/upcoming_expenses/list", api.CheckAuth, h.UpcomingExpenses)
	router.PUT("/api/v1/upcoming_expenses/pay", api.CheckAuth, h.PayUpcomingExpense)

	return router
}
