package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/justteddy/wallet/handlers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

var (
	env             = flag.String("env", "dev", "environment (dev|prod)")
	port            = flag.String("port", ":8080", "port for http connections")
	shutdownTimeout = flag.Duration("shutdown-timeout", time.Second*5, "shutdown timeout")

	dbDSN      = flag.String("db-dsn", "postgres://postgres:postgres@localhost:5432?sslmode=disable", "database dsn")
	dbConnPool = flag.Int("db-conn-pool", 10, "database connection pool")
)

func main() {
	flag.Parse()
	setupLogger(*env)

	dbConn, err := setupDatabase(*dbDSN, *dbConnPool)
	mustNoError(err)

	httpServer := setupHTTPServer(*port, setupRouter(nil))

	httpErrCh := startHTTPServer(httpServer)
	log.Infof("service is ready to accept connections on port %s", *port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		log.Info("received signal to stop service")
		shutdown(httpServer, dbConn, *shutdownTimeout)
	case <-httpErrCh:
		log.WithError(err).Error("http server error")
		shutdown(httpServer, dbConn, *shutdownTimeout)
	}

	log.Info("bye 👋")
}

func setupHTTPServer(port string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: router,
	}
}

func startHTTPServer(server *http.Server) <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		if err := server.ListenAndServe(); err != nil {
			ch <- err
			return
		}
	}()

	return ch
}

func setupDatabase(dsn string, connPool int) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", *dbDSN)
	if err != nil {
		return nil, errors.Wrap(err, "database connection")
	}

	conn.SetMaxOpenConns(connPool)
	conn.SetMaxIdleConns(connPool)

	return conn, errors.Wrap(conn.Ping(), "ping database")
}

func setupLogger(env string) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if env == "prod" {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.JSONFormatter{})
	}

}

func setupRouter(handler *handlers.Handler) http.Handler {
	router := httprouter.New()
	router.POST("/wallet", handler.HandleCreateWallet)
	router.POST("/deposit/:wallet", handler.HandleDeposit)
	router.POST("/transfer", handler.HandleTransfer)

	return router
}

func shutdown(httpServer *http.Server, dbConn *sqlx.DB, shutdownTimeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.WithError(err).Error("http server shutdown")
	}
	log.Info("http server stopped")

	if err := dbConn.Close(); err != nil {
		log.WithError(err).Error("db conn close")
	}
	log.Info("db connection closed")
}

func mustNoError(err error) {
	if err != nil {
		log.WithError(err).Fatal("error occured")
	}
}
