package main

import (
	"flag"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mrlightwood/golang-products-api/api"
	"github.com/mrlightwood/golang-products-api/config"
	"github.com/mrlightwood/golang-products-api/db"
	"github.com/mrlightwood/golang-products-api/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	log.SetFormatter(&log.JSONFormatter{})
	// Setting app launch flags
	configFile := flag.String("conf", "./config/config.yaml", "Path to config file")
	flag.Parse()
	// Config load
	var conf *config.Config
	if conf, err = config.NewConfig(*configFile); err != nil {
		log.Fatal(err)
	}
	// Logger
	log.SetLevel(log.Level(conf.LogLevel))
	log.Info("Starting service with configuration: ", conf.ConfigFile)

	// Storage creation
	store, err := db.NewStore(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	log.Info("Store created successfully")

	// Initialization of services
	cs := service.NewCategoryService(store)
	ps := service.NewProductService(store)
	log.Info("Services created successfully")

	// Initialization of an API
	api := api.NewApi(conf, cs, ps)
	log.WithField("address", api.GetApiInfo().Address).
		WithField("mw", api.GetApiInfo().MW).
		WithField("routes", api.GetApiInfo().Routes).
		Info("Starting api")
	log.Fatal(api.Start())
}
