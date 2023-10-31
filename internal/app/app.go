package app

import (
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	handler "namer/internal/delivery/http"
	"namer/internal/delivery/http/middlewares"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type app struct {
	db     *sql.DB
	server *http.Server
}

func Start() {
	ctx := context.Background()

	a := newApp(ctx)

	go func() {
		log.Info("starting server at ", a.server.Addr)

		if err := a.server.ListenAndServe(); err != nil {
			log.Error(errors.Wrap(err, "Start #1"))
		}
	}()

	a.quit(ctx)
}

func newApp(ctx context.Context) *app {
	var (
		a   app
		err error
	)

	if err = godotenv.Load(".env"); err != nil {
		log.Fatal("failed to load env: ", errors.Wrap(err, "newApp #1"))
	}

	if a.db, err = connectToDB(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", errors.Wrap(err, "newApp #2"))
	}

	a.server = newServer(initRouter(a.db))

	return &a
}

func initRouter(db *sql.DB) *echo.Echo {
	e := echo.New()

	h := handler.NewHandler(db)

	api := e.Group("/api/person", middleware.Logger(), middlewares.ErrLogger())

	api.POST("", h.NewPerson)
	api.GET("/:id", h.GetPerson)
	api.POST("/filter", h.GetPersons)
	api.PUT("/:id", h.UpdatePerson)
	api.DELETE("/:id", h.DeletePerson)

	return e
}

func (a *app) quit(ctx context.Context) {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	log.Warnf("started shutdown")

	a.shutdown(ctx)
}

func (a *app) shutdown(ctx context.Context) {
	if err := a.server.Shutdown(ctx); err != nil {
		log.Error(errors.Wrap(err, "shutdown #1"))
	}

	if err := a.db.Close(); err != nil {
		log.Error(errors.Wrap(err, "shutdown #2"))
	}

	log.Info("services closed")
}
