package main

import (
	"net/http"
	"time"

	"github.com/cprakhar/gopher-social/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) mount() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{app.config.WebURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Version = app.config.Version
	docs.SwaggerInfo.Host = app.config.ApiURL

	api := r.Group("/v1")
	{
		api.GET("/health", app.handler.BasicAuthMiddleware, app.handler.HealthCheckHandler)

		api.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		users := api.Group("/users")
		{
			users.PUT("/activate/:token", app.handler.ActivateUserHandler)
			userfeed := users.Group("/feed")
			{
				userfeed.Use(app.handler.AuthTokenMiddleware)
				userfeed.GET("/", app.handler.GetUserFeedHandler)
			}
			usersID := users.Group("/:id")
			{
				usersID.Use(app.handler.AuthTokenMiddleware)
				usersID.GET("/", app.handler.GetUserHandler)
				usersID.PUT("/follow", app.handler.FollowUserHandler)
				usersID.PUT("/unfollow", app.handler.UnfollowUserHandler)
			}
		}
		authenticate := api.Group("/authenticate")
		{
			authenticate.POST("/user", app.handler.RegisterUserHandler)
		}
		posts := api.Group("/posts")
		{
			posts.Use(app.handler.AuthTokenMiddleware)
			posts.POST("/", app.handler.CreatePostHandler)
			postsID := posts.Group("/:id")
			{
				postsID.Use(app.handler.PostsContextMiddleware)
				postsID.GET("/", app.handler.GetPostHandler)
				postsID.PATCH("/", app.handler.CheckPostOwnership("moderator", app.handler.UpdatePostHandler))
				postsID.DELETE("/", app.handler.CheckPostOwnership("admin", app.handler.DeletePostHandler))

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
