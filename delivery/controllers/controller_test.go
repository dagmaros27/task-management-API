package controllers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"task_managment_api/delivery/controllers"
	"task_managment_api/domain"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock for TaskUsecase
type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) GetTasks(c context.Context) ([]domain.Task, domain.CustomError) {
	args := m.Called(c)
	return args.Get(0).([]domain.Task), args.Get(1).(domain.CustomError)
}

func (m *MockTaskUsecase) GetTaskByID(c context.Context, id string) (domain.Task, domain.CustomError) {
	args := m.Called(c, id)
	return args.Get(0).(domain.Task), args.Get(1).(domain.CustomError)
}

func (m *MockTaskUsecase) UpdateTaskByID(c context.Context, id string, task domain.Task) domain.CustomError {
	args := m.Called(c, id, task)
	return args.Get(0).(domain.CustomError)
}

func (m *MockTaskUsecase) DeleteTaskByID(c context.Context, id string) domain.CustomError {
	args := m.Called(c, id)
	return args.Get(0).(domain.CustomError)
}

func (m *MockTaskUsecase) CreateTask(c context.Context, task domain.Task) domain.CustomError {
	args := m.Called(c, task)
	return args.Get(0).(domain.CustomError)
}

// TaskControllerTestSuite defines a suite of tests for the TaskController
type TaskControllerTestSuite struct {
	suite.Suite
	controller    *controllers.TaskController
	mockTaskUsecase *MockTaskUsecase
}

// SetupTest sets up the test environment before each test
func (suite *TaskControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockTaskUsecase = new(MockTaskUsecase)
	suite.controller = controllers.NewTaskController(suite.mockTaskUsecase)
}

func (suite *TaskControllerTestSuite) TearDownTest() {
	suite.mockTaskUsecase.AssertExpectations(suite.T())
}
// TestGetTasks tests the GetTasks method
func (suite *TaskControllerTestSuite) TestGetTasks() {
	mockTasks := []domain.Task{
		{ID: "1", Title: "Task 1", Description: "Description 1"},
		{ID: "2", Title: "Task 2", Description: "Description 2"},
	}

	suite.mockTaskUsecase.On("GetTasks", mock.Anything).Return(mockTasks, domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	suite.controller.GetTasks(c)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "Task 1")
}

// TestCreateTask tests the CreateTask method
func (suite *TaskControllerTestSuite) TestCreateTask() {
	taskJSON := `{"title": "New Task", "description": "New Description"}`

	suite.mockTaskUsecase.On("CreateTask", mock.Anything, mock.Anything).Return(domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.controller.CreateTask(c)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

//TestGetTaskById tests the GetTaskByID method
func (suite *TaskControllerTestSuite) TestGetTaskByID() {
	mockTask := domain.Task{ID: "1", Title: "Task 1", Description: "Description 1"}

	suite.mockTaskUsecase.On("GetTaskByID", mock.Anything, "1").Return(mockTask, domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	suite.controller.GetTaskByID(c)

	suite.Equal( http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "Task 1")
}

// TestUpdateTaskByID tests the UpdateTaskByID method
func (suite *TaskControllerTestSuite) TestUpdateTaskByID() {
	taskJSON := `{"title": "Updated Task", "description": "Updated Description"}`

	suite.mockTaskUsecase.On("UpdateTaskByID", mock.Anything, "1", mock.AnythingOfType("domain.Task")).Return(domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(taskJSON))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	suite.controller.UpdateTaskByID(c)
	suite.Equal(http.StatusOK, w.Code)
}

// TestDeleteTaskByID tests the DeleteTaskByID method
func (suite *TaskControllerTestSuite) TestDeleteTaskByID() {
	suite.mockTaskUsecase.On("DeleteTaskByID", mock.Anything, "1").Return(domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	suite.controller.DeleteTaskByID(c)
	suite.Equal(http.StatusOK, w.Code)
}



// TestTaskControllerTestSuite runs the suite of tests




// Mock for UserUsecase
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) RegisterUser(c context.Context, user domain.User) domain.CustomError {
	args := m.Called(c, user)
	return args.Get(0).(domain.CustomError)
}

func (m *MockUserUsecase) AuthenticateUser(c context.Context, username, password string) (string, domain.CustomError) {
	args := m.Called(c, username, password)
	return args.String(0), args.Get(1).(domain.CustomError)
}

func (m *MockUserUsecase) PromoteUser(c context.Context, username string) domain.CustomError {
	args := m.Called(c, username)
	return args.Get(0).(domain.CustomError)
}

// UserControllerTestSuite defines a suite of tests for the UserController
type UserControllerTestSuite struct {
	suite.Suite
	controller    *controllers.UserController
	mockUserUsecase *MockUserUsecase
}

// SetupTest sets up the test environment before each test
func (suite *UserControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockUserUsecase = new(MockUserUsecase)
	suite.controller = controllers.NewUserController(suite.mockUserUsecase)
}

func (suite *UserControllerTestSuite) TearDownTest() {
	suite.mockUserUsecase.AssertExpectations(suite.T())
}


// TestRegisterUser tests the RegisterUser method
func (suite *UserControllerTestSuite) TestRegisterUser() {
	userJSON := `{"username": "newuser", "password": "password"}`

	suite.mockUserUsecase.On("RegisterUser", mock.Anything, mock.Anything).Return(domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register", strings.NewReader(userJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.controller.RegisterUser(c)

	suite.Equal(http.StatusCreated, w.Code)
}

// TestLoginUser tests the LoginUser method
func (suite *UserControllerTestSuite) TestLoginUser() {
	loginJSON := `{"username": "user1", "password": "password"}`

	suite.mockUserUsecase.On("AuthenticateUser", mock.Anything, "user1", "password").Return("mocked_token", domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.controller.LoginUser(c)

	suite.Equal(http.StatusOK, w.Code)
	suite.JSONEq(`{"token": "mocked_token"}`, w.Body.String())
}

// TestPromoteUser tests the PromoteUser method
func (suite *UserControllerTestSuite) TestPromoteUser() {
	promoteJSON := `{"username": "user1"}`

	suite.mockUserUsecase.On("PromoteUser", mock.Anything, "user1").Return(domain.CustomError{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/promote", strings.NewReader(promoteJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.controller.PromoteUser(c)

	suite.Equal( http.StatusOK, w.Code)
	suite.JSONEq( `{"message": "User promoted successfully"}`, w.Body.String())
}

// TestControllerTestSuite runs the suites of the task tests and user tests
func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
	suite.Run(t, new(UserControllerTestSuite))
}
