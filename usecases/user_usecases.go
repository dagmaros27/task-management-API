package usecases

import (
	"context"
	"fmt"
	"net/http"

	"task_managment_api/domain"
	"task_managment_api/infrastructure"
	"time"
)

type userUsecase struct {
	userRepository domain.UserRepository
	ctxTimeout time.Duration
	passwordService infrastructure.PasswordService
}

func NewUserUsecase(userRepository domain.UserRepository, ctxTimeout time.Duration, passwordService infrastructure.PasswordService) domain.UserUsecase {
	return &userUsecase{userRepository: userRepository, ctxTimeout: ctxTimeout, passwordService: passwordService}
}


func (uc *userUsecase)RegisterUser(c context.Context, user domain.User) domain.CustomError{
	
	_ ,err := uc.userRepository.GetUserByUsername(c, user.Username)

	if err.ErrCode == 0{
		return domain.CustomError{ErrCode: http.StatusConflict, ErrMessage: "User already exists"}
	}
	if err.ErrMessage !=  "User not found" {
		fmt.Println(err.ErrMessage)
		return domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while checking user existence"}
	}


	count, err := uc.userRepository.GetUserCount(c)

	if err.ErrCode != 0 {
		return err
	}

	if count == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}	
	hashed,err :=  uc.passwordService.HashPassword(user.Password)
	if err.ErrCode != 0 {
		return err
	}
	
	user.Password = hashed
	
	return uc.userRepository.CreateUser(c, user)
}


func (uc *userUsecase)AuthenticateUser(c context.Context, username, password string) (string, domain.CustomError){
	
	user, err := uc.userRepository.GetUserByUsername(c, username)

	if err.ErrCode != 0 {
		if err.ErrCode ==  500{
			return "", domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while checking user"}

		}
		return "", domain.CustomError{ErrCode: http.StatusUnauthorized, ErrMessage: "Invalid username or password"}
	}

	err = uc.passwordService.VerifyPassword(user, password)

	if err.ErrCode != 0 { 
		return "", err
	}

	return infrastructure.NewJWTService().GenerateUserToken(user)
}


func (uc *userUsecase)PromoteUser(c context.Context, username string) domain.CustomError{
	user, err := uc.userRepository.GetUserByUsername(c, username)
	if err.ErrCode != 0 {
		return err
	}
	user.Role = "admin"
	return uc.userRepository.UpdateUser(c, user)
}

