package httpHandler

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/config/postgres"

	crudHandler "github.com/novalwardhana/golang-boilerplate/module/crud/handler"
	crudRepository "github.com/novalwardhana/golang-boilerplate/module/crud/repository"
	crudUsecase "github.com/novalwardhana/golang-boilerplate/module/crud/usecase"

	userAuthenticationHandler "github.com/novalwardhana/golang-boilerplate/module/user-authentication/handler"
	userAuthenticationRepository "github.com/novalwardhana/golang-boilerplate/module/user-authentication/repository"
	userAuthenticationUsecase "github.com/novalwardhana/golang-boilerplate/module/user-authentication/usecase"

	userManagementHandler "github.com/novalwardhana/golang-boilerplate/module/user-management/handler"
	userManagementRepository "github.com/novalwardhana/golang-boilerplate/module/user-management/repository"
	userManagementUsecase "github.com/novalwardhana/golang-boilerplate/module/user-management/usecase"

	fileHandler "github.com/novalwardhana/golang-boilerplate/module/file/handler"
	fileRepository "github.com/novalwardhana/golang-boilerplate/module/file/repository"
	fileUsecase "github.com/novalwardhana/golang-boilerplate/module/file/usecase"

	advanceCrudHandler "github.com/novalwardhana/golang-boilerplate/module/advance-crud/handler"
	advanceCrudRepository "github.com/novalwardhana/golang-boilerplate/module/advance-crud/repository"
	advanceCrudUsecase "github.com/novalwardhana/golang-boilerplate/module/advance-crud/usecase"

	httpClientHandler "github.com/novalwardhana/golang-boilerplate/module/http-client/handler"
	httpClientRepository "github.com/novalwardhana/golang-boilerplate/module/http-client/repository"
	httpClientUsecase "github.com/novalwardhana/golang-boilerplate/module/http-client/usecase"

	emailHandler "github.com/novalwardhana/golang-boilerplate/module/email/handler"
	emailRepository "github.com/novalwardhana/golang-boilerplate/module/email/repository"
	emailUsecase "github.com/novalwardhana/golang-boilerplate/module/email/usecase"

	sftpHandler "github.com/novalwardhana/golang-boilerplate/module/sftp/handler"
	sftpRepository "github.com/novalwardhana/golang-boilerplate/module/sftp/repository"
	sftpUsecase "github.com/novalwardhana/golang-boilerplate/module/sftp/usecase"
)

func RunHTTPHandler() {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error load env: ", err.Error())
		return
	}

	e := echo.New()

	/* Create DB Connection */
	dbMaster := postgres.DBMasterConnection()

	/* CRUD module */
	crudRepository := crudRepository.NewRepository(dbMaster)
	crudUsecase := crudUsecase.NewUsecase(crudRepository)
	crudHandler := crudHandler.NewHandler(crudUsecase)
	crudHandler.Mount(e.Group("/api/v1/crud"))

	/* User Authentication */
	userAuthenticationRepository := userAuthenticationRepository.NewRepository(dbMaster)
	userAuthenticationUsecase := userAuthenticationUsecase.NewUsecase(userAuthenticationRepository)
	userAuthenticationHandler := userAuthenticationHandler.NewHandler(userAuthenticationUsecase)
	userAuthenticationHandler.Mount(e.Group("/api/v1/user-authentication"))

	/* User Management */
	userManagenentRepository := userManagementRepository.NewRepository(dbMaster)
	userManagementUsecase := userManagementUsecase.NewUsecase(userManagenentRepository)
	userManagementHandler := userManagementHandler.NewHandler(userManagementUsecase)
	userManagementHandler.Mount(e.Group("/api/v1/user-management"))

	/* File */
	fileRepository := fileRepository.NewRepository(dbMaster)
	fileUsecase := fileUsecase.NewUsecase(fileRepository)
	fileHandler := fileHandler.NewHandler(fileUsecase)
	fileHandler.Mount(e.Group("/api/v1/file"))

	/* Advance CRUD */
	advanceCrudRepository := advanceCrudRepository.NewRepository(dbMaster)
	advanceCrudUsecase := advanceCrudUsecase.NewUsecase(advanceCrudRepository)
	advanceCrudHandler := advanceCrudHandler.NewHandler(advanceCrudUsecase)
	advanceCrudHandler.Mount(e.Group("/api/v1/advance-crud"))

	/* HTTP Client */
	httpClientRepository := httpClientRepository.NewRepository()
	httpClientUsecase := httpClientUsecase.NewUsecase(httpClientRepository)
	httpClientHandler := httpClientHandler.NewHandler(httpClientUsecase)
	httpClientHandler.Mount(e.Group("/api/v1/http-client"))

	/* Email */
	emailRepository := emailRepository.NewRepository()
	emailUsecase := emailUsecase.NewUsecase(emailRepository)
	emailHandler := emailHandler.NewHandler(emailUsecase)
	emailHandler.Mount(e.Group("/api/v1/email"))

	/* SFTP */
	sftpRepository := sftpRepository.NewRepository()
	sftpUsecase := sftpUsecase.NewUsecase(sftpRepository)
	sftpHandler := sftpHandler.NewHandler(sftpUsecase)
	sftpHandler.Mount(e.Group("/api/v1/sftp"))

	e.Start(fmt.Sprintf("localhost:%s", os.Getenv(env.EnvPort)))
}
