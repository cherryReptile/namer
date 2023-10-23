package app

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	handler "namer/internal/delivery/http"
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
		if err := a.server.ListenAndServe(); err != nil {
			log.Error(errors.Wrap(err, "Start #2"))
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

	if a.db, err = ConnectToDB(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", errors.Wrap(err, "newApp #2"))
	}

	a.server = NewServer(initRouter(a.db))

	return &a
}

func initRouter(db *sql.DB) *gin.Engine {
	g := gin.New()

	h := handler.NewHandler(db)

	api := g.Group("/api/person", gin.Logger())

	api.Handle(http.MethodPost, "", h.NewPerson)
	api.Handle(http.MethodGet, "/:id", h.GetPerson)
	api.Handle(http.MethodPost, "/filter", h.GetPersons)
	api.Handle(http.MethodPut, "/:id", h.UpdatePerson)
	api.Handle(http.MethodDelete, "/:id", h.DeletePerson)

	return g
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
