package router

import (
	"task_managment_api/delivery/controllers"
	"task_managment_api/infrastructure"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database, taskController *controllers.TaskController, userController *controllers.UserController, authService infrastructure.AuthMiddlewareService) *gin.Engine {

	
	router := gin.Default()

	// public routes
	router.POST("/register", userController.RegisterUser)
	router.POST("/login", userController.LoginUser)



	// private routes
	authorized := router.Group("/")
	authorized.Use(authService.AuthMiddleware())

	// task routes
	authorized.GET("/tasks", taskController.GetTasks)
	authorized.GET("/tasks/:id", taskController.GetTaskByID)
	authorized.POST("/tasks", authService.AdminMiddleware(), taskController.CreateTask)
	authorized.PUT("/tasks/:id", authService.AdminMiddleware(), taskController.UpdateTaskByID)
	authorized.DELETE("/tasks/:id", authService.AdminMiddleware(), taskController.DeleteTaskByID)

	// user promotion route
	authorized.POST("/promote", authService.AdminMiddleware(), userController.PromoteUser)

	return router
}
