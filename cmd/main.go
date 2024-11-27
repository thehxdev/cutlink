package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/thehxdev/cutlink/models"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberSQLstore "github.com/gofiber/storage/sqlite3/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/viper"
)

type cutlink struct {
	App      *fiber.App
	Conn     *models.Conn
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Store    *session.Store
	Cfg      *Config
}

func main() {
	configPath := flag.String("cfg", "", "Path to config file")
	uiPath := flag.String("ui", "./ui", "Path to ui directory")
	flag.Parse()

	cfg := &Config{}
	setupViper(cfg, *configPath)

	var dbIsNew bool = false
	dbPath := viper.GetString("database.mainDB")
	sessDB := viper.GetString("database.sessionsDB")

	if _, err := os.Stat(dbPath); err != nil {
		dbIsNew = true
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?parseTime=true", dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	htmlFiles := fmt.Sprintf("%s/html", *uiPath)
	cl := &cutlink{
		App: fiber.New(fiber.Config{
			Views: html.New(htmlFiles, ".html"),
		}),

		ErrorLog: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),

		InfoLog: log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),

		Conn: &models.Conn{DB: db},

		Store: session.New(session.Config{
			Expiration:     12 * time.Hour,
			CookieHTTPOnly: true,
			CookieSecure:   true,
			Storage: fiberSQLstore.New(fiberSQLstore.Config{
				Database:   sessDB,
				GCInterval: 30 * time.Second,
			}),
		}),

		Cfg: cfg,
	}

	if dbIsNew {
		err = cl.Conn.MigrateDB()
		if err != nil {
			cl.ErrorLog.Println("Cannot create database tables")
			return
		}
	}

	// setup static file server first
	staticFiles := fmt.Sprintf("%s/static", *uiPath)
	cl.App.Static("/static", staticFiles, fiber.Static{
		Browse:        false,
		CacheDuration: 10 * time.Second,
	})

	cl.setupMiddlewares()
	cl.setupRoutes()

	fullAddr := fmt.Sprintf("%s:%d", cl.Cfg.Server.Addr, cl.Cfg.Server.Port)
	cl.InfoLog.Println("Server listening on", fullAddr)

	if cl.Cfg.Tls.Cert != "" && cl.Cfg.Tls.Key != "" {
		cl.InfoLog.Println("TLS encryption enabled by Cutlink")
		err = cl.App.ListenTLS(fullAddr, cl.Cfg.Tls.Cert, cl.Cfg.Tls.Key)
	} else {
		err = cl.App.Listen(fullAddr)
	}

	cl.ErrorLog.Println(err)
}
