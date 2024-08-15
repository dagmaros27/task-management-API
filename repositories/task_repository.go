package repositories

import (
	"context"
	//"errors"
	"net/http"
	"task_managment_api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database, taskCollectionString string) domain.TaskRepository {
	return &taskRepository{
		collection: db.Collection(taskCollectionString),
	}
}

// GetTasks retrieves all tasks from the database.
func (ts *taskRepository) GetTasks(c context.Context) ([]domain.Task, domain.CustomError) {
	var tasks []domain.Task
	cursor, err := ts.collection.Find(c, bson.D{})
	if err != nil {
		return nil, domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: err.Error()}
	}

	if err := cursor.All(c, &tasks); err != nil {
		return nil, domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: err.Error()}
	}
	return tasks, domain.CustomError{}
}

// GetTaskByID retrieves a task from the database by its ID.
func (ts *taskRepository) GetTaskByID(c context.Context, taskID string) (domain.Task, domain.CustomError) {
	var task domain.Task
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return domain.Task{}, domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Invalid task ID"}
	}

	err = ts.collection.FindOne(c, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, domain.CustomError{ErrCode: http.StatusNotFound, ErrMessage: "Task not found"}
		}
		return domain.Task{}, domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while retriving task"}
	}

	return task, domain.CustomError{}
}

// CreateTask creates a new task in the database.
func (ts *taskRepository) CreateTask(c context.Context, task domain.Task) domain.CustomError {
	_, err := ts.collection.InsertOne(c, task)
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while creating task"}
	}
	return domain.CustomError{}
}

// UpdateTaskByID updates a task in the database by its ID.
func (ts *taskRepository) UpdateTaskByID(c context.Context, updatedTask domain.Task) domain.CustomError {
	objectID, err := primitive.ObjectIDFromHex(updatedTask.ID)
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusBadRequest, ErrMessage: "Invalid task ID"}
	}

	update := bson.M{}

	if updatedTask.Title != "" {
		update["title"] = updatedTask.Title
	}
	if updatedTask.Description != "" {
		update["description"] = updatedTask.Description
	}
	if updatedTask.DueDate != "" {
		update["due_date"] = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		update["status"] = updatedTask.Status
	}

	result, err := ts.collection.UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while updating task"}
	}
	if result.MatchedCount == 0 {
		return domain.CustomError{ErrCode: http.StatusNotFound, ErrMessage: "Task not found"}
	}
	return domain.CustomError{}
}

// DeleteTaskByID deletes a task from the database by its ID.
func (ts *taskRepository) DeleteTaskByID(c context.Context, taskID string) domain.CustomError {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusBadRequest, ErrMessage: "Invalid task id"}
	}

	result, err := ts.collection.DeleteOne(c, bson.M{"_id": objectID})
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while deleting task"}
	}

	if result.DeletedCount == 0 {
		return domain.CustomError{ErrCode: http.StatusNotFound, ErrMessage: "Task not found"}
	}
	return domain.CustomError{}
}
