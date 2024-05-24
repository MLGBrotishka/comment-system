// Package app configures and runs application.
package app

import (
	"comment-system/config"
	"comment-system/internal/adapters/graphql"
	"comment-system/internal/adapters/repo/memory"
	"comment-system/internal/adapters/repo/postgre"
	"comment-system/internal/adapters/server"
	"comment-system/internal/usecase"
	"comment-system/pkg/httpserver"
	"comment-system/pkg/logger"
	"comment-system/pkg/postgres"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	var err error
	l := logger.New(cfg.Log.Level)

	var commsRepo usecase.CommentsRepo
	var postsRepo usecase.PostsRepo
	if cfg.App.InMemory {
		commsRepo = memory.NewMemoryCommentsRepo(cfg.CommsRepo.Limit, cfg.CommsRepo.Offset)
		postsRepo = memory.NewMemoryPostsRepo(cfg.PostsRepo.Limit, cfg.PostsRepo.Offset)
	} else {
		Migrate(cfg.PG.URL)

		pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
		}
		defer pg.Close()

		commsRepo = postgre.NewComments(pg, cfg.CommsRepo.Limit, cfg.CommsRepo.Offset)
		postsRepo = postgre.NewPostsRepo(pg, cfg.CommsRepo.Limit, cfg.CommsRepo.Offset)
	}
	postsUC := usecase.NewPosts(postsRepo)
	resolver := graphql.NewResolver(
		postsUC,
		usecase.NewComment(commsRepo, postsUC),
	)

	router := http.NewServeMux()

	qhandler := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	server.NewRouter(router, qhandler)
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
