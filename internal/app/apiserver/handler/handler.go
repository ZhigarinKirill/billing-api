package handler

import (
	"github.com/ZhigarinKirill/billing-api/internal/app/service"
	"github.com/gin-gonic/gin"
)

// Handler - handler structure
type Handler struct {
	Services *service.Service
}

// NewHandler - handler constructor
func NewHandler(services *service.Service) *Handler {
	return &Handler{Services: services}
}

// InitRoutes - endpoints handlers initialization function
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// User account routes
	users := router.Group("/users/:user_id")
	{
		account := users.Group("/account")
		{
			account.GET("", h.getUserBillAmount)
			account.POST("/deposit", h.changeUserBillAmount)
			account.POST("/withdraw", h.changeUserBillAmount)
			account.POST("/transfer", h.transferMoneyBetweenUsers)
		}
	}

	transactions := router.Group("/transactions")
	{
		transactions.POST("/reserve", h.reserveMoneyFromUser)
		transactions.POST("/complete", h.completeReservationTransaction)
	}

	router.POST("/users/auth/sign-up", h.createUser) // New user registration
	router.GET("/users", h.getAllUsers)
	router.GET("/month_report", h.getReport)
	return router
}
