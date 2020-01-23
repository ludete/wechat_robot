package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	advert   []string
	coins    map[string]struct{}
}

func NewRobotApp(cfg *toml.Tree) (*RobotApp, error) {
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
		coins: make(map[string]struct{}, 300),
	}

	router := registerHandler(app)
	httpSvr := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetDefault("port", 9789).(int64)),
		Handler:      router,
		ReadTimeout:  READTIMEOUT * 20 * time.Second,
		WriteTimeout: WRITETIMEOUT * 20 * time.Second,
	}
	app.server = httpSvr
	if err := app.readCoinSymbols(cfg); err != nil {
		return nil, err
	}
	if err := app.initAdvertMsg(cfg); err != nil {
		return nil, err
	}
	return app, nil
}

func (app *RobotApp) initAdvertMsg(cfg *toml.Tree) error {
	advertFile := cfg.GetDefault("advert_file", "").(string)
	file, err := os.Open(advertFile)
	if err != nil {
		return err
	}
	defer file.Close()
	bz, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	app.advert = strings.Split(string(bz), "\n")
	return nil
}

func (app *RobotApp) readCoinSymbols(cfg *toml.Tree) error {
	filePath := cfg.GetDefault("symbols", "").(string)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	bz, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	coins := strings.Split(string(bz), "\n")
	for _, c := range coins {
		app.coins["="+strings.ToLower(c)] = struct{}{}
	}
	return nil
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
