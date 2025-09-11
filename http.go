package main

import (
	"net/http"
	"time"

	"github.com/cprakhar/gopher-social/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) mount() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Version = app.config.Version
	docs.SwaggerInfo.Host = app.config.ApiURL

	api := r.Group("/v1")
	{
		api.GET("/health", app.handler.HealthCheckHandler)

		api.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		users := api.Group("/users")
		{
			users.POST("/", app.handler.RegisterUserHandler)
			userfeed := users.Group("/feed")
			{
				userfeed.GET("/", app.handler.GetUserFeedHandler)
			}
			usersID := users.Group("/:id")
			{
				usersID.Use(app.handler.UsersContextMiddleware)
				usersID.GET("/", app.handler.GetUserHandler)
				usersID.PUT("/follow", app.handler.FollowUserHandler)
				usersID.PUT("/unfollow", app.handler.UnfollowUserHandler)
			}
		}
		posts := api.Group("/posts")
		{
			posts.POST("/", app.handler.CreatePostHandler)
			postsID := posts.Group("/:id")
			{
				postsID.Use(app.handler.PostsContextMiddleware)
				postsID.GET("/", app.handler.GetPostHandler)
				postsID.PATCH("/", app.handler.UpdatePostHandler)
				postsID.DELETE("/", app.handler.DeletePostHandler)

			}
		}
	}

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
