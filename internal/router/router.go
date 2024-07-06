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
	}
	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", controllers.GetUsers)
		protected.POST("/users", controllers.CreateUser)
		protected.GET("/users/:id", controllers.GetUserByID)
		protected.PUT("/users/:id", controllers.UpdateUser)
		protected.DELETE("/users/:id", controllers.DeleteUser)

		protected.GET("/bill-mates", controllers.GetBillMates)
		protected.GET("/bill-mates/:id", controllers.GetBillMateById)
		protected.POST("/bill-mates", controllers.CreateBillMate)

		protected.GET("/bill-groups", controllers.GetBillGroups)
	}
	return router
}
