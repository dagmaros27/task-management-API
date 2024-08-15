package controllers

import (
	"net/http"
	"task_managment_api/domain"

	//"time"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskUsecase domain.TaskUsecase
}

type UserController struct {
	userUsecase domain.UserUsecase
}

//task controllers

func NewTaskController(taskUsecase domain.TaskUsecase) *TaskController {
	return &TaskController{
		taskUsecase: taskUsecase,
	}
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.taskUsecase.GetTasks(c)
	if err.ErrCode != 0  {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}
	if len(tasks) == 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.taskUsecase.GetTaskByID(c, id)
	if err.ErrCode != 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}
	err := tc.taskUsecase.UpdateTaskByID(c,id, task)
	if err.ErrCode != 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (tc *TaskController) DeleteTaskByID(c *gin.Context) {
	id := c.Param("id")
	err := tc.taskUsecase.DeleteTaskByID(c,id)
	if err.ErrCode != 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}
	
	err := tc.taskUsecase.CreateTask(c,task)
	if err.ErrCode != 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

//user controllers

func NewUserController(userUsecase domain.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var user domain.User
	
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	err := uc.userUsecase.RegisterUser(c,user)
	
	// TODO: should return statusConflict if err is user already created
	if err.ErrCode != 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid JSON"})
		return
	}

	token, err := uc.userUsecase.AuthenticateUser(c,user.Username, user.Password)
	if err.ErrCode != 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}



func (uc *UserController) PromoteUser(c *gin.Context) {
	
	var user domain.UserToPromote

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	err := uc.userUsecase.PromoteUser(c,user.Username)
	if err.ErrCode != 0 {
		c.JSON(err.ErrCode, gin.H{"message": err.ErrMessage})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted successfully"})
}
