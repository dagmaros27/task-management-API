package domain

import (
	"context"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId string `json:"userId"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}


type Task struct {
	ID          string `json:"_id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	DueDate     string `json:"due_date" bson:"due_date"`
	Status      string `json:"status" bson:"status"`
}

type User struct {
	ID       string `json:"_id " bson:"_id,omitempty"`
	Username string `json:"username" binding:"required" bson:"username"`
	Password string `json:"password" binding:"required" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

type UserToPromote struct {
	Username string `json:"username" binding:"required"`
}

type CustomError struct{
	ErrCode int
	ErrMessage string
}


type TaskRepository interface {
	GetTasks(c context.Context) ([]Task, CustomError)
	GetTaskByID(c context.Context, taskID string) (Task, CustomError)
	CreateTask(c context.Context, task Task) CustomError
	UpdateTaskByID(c context.Context, updatedTask Task) CustomError
	DeleteTaskByID(c context.Context, taskID string) CustomError
}


type TaskUsecase interface {
	GetTasks(c context.Context) ([]Task, CustomError)
	GetTaskByID(c context.Context, taskID string) (Task, CustomError)
	CreateTask(c context.Context, task Task) CustomError
	UpdateTaskByID(c context.Context, taskID string, updatedTask Task) CustomError
	DeleteTaskByID(c context.Context, taskID string) CustomError
}

type UserRepository interface {
	CreateUser(c context.Context, user User) CustomError
	GetUserByUsername(c context.Context, username string) (User, CustomError)
	UpdateUser(c context.Context, user User) CustomError
	GetUserCount(c context.Context)(int64,CustomError)
}

type UserUsecase interface {
	RegisterUser(c context.Context, user User) CustomError
	AuthenticateUser(c context.Context, username string, password string) (string, CustomError)
	PromoteUser(c context.Context, username string) CustomError
}



