package main

import (
	"context"
	"fmt"
	"github.com/PhamDuyKhang/littledetective/internal/api"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	r := api.Init()

	logger := flog.New()
	logger.SetLocal("main")
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
