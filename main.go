package main

import (
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/namrahov/hesen-go/config"
	"github.com/namrahov/hesen-go/handler"
	"github.com/namrahov/hesen-go/repo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var opts struct {
	Profile string `short:"p" long:"profile" default:"dev" description:"Application run profile"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	initLogger()
	initEnvVars()
	config.LoadConfig()
	applyLogLevel()

	log.Info("Application is starting with profile: ", opts.Profile)

	err = repo.MigrateDb()
	if err != nil {
		log.Fatal(err)
	}
	repo.InitDb()

	router := mux.NewRouter()
	handler.ApplicationHandler(router)

	log.Info("Starting server at port: ", config.Props.Port)
	log.Fatal(http.ListenAndServe(":"+config.Props.Port, router))

}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	if opts.Profile == "default" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func initEnvVars() {
	if godotenv.Load("profiles/default.env") != nil {
		log.Fatal("Error in loading environment variables from: profiles/default.env")
	} else {
		log.Info("Environment variables loaded from: profiles/default.env")
	}

	if opts.Profile != "default" {
		profileFileName := "profiles/" + opts.Profile + ".env"
		if godotenv.Overload(profileFileName) != nil {
			log.Fatal("Error in loading environment variables from: ", profileFileName)
		} else {
			log.Info("Environment variables overloaded from: ", profileFileName)
		}
	}
}

func applyLogLevel() {
	log.SetLevel(config.Props.LogLevel)
}
