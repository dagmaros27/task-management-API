package bootstrap

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	DbUri                 string `mapstructure:"DB_URI"`
	DbName                 string `mapstructure:"DB_NAME"`
	DbTaskCollection                 string `mapstructure:"DB_TASK_COLLECTION"`
	DbUserCollection                 string `mapstructure:"DB_USER_COLLECTION"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
}

func NewEnv() *Env {
	env := Env{}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Can't find the path of .env : ", err)
	}


	viper.SetConfigFile(dir + "/../.env")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}