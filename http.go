package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *application) mount() *gin.Engine {
	r := gin.Default()

	api := r.Group("/v1")
	api.GET("/health", app.handler.HealthCheckHandler)

	users := api.Group("/users")
	users.POST("/", app.handler.RegisterUserHandler)
	usersID := users.Group("/:id")
	usersID.Use(app.handler.UsersContextMiddleware)
	usersID.GET("/", app.handler.GetUserHandler)
	usersID.PUT("/follow", app.handler.FollowUserHandler)
	usersID.PUT("/unfollow", app.handler.UnfollowUserHandler)

	posts := api.Group("/posts")
	posts.POST("/", app.handler.CreatePostHandler)
	postsID := posts.Group("/:id")
	postsID.Use(app.handler.PostsContextMiddleware)
	postsID.GET("/", app.handler.GetPostHandler)
	postsID.PATCH("/", app.handler.UpdatePostHandler)
	postsID.DELETE("/", app.handler.DeletePostHandler)

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
