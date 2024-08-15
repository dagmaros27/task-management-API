package router

import (
	"task_managment_api/delivery/controllers"
	"task_managment_api/infrastructure"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database, taskController *controllers.TaskController, userController *controllers.UserController) *gin.Engine {

	
	router := gin.Default()

	// public routes
	router.POST("/register", userController.RegisterUser)
	router.POST("/login", userController.LoginUser)

	// private routes
	authorized := router.Group("/")
	authorized.Use(infrastructure.AuthMiddleware())

	// task routes
	authorized.GET("/tasks", taskController.GetTasks)
	authorized.GET("/tasks/:id", taskController.GetTaskByID)
	authorized.POST("/tasks", infrastructure.AdminMiddleware(), taskController.CreateTask)
	authorized.PUT("/tasks/:id", infrastructure.AdminMiddleware(), taskController.UpdateTaskByID)
	authorized.DELETE("/tasks/:id", infrastructure.AdminMiddleware(), taskController.DeleteTaskByID)

	// user promotion route
	authorized.POST("/promote", infrastructure.AdminMiddleware(), userController.PromoteUser)

	return router
}
