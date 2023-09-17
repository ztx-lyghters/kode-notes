package main

import (
	"net/http"

	"github.com/ztx-lyghters/kode-notes/config"
	"github.com/ztx-lyghters/kode-notes/handlers"
	"github.com/ztx-lyghters/kode-notes/repository"
	"github.com/ztx-lyghters/kode-notes/server"
	"github.com/ztx-lyghters/kode-notes/services"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	srv_cfg := &server.Config{
		Host: "localhost",
		Port: "8080",
	}
	db_cfg := &repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		DBName:   "kode_notes",
		SSLMode:  "disable",
	}
	app_cfg := &config.Config{
		Spellcheck: true,
	}

	err := config.ParseArgs(srv_cfg, db_cfg, app_cfg)
	if err != nil {
		logrus.Fatalf("Couldn't parse arguments: %s", err)
	}

	db, err := repository.NewDBPostgres(db_cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		logrus.Fatalf("database is down: %s", err.Error())
	}

	repo := repository.New(db)
	if err != nil {
		logrus.Fatal(err)
	}

	mux := http.NewServeMux()
	services := services.New(repo)
	handler := handlers.New(services, mux, app_cfg)

	server := server.New(srv_cfg)
	err = server.Run(handler)
	if err != nil {
		logrus.Fatal(err)
	}
}
