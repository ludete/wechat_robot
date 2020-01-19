package app

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/ludete/wechat_robot/storage"
	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

type RobotApp struct {
	dbMutex sync.Mutex
	db      storage.DB
	server  *http.Server

	userID   string
	password string
}

func NewRobotApp(cfg *toml.Tree) *RobotApp {
	dbPath := cfg.GetDefault("db", "data").(string)
	app := &RobotApp{
		dbMutex:  sync.Mutex{},
		db:       storage.NewDB(dbPath),
		userID:   cfg.Get("user-name").(string),
		password: cfg.Get("password").(string),
	}

	router := registerHandler(app)
	httpSvr := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetDefault("port", 9789).(int64)),
		Handler:      router,
		ReadTimeout:  READTIMEOUT * time.Second,
		WriteTimeout: WRITETIMEOUT * time.Second,
	}
	app.server = httpSvr
	return app
}

func registerHandler(app *RobotApp) *mux.Router {
	route := mux.NewRouter()
	route.HandleFunc("/", handler(app)).Methods("POST")
	return route
}

func (app *RobotApp) Start() {
	log.Info("robot begin start")
	go func() {
		if err := app.server.ListenAndServe(); err != nil {

		}
	}()
}
