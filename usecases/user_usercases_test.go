package usecases_test

import (
	"context"
	"task_managment_api/domain"
	"task_managment_api/usecases"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByUsername(c context.Context, username string) (domain.User, domain.CustomError) {
	args := m.Called(c, username)
	return args.Get(0).(domain.User), args.Get(1).(domain.CustomError)
}

func (m *MockUserRepository) GetUserCount(c context.Context) (int64, domain.CustomError) {
	args := m.Called(c)
	return args.Get(0).(int64), args.Get(1).(domain.CustomError)
}

func (m *MockUserRepository) CreateUser(c context.Context, user domain.User) domain.CustomError {
	args := m.Called(c, user)
	return args.Get(0).(domain.CustomError)
}

func (m *MockUserRepository) UpdateUser(c context.Context, user domain.User) domain.CustomError {
	args := m.Called(c, user)
	return args.Get(0).(domain.CustomError)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, domain.CustomError) {
	args := m.Called(password)
	return args.Get(0).(string), args.Get(1).(domain.CustomError)
}

func (m *MockPasswordService) VerifyPassword(user domain.User, password string) domain.CustomError {
	args := m.Called(user, password)
	return args.Get(0).(domain.CustomError)
}

type MockJWTService struct {
	mock.Mock
}	

func (m *MockJWTService) GenerateUserToken(user domain.User) (string, domain.CustomError) {
	args := m.Called(user)
	return args.Get(0).(string), args.Get(1).(domain.CustomError)
}

func (m *MockJWTService) ValidateToken(token string) (jwt.MapClaims, domain.CustomError) {
	args := m.Called(token)
	return args.Get(0).(jwt.MapClaims), args.Get(1).(domain.CustomError)
}

// Test Suite for UserUsecase
type UserUsecaseSuite struct {
	suite.Suite
	mockRepo        *MockUserRepository
	mockPasswordSvc *MockPasswordService
	mockJwtService *MockJWTService
	usecase         domain.UserUsecase
}

func (suite *UserUsecaseSuite) SetupTest() {
	suite.mockRepo = new(MockUserRepository)
	suite.mockPasswordSvc = new(MockPasswordService)
	suite.mockJwtService = new(MockJWTService)
	suite.usecase = usecases.NewUserUsecase(suite.mockRepo, suite.mockJwtService,suite.mockPasswordSvc )
}

func (suite *UserUsecaseSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockPasswordSvc.AssertExpectations(suite.T())
}

// Test RegisterUser
func (suite *UserUsecaseSuite) TestRegisterUser() {
	user := domain.User{Username: "testuser", Password: "password"}

	suite.mockRepo.On("GetUserByUsername", mock.Anything, user.Username).Return(domain.User{}, domain.CustomError{ErrCode: 400, ErrMessage: "User not found"})
	suite.mockRepo.On("GetUserCount", mock.Anything).Return(int64(1), domain.CustomError{})
	suite.mockPasswordSvc.On("HashPassword", user.Password).Return("hashedpassword", domain.CustomError{})
	suite.mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(domain.CustomError{})

	err := suite.usecase.RegisterUser(context.TODO(), user)

	suite.Empty(err.ErrCode)
	suite.mockRepo.AssertCalled(suite.T(), "GetUserByUsername", mock.Anything, user.Username)
	suite.mockRepo.AssertCalled(suite.T(), "GetUserCount", mock.Anything)
	suite.mockPasswordSvc.AssertCalled(suite.T(), "HashPassword", user.Password)
	suite.mockRepo.AssertCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("domain.User"))
}

// Test RegisterUser when User already exists
func (suite *UserUsecaseSuite) TestRegisterUser_UserAlreadyExists() {
	user := domain.User{Username: "testuser", Password: "password"}

	suite.mockRepo.On("GetUserByUsername", mock.Anything, user.Username).Return(user, domain.CustomError{})

	err := suite.usecase.RegisterUser(context.TODO(), user)

	suite.Equal(409, err.ErrCode)
	suite.Equal("User already exists", err.ErrMessage)
	suite.mockRepo.AssertCalled(suite.T(), "GetUserByUsername", mock.Anything, user.Username)
	suite.mockRepo.AssertNotCalled(suite.T(), "GetUserCount", mock.Anything)
	suite.mockPasswordSvc.AssertNotCalled(suite.T(), "HashPassword", user.Password)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("domain.User"))
}

// Test RegisterUser when hashing password fails
func (suite *UserUsecaseSuite) TestRegisterUser_HashPasswordError() {
	user := domain.User{Username: "testuser", Password: "password"}
	suite.mockRepo.On("GetUserByUsername", mock.Anything, user.Username).Return(domain.User{}, domain.CustomError{ErrCode: 400, ErrMessage: "User not found"})
	suite.mockRepo.On("GetUserCount", mock.Anything).Return(int64(0), domain.CustomError{})
	suite.mockPasswordSvc.On("HashPassword", user.Password).Return("", domain.CustomError{ErrCode: 500, ErrMessage: "Error while hashing password"})

	err := suite.usecase.RegisterUser(context.TODO(), user)
	suite.Equal(500, err.ErrCode)

	suite.mockRepo.AssertCalled(suite.T(), "GetUserByUsername", mock.Anything, user.Username)
	suite.mockRepo.AssertCalled(suite.T(), "GetUserCount", mock.Anything)
	suite.mockPasswordSvc.AssertCalled(suite.T(), "HashPassword", user.Password)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("domain.User"))
}

// Test AuthenticateUser
func (suite *UserUsecaseSuite) TestAuthenticateUser() {
	user := domain.User{Username: "testuser", Password: "hashedpassword"}

	suite.mockRepo.On("GetUserByUsername", mock.Anything, user.Username).Return(user, domain.CustomError{})
	suite.mockPasswordSvc.On("VerifyPassword", user, "password").Return(domain.CustomError{})
	suite.mockJwtService.On("GenerateUserToken", user).Return("token", domain.CustomError{})

	token, err := suite.usecase.AuthenticateUser(context.TODO(), user.Username, "password")

	suite.Empty(err.ErrMessage)
	suite.NotEmpty(token)
}

// Test AuthenticateUser with Invalid Credentials
func (suite *UserUsecaseSuite) TestAuthenticateUser_InvalidCredentials() {
	user := domain.User{Username: "testuser"}
	suite.mockRepo.On("GetUserByUsername", mock.Anything, user.Username).Return(domain.User{}, domain.CustomError{ErrCode: 404, ErrMessage: "User not found"})

	token, err := suite.usecase.AuthenticateUser(context.TODO(), user.Username, "password")

	suite.Equal("", token)
	suite.Equal(401, err.ErrCode)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test PromoteUser
func (suite *UserUsecaseSuite) TestPromoteUser() {
	user := domain.User{Username: "testuser", Role: "user"}

	suite.mockRepo.On("GetUserByUsername", mock.Anything, user.Username).Return(user, domain.CustomError{})
	suite.mockRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(domain.CustomError{})

	err := suite.usecase.PromoteUser(context.TODO(), user.Username)

	suite.Empty(err.ErrMessage)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Run the test suite
func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}
