package router

import (
	"github.com/brevd/equalizer/internal/controllers"
	"github.com/brevd/equalizer/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// Public routes
	public := router.Group("/")
	{
		public.GET("/", controllers.Index)
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}
	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", controllers.Logout)

		protected.GET("/users", controllers.GetUsers)
		protected.GET("/users/:id", controllers.GetUserByID)
		protected.PUT("/users/:id", controllers.UpdateUser)
		protected.DELETE("/users/:id", controllers.DeleteUser)

		protected.GET("/bill-mates", controllers.GetBillMates)
		protected.GET("/bill-mates/:id", controllers.GetBillMateById)
		protected.POST("/bill-mates", controllers.CreateBillMate)

		protected.GET("/bill-groups", controllers.GetBillGroups)

		protected.GET("/categories", controllers.GetCategories)
		protected.POST("/categories", controllers.CreateCategory)

		protected.POST("/expenses", controllers.CreateExpense)
		protected.GET("/expenses", controllers.GetExpenses)
		protected.GET("/expenses/:id", controllers.GetExpenseByID)

		protected.POST("/budgets", controllers.CreateBudget)
		protected.GET("/budgets/", controllers.GetBudgets)
	}
	return router
}
