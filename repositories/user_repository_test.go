package repositories_test

import (
	"context"
	"net/http"
	"task_managment_api/domain"
	"task_managment_api/repositories"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type UserRepositorySuite struct {
	suite.Suite
	db         *mongo.Database
	collection *mongo.Collection
	repo       domain.UserRepository
}


func (suite *UserRepositorySuite) SetupTest() {
	// Clear the collection before each test
	suite.collection.DeleteMany(context.TODO(), bson.D{})
}

func (suite *UserRepositorySuite) SetupSuite() {
	// Set up a test MongoDB instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	suite.Require().NoError(err)

	suite.db = client.Database("task_management_test")
	suite.collection = suite.db.Collection("users")

	suite.repo = repositories.NewUserRepository(suite.db, "users")
}

func (suite *UserRepositorySuite) TearDownSuite() {
	// Drop the test database
	//suite.Require().NoError(suite.db.Drop(context.TODO()))
}

// Test CreateUser
func (suite *UserRepositorySuite) TestCreateUser() {
	user := domain.User{
		Username: "Test User",
		Password: "hashed password",
		Role: "admin",}

	err := suite.repo.CreateUser(context.TODO(), user)
	suite.Empty(err.ErrCode)

	var result domain.User
	dbError := suite.collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&result)
	suite.NoError(dbError)
	suite.Equal(user.Username, result.Username)
}


//Test GetUserByUsername
func (suite *UserRepositorySuite) TestGetUserByUsername(){
	user := domain.User{
		Username: "Test User",
		Password: "hashed password",
		Role: "admin",}

	_, dbError := suite.collection.InsertOne(context.TODO(), user)
	suite.NoError(dbError)

	result, err := suite.repo.GetUserByUsername(context.TODO(), user.Username)
	suite.Empty(err.ErrCode)
	suite.Equal(user.Username, result.Username)
}

//Test GetUserByUsername_NotFound
func (suite *UserRepositorySuite) TestGetUserByUsername_NotFound(){
	result, err := suite.repo.GetUserByUsername(context.TODO(), "invalid username")
	suite.Empty(result.Username)
	suite.Equal(http.StatusBadRequest, err.ErrCode)
}

//Test GetUserByUsername_Error


//Test update user
func (suite *UserRepositorySuite) TestUpdateUser(){
	user := domain.User{
		Username: "Test User",
		Password: "hashed password",
		Role: "admin",}

	insertedResult, dbError := suite.collection.InsertOne(context.TODO(), user)
	suite.NoError(dbError)
	user.ID = insertedResult.InsertedID.(primitive.ObjectID).Hex()
	user.Role = "user"

	err := suite.repo.UpdateUser(context.TODO(), user)
	suite.Empty(err.ErrMessage)

	var result domain.User
	dbError = suite.collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&result)
	suite.NoError(dbError)
	suite.Equal(user.Role, result.Role)
}

//Test UpdateUser_Error
func (suite *UserRepositorySuite) TestUpdateUser_Error(){
	user := domain.User{
		Username: "Test User",
		Password: "hashed password",
		Role: "admin",}

	insertedResult, dbError := suite.collection.InsertOne(context.TODO(), user)
	suite.NoError(dbError)
	user.ID = insertedResult.InsertedID.(primitive.ObjectID).Hex() + "invalid"
	user.Role = "user"

	err := suite.repo.UpdateUser(context.TODO(), user)
	suite.Equal(http.StatusInternalServerError, err.ErrCode)
}

//Test UpdateUser_NotFound
func (suite *UserRepositorySuite) TestUpdateUser_NotFound(){
	user := domain.User{
		ID: primitive.NewObjectID().Hex(),
		Username: "Test User",
		Password: "hashed password",
		Role: "admin",}
	err := suite.repo.UpdateUser(context.TODO(), user)
	suite.Equal(http.StatusNotFound, err.ErrCode)
}
	
	


//Test get user by count
func (suite *UserRepositorySuite) TestGetUserByCount(){
	user := domain.User{
		Username: "Test User",
		Password: "hashed password",
		Role: "admin",}

	_, dbError := suite.collection.InsertOne(context.TODO(), user)
	suite.NoError(dbError)

	count, err := suite.collection.CountDocuments(context.TODO(), bson.M{"username": user.Username})
	suite.NoError(err)
	suite.Equal(int64(1), count)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}