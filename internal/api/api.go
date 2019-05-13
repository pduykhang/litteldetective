package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/PhamDuyKhang/littledetective/internal/externalconnection"
	"github.com/PhamDuyKhang/littledetective/internal/film"
	"github.com/PhamDuyKhang/littledetective/internal/handler"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

type (
	ExternalConn struct {
		Database     externalconnection.DBConnection
		SearchEngine externalconnection.SearchConnection
	}
	middlewareFunc = func(http.HandlerFunc) http.HandlerFunc
	route          struct {
		url        string
		method     string
		handler    http.HandlerFunc
		middleware []middlewareFunc
	}
)

func Init(exc *ExternalConn) *mux.Router {
	logger := flog.New()
	logger.SetLocal("api")
	// get connection, that is all database is supported
	// Elastic search

	elasticRepo := film.NewElasticFilmRepository(exc.SearchEngine.Client)
	searchService := film.NewFilmService(elasticRepo, logger)

	searchHandler := handler.NewSearchHandler(searchService, logger)

	//mongodb
	mongoRepo := film.NewFilmMongoRepository(exc.Database.Session)
	mongoService := film.NewFilmMongoService(mongoRepo)
	filmHandler := handler.NewFilmHandler(mongoService, logger, searchService)

	r := mux.NewRouter()

	routes := []route{
		{
			url:     "/api/v1/test",
			handler: searchHandler.Test,
			method:  get,
		},
		{
			url:     "/api/v1/films/init",
			handler: filmHandler.GenerateData,
			method:  get,
		},
		{
			url:     "/api/v1/films/sync",
			handler: filmHandler.Sync,
			method:  get,
		},
		{
			url:     "/api/v1/films/search",
			handler: searchHandler.Search,
			method:  get,
		},
		{
			url:     "/api/v1/films",
			handler: filmHandler.GetFilm,
			method:  get,
		},
		{
			url:     "/api/v1/films",
			handler: filmHandler.AddFilm,
			method:  post,
		},
		{
			url:     "/api/v1/film/{id}",
			handler: filmHandler.GetFilmWithID,
			method:  get,
		},
		{
			url:     "/api/v1/film/search",
			handler: searchHandler.AddFilmToElastic,
			method:  post,
		},
	}
	logger.Infof("init router ")

	for _, rt := range routes {
		h := rt.handler
		for i := len(rt.middleware) - 1; i >= 0; i-- {
			h = rt.middleware[i](h)
		}
		r.Path(rt.url).Methods(rt.method).HandlerFunc(h)
	}
	return r

}
