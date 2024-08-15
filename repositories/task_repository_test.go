package repositories_test

import (
	"context"
	"task_managment_api/domain"
	"task_managment_api/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepositorySuite struct {
	suite.Suite
	db         *mongo.Database
	collection *mongo.Collection
	repo       domain.TaskRepository
}

func (suite *TaskRepositorySuite) SetupTest() {
	// Clear the collection before each test
	suite.collection.DeleteMany(context.TODO(), bson.D{})
}

func (suite *TaskRepositorySuite) SetupSuite() {
	// Set up a test MongoDB instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	suite.Require().NoError(err)

	suite.db = client.Database("task_management_test")
	suite.collection = suite.db.Collection("tasks")

	suite.repo = repositories.NewTaskRepository(suite.db, "tasks")
}

func (suite *TaskRepositorySuite) TearDownSuite() {
	// Drop the test database
	//suite.Require().NoError(suite.db.Drop(context.TODO()))
}

// Test CreateTask
func (suite *TaskRepositorySuite) TestCreateTask() {
	task := domain.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		DueDate:     time.Now().Format(time.RFC3339),
		Status:      "Pending",
	}

	err := suite.repo.CreateTask(context.TODO(), task)
	suite.Empty(err.ErrCode)

	var result domain.Task
	dbError := suite.collection.FindOne(context.TODO(), bson.M{"title": task.Title}).Decode(&result)
	suite.NoError(dbError)
	suite.Equal(task.Title, result.Title)
}

// Test GetTasks
func (suite *TaskRepositorySuite) TestGetTasks() {
	_, dbError := suite.collection.InsertOne(context.TODO(), domain.Task{
		Title:       "Sample Task",
		Description: "Sample Description",
		DueDate:     time.Now().Format(time.RFC3339),
		Status:      "Pending",
	})
	suite.NoError(dbError)

	tasks, err := suite.repo.GetTasks(context.TODO())
	suite.Empty(err.ErrCode)
	suite.NotEmpty(tasks)
}

// Test GetTaskByID
func (suite *TaskRepositorySuite) TestGetTaskByID() {
	task := domain.Task{
		Title:       "GetTaskByID Test",
		Description: "This task is for testing GetTaskByID",
		DueDate:     time.Now().Format(time.RFC3339),
		Status:      "Pending",
	}

	insertedResult, dbError := suite.collection.InsertOne(context.TODO(), task)
	suite.NoError(dbError)

	result, err := suite.repo.GetTaskByID(context.TODO(), insertedResult.InsertedID.(primitive.ObjectID).Hex())
	suite.Empty(err.ErrCode)
	suite.Equal(task.Title, result.Title)
}

// Test UpdateTaskByID
func (suite *TaskRepositorySuite) TestUpdateTaskByID() {
	task := domain.Task{
		Title:       "Update Task",
		Description: "This is a task to update",
		DueDate:     time.Now().Format(time.RFC3339),
		Status:      "Pending",
	}

	insertedResult, dbError := suite.collection.InsertOne(context.TODO(), task)
	suite.NoError(dbError)

	updatedTask := domain.Task{
		ID:          insertedResult.InsertedID.(primitive.ObjectID).Hex(),
		Title:       "Updated Task",
		Description: "This task has been updated",
		Status:      "Completed",
	}

	err := suite.repo.UpdateTaskByID(context.TODO(), updatedTask)
	suite.Empty(err.ErrCode)

	var result domain.Task
	dbError = suite.collection.FindOne(context.TODO(), bson.M{"_id": insertedResult.InsertedID}).Decode(&result)
	suite.NoError(dbError)
	suite.Equal(updatedTask.Title, result.Title)
	suite.Equal(updatedTask.Status, result.Status)
}

// Test DeleteTaskByID
func (suite *TaskRepositorySuite) TestDeleteTaskByID() {
	task := domain.Task{
		Title:       "Delete Task",
		Description: "This task will be deleted",
		DueDate:     time.Now().Format(time.RFC3339),
		Status:      "Pending",
	}

	insertedResult, dbError := suite.collection.InsertOne(context.TODO(), task)
	suite.NoError(dbError)

	err := suite.repo.DeleteTaskByID(context.TODO(), insertedResult.InsertedID.(primitive.ObjectID).Hex())
	suite.Empty(err.ErrCode)

	dbError = suite.collection.FindOne(context.TODO(), bson.M{"_id": task.ID}).Err()
	suite.Equal(mongo.ErrNoDocuments, dbError)
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}
