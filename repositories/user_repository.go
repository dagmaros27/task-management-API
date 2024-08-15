package repositories

import (
	"context"
	"net/http"
	"task_managment_api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository instance.
func NewUserRepository(db *mongo.Database, userCollectionString string) domain.UserRepository {
	return &userRepository{
		collection: db.Collection(userCollectionString),
	}
}

// CreateUser inserts a new user into the database.
func (us *userRepository) CreateUser(c context.Context, user domain.User) domain.CustomError {
	_, err := us.collection.InsertOne(c, user)
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while creating user"}
	}
	return domain.CustomError{}
}

// GetUserByUsername retrieves a user from the database based on the username.
func (us *userRepository) GetUserByUsername(c context.Context, username string) (domain.User, domain.CustomError) {
	var user domain.User
	err := us.collection.FindOne(c, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, domain.CustomError{ErrCode: http.StatusBadRequest, ErrMessage: "User not found"}
		}
		return domain.User{}, domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while retrieving user"}
	}
	return user, domain.CustomError{}
}

// UpdateUser updates an existing user in the database.
func (us *userRepository) UpdateUser(c context.Context, user domain.User) domain.CustomError {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	user.ID = ""
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while updating user"}
	}
	updatedUser, err := us.collection.UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": user})
	if err != nil {
		return domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: err.Error()}
	}
	if updatedUser.MatchedCount == 0{
		return domain.CustomError{ErrCode: http.StatusNotFound, ErrMessage: "User not found"}
	}
	return domain.CustomError{}
}

// GetUserCount returns the total number of users in the database.
func (us *userRepository) GetUserCount(c context.Context) (int64, domain.CustomError) {
	count, err := us.collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return 0, domain.CustomError{ErrCode: http.StatusInternalServerError, ErrMessage: "Error while getting user count"}
	}
	return count, domain.CustomError{}
}
