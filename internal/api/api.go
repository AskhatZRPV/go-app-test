package api

import (
	"database/sql"
	"golang-app/internal/config"
	"golang-app/internal/store"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Server struct {
	handler http.Handler
	router  *httprouter.Router
	logger  *logrus.Logger
	store   store.Store
	db      *sql.DB
}

func (s *Server) WithRouter() *Server {
	router := httprouter.New()
	router.POST("/visit", Log(NewVisitDoctor(s.db, s.store).Handle))
	s.handler = router
	return s
}

func Start(config *config.Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		logrus.Info(err)
	}
	defer db.Close()
	store := store.New(db)
	store.Repository = store.GetNewRepository()
	srv := newServer(store)
	go srv.notificationsWithDelay(5)
	srv.logger.Info("Server started on http://127.0.0.1", config.ServerPort)
	return http.ListenAndServe(config.ServerPort, srv.WithRouter().handler)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func newServer(store *store.Store) *Server {
	s := &Server{
		router: httprouter.New(),
		logger: logrus.New(),
		store:  *store,
	}
	return s
}

func (s *Server) notificationsWithDelay(n time.Duration) {
	for {
		for range time.Tick(n * time.Minute) {
			err := s.store.Repository.Notification()
			if err != nil {
				s.logger.Error(err)
			}
		}
	}
}
