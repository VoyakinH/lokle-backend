package main

import (

	// "database/sql"

	// "github.com/xuri/excelize/v2"
	"net/http"

	"github.com/VoyakinH/lokle_backend/config"
	file_manager "github.com/VoyakinH/lokle_backend/internal/file"
	"github.com/VoyakinH/lokle_backend/internal/pkg/middleware"
	reg_req_delivery "github.com/VoyakinH/lokle_backend/internal/reg_req/delivery"
	reg_req_repository "github.com/VoyakinH/lokle_backend/internal/reg_req/repository"
	reg_req_usecase "github.com/VoyakinH/lokle_backend/internal/reg_req/usecase"
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
	ur := user_repository.NewPostgresqlRepository(config.Postgres, *logger)
	rsr := user_repository.NewRedisSessionRepository(config.RedisSession, *logger)
	rur := user_repository.NewRedisUserRepository(config.RedisUser, *logger)
	rrr := reg_req_repository.NewPostgresqlRepository(config.Postgres, *logger)

	// router
	router := mux.NewRouter()

	// usecase
	uu := user_usecase.NewUserUsecase(ur, rsr, rur, *logger)

	// middlewars
	auth := middleware.NewAuthMiddleware(uu, *logger)
	roleMw := middleware.NewRoleMiddleware(uu, *logger)

	// files
	fm := file_manager.SetFileRouting(router, uu, auth, *logger)

	// usecase
	rru := reg_req_usecase.NewRegReqUsecase(rrr, ur, fm, *logger)

	// delivery
	user_delivery.SetUserRouting(router, uu, auth, roleMw, *logger)
	reg_req_delivery.SetRegReqRouting(router, rru, auth, roleMw, *logger)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Lokle.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}
	logger.Infof("starting server at %s\n", srv.Addr)

	logger.Fatal(srv.ListenAndServeTLS("kit-lokle.crt", "kit-lokle.key"))
}
