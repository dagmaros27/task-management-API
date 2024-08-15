package main

import (
	"context"
	"log"
	bootstrap "task_managment_api"
	"task_managment_api/delivery/controllers"
	"task_managment_api/delivery/router"
	"task_managment_api/infrastructure"
	"task_managment_api/repositories"
	"task_managment_api/usecases"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Application struct {
	Db  *mongo.Database
	Env *bootstrap.Env
}

func App() Application {
	app := &Application{}
	app.Env = bootstrap.NewEnv()
	app.Db = NewMongoDatabase(app.Env)
	return *app
}

//initialize a new database connection and return the database instance
func NewMongoDatabase(env *bootstrap.Env) *mongo.Database {
	clientOptions := options.Client().ApplyURI(env.DbUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(env.DbName)

	err = EnsureIndexes(db, env.DbUserCollection)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

//make sure username is unique in database level
func EnsureIndexes(db *mongo.Database, userCollectionString string) error {
	userCollection := db.Collection(userCollectionString)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCollection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}


func main() {

	app := App()

	tr := repositories.NewTaskRepository(app.Db, app.Env.DbTaskCollection)
	tc := repositories.NewUserRepository(app.Db, app.Env.DbUserCollection)
	ps := infrastructure.NewPasswordService()

	js := infrastructure.NewJWTService(app.Env.AccessTokenSecret)	
	as := infrastructure.NewAuthService(js)
	taskController := controllers.NewTaskController(usecases.NewTaskUsecase(tr)) 
	userController := controllers.NewUserController(usecases.NewUserUsecase(tc, js,ps))


	r := router.SetupRouter(app.Db, taskController, userController,as )
	r.Run(":8080")	
}
