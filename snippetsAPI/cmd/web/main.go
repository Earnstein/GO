package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewApplicaton(infoLogger, errorLogger *log.Logger) *Application {
	return &Application{
		InfoLog:  infoLogger,
		ErrorLog: errorLogger,
	}
}

func NewServer(addr string, logger *log.Logger, app *Application) *http.Server {
	return &http.Server{
		Addr:     addr,
		ErrorLog: logger,
		Handler:  app.routes(),
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
	db , err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// server
	app := NewApplicaton(infoLog, errorLog)
	server := NewServer(*addr, errorLog, app)
	infoLog.Printf("server is listening on port %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
