//	@title			Gopher Social API
//	@description	This is a server for a social media application called Gopher Social.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"time"

	"github.com/cprakhar/gopher-social/internal/auth"
	"github.com/cprakhar/gopher-social/internal/config"
	"github.com/cprakhar/gopher-social/internal/db"
	"github.com/cprakhar/gopher-social/internal/handler"
	"github.com/cprakhar/gopher-social/internal/mail"
	"github.com/cprakhar/gopher-social/internal/store"
	"go.uber.org/zap"
)

type application struct {
	config  config.Config
	handler handler.Handler
	logger  *zap.SugaredLogger
}

func main() {

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database connection pool
	db, err := db.New(ctx,
		cfg.DB.Addr,
		cfg.DB.MaxConns,
		cfg.DB.MinConns,
		cfg.DB.MaxIdleTime,
		cfg.DB.MaxConnLifetime,
	)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	store := store.NewStore(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.Auth.Token.Secret, cfg.Auth.Token.Aud, cfg.Auth.Token.Iss)

	mailer := mail.NewSendGrid(cfg.Mail.Sender, cfg.Mail.ApiKey)
	app := &application{
		config: cfg,
		handler: handler.Handler{
			Cfg:           cfg,
			Store:         store,
			Logger:        logger,
			Mailer:        mailer,
			Authenticator: jwtAuthenticator,
		},
		logger: logger,
	}

	mux := app.mount()
	logger.Infow("server is running", "addr", app.config.Addr, "env", app.config.Env)
	logger.Fatal(app.run(mux))
}
