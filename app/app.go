package app

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ludete/wechat_robot/wallets"

	"github.com/ludete/wechat_robot/exchanges"

	"github.com/gorilla/mux"
	"github.com/ludete/wechat_robot/storage"
	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

type RobotApp struct {
	dbMutex sync.Mutex
	db      storage.DB

	server   *http.Server
	exchange exchanges.Exchange
	wallet   wallets.WalletInterface
	advert   string
	resURL   string
}

func NewRobotApp(cfg *toml.Tree) *RobotApp {
	dbPath := cfg.GetDefault("db", "data").(string)
	app := &RobotApp{
		dbMutex: sync.Mutex{},
		db:      storage.NewDB(dbPath),
		exchange: exchanges.NewExchanges(
			cfg.GetDefault("coinex", "").(string),
			cfg.GetDefault("binance", "").(string),
		),
		wallet: wallets.NewWallet(
			cfg.GetDefault("wallet", "").(string),
			cfg.GetDefault("apikey", "").(string),
			cfg.GetDefault("secretkey", "").(string),
		),
		resURL: cfg.GetDefault("proxy", "").(string),
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
			log.Errorf("listen server error : %s\n", err.Error())
		}
	}()
}
