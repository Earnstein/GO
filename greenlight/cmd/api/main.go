package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func NewApplication(config config) *application {
	return &application{
		config: config,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func NewServer(config config, app *application) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	app := NewApplication(cfg)
	srv := NewServer(cfg, app)

	app.logger.Printf("Starting %s server on port %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	app.logger.Fatal(err)
}
