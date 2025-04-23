package main

import (
	todo "helloapp"
	"helloapp/pkg/handler"
	"helloapp/pkg/repository"
	"helloapp/pkg/service"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("произошла ошибка при инициализации конфигураций: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("ошибка загрузки env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("не удалось инициализировать базу данных: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	emailSender := service.NewMockEmailSender()
	services := service.NewService(repos, emailSender) // Добавляем emailSender
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("произошла ошибка при запуске http-сервера: %s", err.Error())
	}
	// srv := new(todo.Server)
	// if err := srv.Run("0.0.0.0:8000", handlers.InitRoutes()); err != nil {
	//     log.Fatalf("произошла ошибка при запуске http-сервера: %s", err.Error())
	// }
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
