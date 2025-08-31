package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *application) mount() *gin.Engine {
	r := gin.Default()

	apiGrp := r.Group("/v1")
	apiGrp.GET("/health", app.handler.HealthCheckHandler)

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
