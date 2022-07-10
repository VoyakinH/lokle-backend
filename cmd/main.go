package main

import (

	// "database/sql"

	// "github.com/xuri/excelize/v2"
	"net/http"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/VoyakinH/lokle_backend/internal/user/delivery"
	"github.com/VoyakinH/lokle_backend/internal/user/repository"
	"github.com/VoyakinH/lokle_backend/internal/user/usecase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	config.SetConfig()

	// logger
	logger := logrus.New()

	// repository
	pr := repository.NewPostgresqlRepository(config.Postgres, *logger)
	rr := repository.NewRedisRepository(config.Redis, *logger)

	// usecase
	uu := usecase.NewUserUsecase(pr, rr, *logger)

	// delivery
	router := mux.NewRouter()
	delivery.SetUserRouting(router, uu, *logger)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Lokle.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}
	logger.Infof("starting server at %s\n", srv.Addr)

	logger.Fatal(srv.ListenAndServe())
}
