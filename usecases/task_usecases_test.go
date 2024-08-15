package usecases_test

import (
	"context"
	"net/http"
	"task_managment_api/domain"
	"task_managment_api/usecases"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetTasks(c context.Context) ([]domain.Task, domain.CustomError) {
	args := m.Called(c)
	return args.Get(0).([]domain.Task), args.Get(1).(domain.CustomError)
}

func (m *MockTaskRepository) GetTaskByID(c context.Context, taskId string) (domain.Task, domain.CustomError) {
	args := m.Called(c, taskId)
	return args.Get(0).(domain.Task), args.Get(1).(domain.CustomError)
}

func (m *MockTaskRepository) CreateTask(c context.Context, task domain.Task) domain.CustomError {
	args := m.Called(c, task)
	return args.Get(0).(domain.CustomError)
}

func (m *MockTaskRepository) UpdateTaskByID(c context.Context, updatedTask domain.Task) domain.CustomError {
	args := m.Called(c, updatedTask)
	return args.Get(0).(domain.CustomError)
}

func (m *MockTaskRepository) DeleteTaskByID(c context.Context, taskId string) domain.CustomError {
	args := m.Called(c, taskId)
	return args.Get(0).(domain.CustomError)
}

type TaskUsecaseSuite struct {
	suite.Suite
	mockRepo  *MockTaskRepository
	usecase   domain.TaskUsecase
}

func (suite *TaskUsecaseSuite) SetupTest() {
	suite.mockRepo = new(MockTaskRepository)
	suite.usecase = usecases.NewTaskUsecase(suite.mockRepo)
}

// Test GetTasks
func (suite *TaskUsecaseSuite) TestGetTasks() {
	mockTasks := []domain.Task{
		{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now().Format(time.RFC3339), Status: "Pending"},
	}

	suite.mockRepo.On("GetTasks", mock.Anything).Return(mockTasks, domain.CustomError{})

	tasks, err := suite.usecase.GetTasks(context.TODO())

	suite.Empty(err.ErrCode)
	suite.Equal(1, len(tasks))
	suite.Equal("Task 1", tasks[0].Title)

	suite.mockRepo.AssertExpectations(suite.T())
}

// Test GetTaskByID
func (suite *TaskUsecaseSuite) TestGetTaskByID() {
	mockTask := domain.Task{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now().Format(time.RFC3339), Status: "Pending"}

	suite.mockRepo.On("GetTaskByID", mock.Anything, "1").Return(mockTask, domain.CustomError{})

	task, err := suite.usecase.GetTaskByID(context.TODO(), "1")

	suite.Empty(err.ErrMessage)
	suite.Equal("Task 1", task.Title)

	suite.mockRepo.AssertExpectations(suite.T())
}

// Test CreateTask
func (suite *TaskUsecaseSuite) TestCreateTask() {
	mockTask := domain.Task{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now().Format(time.RFC3339), Status: "Pending"}

	suite.mockRepo.On("CreateTask", mock.Anything, mockTask).Return(domain.CustomError{})

	err := suite.usecase.CreateTask(context.TODO(), mockTask)

	suite.Empty(err.ErrMessage)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test CreateTask with Missing Title
func (suite *TaskUsecaseSuite) TestCreateTask_MissingTitle() {
	mockTask := domain.Task{ID: "1", Description: "First task", DueDate: time.Now().Format(time.RFC3339), Status: "Pending"}

	err := suite.usecase.CreateTask(context.TODO(), mockTask)

	suite.Equal(http.StatusBadRequest, err.ErrCode)
	suite.Equal("title is required", err.ErrMessage)
}

// Test UpdateTaskByID
func (suite *TaskUsecaseSuite) TestUpdateTaskByID() {
	mockTask := domain.Task{ID: "1", Title: "Updated Task", Description: "Updated Description", Status: "Completed"}

	suite.mockRepo.On("UpdateTaskByID", mock.Anything, mockTask).Return(domain.CustomError{})

	err := suite.usecase.UpdateTaskByID(context.TODO(), "1", mockTask)

	suite.Empty(err.ErrMessage)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test DeleteTaskByID
func (suite *TaskUsecaseSuite) TestDeleteTaskByID() {
	suite.mockRepo.On("DeleteTaskByID", mock.Anything, "1").Return(domain.CustomError{})

	err := suite.usecase.DeleteTaskByID(context.TODO(), "1")

	suite.Empty(err.ErrMessage)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}
