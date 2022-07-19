package main

import (

	// "database/sql"

	// "github.com/xuri/excelize/v2"
	"net/http"

	"github.com/VoyakinH/lokle_backend/config"
	file_manager "github.com/VoyakinH/lokle_backend/internal/file"
	"github.com/VoyakinH/lokle_backend/internal/pkg/middleware"
	user_delivery "github.com/VoyakinH/lokle_backend/internal/user/delivery"
	user_repository "github.com/VoyakinH/lokle_backend/internal/user/repository"
	user_usecase "github.com/VoyakinH/lokle_backend/internal/user/usecase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	config.SetConfig()

	// logger
	logger := logrus.New()

	// repository
	pr := user_repository.NewPostgresqlRepository(config.Postgres, *logger)
	rsr := user_repository.NewRedisSessionRepository(config.RedisSession, *logger)
	rur := user_repository.NewRedisUserRepository(config.RedisUser, *logger)

	// usecase
	uu := user_usecase.NewUserUsecase(pr, rsr, rur, *logger)

	// middlewars
	auth := middleware.NewAuthMiddleware(uu, *logger)

	// delivery
	router := mux.NewRouter()
	user_delivery.SetUserRouting(router, uu, auth, *logger)
	file_manager.SetFileRouting(router, uu, auth, *logger)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Lokle.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}
	logger.Infof("starting server at %s\n", srv.Addr)

	logger.Fatal(srv.ListenAndServe())
}
