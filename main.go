package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PhamDuyKhang/littledetective/internal/api"
	"github.com/PhamDuyKhang/littledetective/internal/externalconnection/elasticsearch"
	"github.com/PhamDuyKhang/littledetective/internal/externalconnection/mongodb"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

func main() {
	logger := flog.New()
	logger.SetLocal("main")
	exc, err := PrepareConn(logger)
	defer exc.Close()
	if err != nil {
		log.Fatal("can't connection to necessary servers")
		return
	}

	r := api.Init(exc)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", "", "8080"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panicf("http.ListenAndServe() error: %v", err)
		}
	}()
	logger.Infof("HTTP Server is listening at: %v", srv.Addr)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	logger.Infof("shutting down http server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("http server shutdown with error: %v", err)
	}

}
func PrepareConn(l flog.Logger) (*api.ExternalConn, error) {
	log.Println("connecting to database")
	ec := &api.ExternalConn{}
	ss, err := mongodb.GetConnection()
	if err != nil {
		l.Errorf("error when connect to database with err = %v", err)
		return nil, err
	}
	ec.Database.Session = ss
	client, err := elasticsearch.GetElasticConnection()
	if err != nil {
		l.Errorf("error when connect to database with err = %v", err)
		return nil, err
	}
	ec.SearchEngine.Client = client

	return ec, nil
}
