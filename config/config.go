package config

import (
	"github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

const RootPath = "/v1/ecourt"
const InternalRootPath = "/internal/v1/ecourt"

type Args struct {
	DbHost      string    `arg:"env:DB_CARD_DELIVERY_HOST,required"`
	DbPort      string    `arg:"env:DB_CARD_DELIVERY_PORT,required"`
	DbName      string    `arg:"env:DB_CARD_DELIVERY_NAME,required"`
	DbUser      string    `arg:"env:DB_CARD_DELIVERY_USER,required"`
	DbPass      string    `arg:"env:DB_CARD_DELIVERY_PASS,required"`
	AdminClient string    `arg:"env:MS_ADMIN,required"`
	LogLevel    log.Level `arg:"env:LOG_LEVEL"`
	Port        string    `arg:"env:PORT,required"`
	Hostname    string    `arg:"env:HOSTNAME,required"`
}

var Props Args

func LoadConfig() {
	arg.Parse(&Props)
}
