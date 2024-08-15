package usecases

import (
	"context"
	//"errors"
	"net/http"
	"task_managment_api/domain"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
}

func NewTaskUsecase(taskRepository domain.TaskRepository, ) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
	}
}


func (uc *taskUsecase) GetTasks(c context.Context) ([]domain.Task, domain.CustomError) {
	return uc.taskRepository.GetTasks(c)
}

func (uc *taskUsecase) GetTaskByID(c context.Context, taskId string) (domain.Task, domain.CustomError) {
	return uc.taskRepository.GetTaskByID(c, taskId)
}


func (uc *taskUsecase) CreateTask(c context.Context, task domain.Task) domain.CustomError {
	if task.Title == "" {
		return domain.CustomError{ErrCode: http.StatusBadRequest, ErrMessage: "title is required"}
	}
	return uc.taskRepository.CreateTask(c, task)
}

func (uc *taskUsecase) UpdateTaskByID(c context.Context, taskId string, updatedTask domain.Task) domain.CustomError {
	updatedTask.ID = taskId
	return uc.taskRepository.UpdateTaskByID(c, updatedTask)
}

func (uc taskUsecase) DeleteTaskByID(c context.Context, taskId string) domain.CustomError {
	return uc.taskRepository.DeleteTaskByID(c, taskId)
}