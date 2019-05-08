package api

import (
	"github.com/PhamDuyKhang/littledetective/internal/dataconnection"
	"github.com/PhamDuyKhang/littledetective/internal/film"
	"github.com/PhamDuyKhang/littledetective/internal/handler"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"github.com/gorilla/mux"
	"net/http"
)

func Init() *mux.Router {

	logger := flog.New()
	logger.SetLocal("api")
	// get connection, that is all database is supported
	// Elastic search
	client, err := dataconnection.GetElasticConnection()
	if err != nil {
		logger.Errorf("can't connect to elastic service with err:= %v", err)
	}
	elasticRepo := film.NewElasticFilmRepository(client)
	SearchService := film.NewFilmService(elasticRepo, logger)

	searchHandler := handler.NewSearchHandeler(SearchService, logger)

	//mongodb
	mgoSession, err := dataconnection.GetConnection()
	if err != nil {
		logger.Errorf("can't connect to mongodb service with err:= %v", err)
	}
	mongoRepo := film.NewFilmMongoRepository(mgoSession)
	mongoService := film.NewFilmMongoService(mongoRepo)
	filmHandler := handler.NewFilmHandler(mongoService, logger, SearchService)

	r := mux.NewRouter()
	r.HandleFunc("/search", searchHandler.Search)
	r.HandleFunc("/search/add", searchHandler.AddFilmToElastic).Methods(http.MethodPost)
	r.HandleFunc("/test", searchHandler.Test)
	//
	r.HandleFunc("/film", filmHandler.GetFilm).Methods(http.MethodGet)
	r.HandleFunc("/make", filmHandler.GenerateData).Methods(http.MethodGet)
	r.HandleFunc("/film/{id}", filmHandler.GetFilmWithID).Methods(http.MethodGet)
	r.HandleFunc("/film", filmHandler.AddFilm).Methods(http.MethodPost)
	r.HandleFunc("/sync", filmHandler.Sync).Methods(http.MethodGet)
	return r

}
