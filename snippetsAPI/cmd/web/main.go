package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/earnstein/GO/snippetsAPI/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Application struct {
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	snippets       *models.SnippetModel
	sessionManager *scs.SessionManager
	userManager    *models.UserModel
}

func newApplicaton(infoLogger, errorLogger *log.Logger, db *sql.DB, sessionManager *scs.SessionManager) *Application {
	return &Application{
		InfoLog:        infoLogger,
		ErrorLog:       errorLogger,
		snippets:       &models.SnippetModel{DB: db},
		userManager:    &models.UserModel{DB: db},
		sessionManager: sessionManager,
	}
}

func newServer(addr string, logger *log.Logger, app *Application, tlsconfig *tls.Config) *http.Server {
	return &http.Server{
		Addr:         addr,
		ErrorLog:     logger,
		Handler:      app.routes(),
		TLSConfig:    tlsconfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func newTls() *tls.Config {
	return &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		errorLog.Fatal(err)
	}

	// command flags
	addr := flag.String("addr", os.Getenv("ADDR"), "HTTP network port address")
	dsn := flag.String("dsn", os.Getenv("DB_CONN_STRING"), "MYSQL data source name")
	flag.Parse()

	// database configuration
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// session configuration
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// tls configuration
	tlsconfig := newTls()

	// server
	app := newApplicaton(infoLog, errorLog, db, sessionManager)
	server := newServer(*addr, errorLog, app, tlsconfig)
	infoLog.Printf("server is listening on port %s", *addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
