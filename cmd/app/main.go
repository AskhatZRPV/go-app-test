package mainapp

import (
	"flag"
	"golang-app/internal/api"
	"golang-app/internal/config"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var flagConfig = flag.String("config", "./config/config.yml", "path to the config file")

func Start() {

	logger := logrus.New()
	config, err := config.Load(*flagConfig)
	if err != nil {
		logger.Error(err)
	}
	err = api.Start(config)
	if err != nil {
		panic(err)
	}
}
