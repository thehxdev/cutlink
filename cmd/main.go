package main

import (
    "os"
    "log"
    "fmt"
    "flag"
    "time"
    "database/sql"
    "cutlink/models"

    _ "github.com/mattn/go-sqlite3"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html/v2"
    "github.com/gofiber/fiber/v2/middleware/session"
    fiberSQLstore "github.com/gofiber/storage/sqlite3/v2"
)


type cutlink struct {
    App             *fiber.App
    Conn            *models.Conn
    ErrorLog        *log.Logger
    InfoLog         *log.Logger
    Store           *session.Store
    DisableSignup   bool
}


func main() {
    addr     := flag.String("addr", ":5000", "Listening Address")
    dbFile   := flag.String("db", "./database.db", "Path to database file")
    noSignUp := flag.Bool("disable-signup", false, "Disable user signup")
    flag.Parse()


    db, err := sql.Open("sqlite3", fmt.Sprintf("%s?parseTime=true", *dbFile))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()


    cl := &cutlink{
        App: fiber.New(fiber.Config{
            Views: html.New("./ui/html", ".html"),
        }),

        ErrorLog: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),

        InfoLog: log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),

        Conn: &models.Conn{ DB: db, },

        Store: session.New(session.Config{
            Expiration: 12 * time.Hour,
            CookieHTTPOnly: true,
            CookieSecure: true,
            Storage: fiberSQLstore.New(fiberSQLstore.Config{
                Database: "./sessions.db",
                GCInterval: 30 * time.Second,
            }),
        }),

        DisableSignup: *noSignUp,
    }

    // setup static file server first
    cl.App.Static("/static", "./ui/static", fiber.Static{
        Browse: false,
        CacheDuration: 10 * time.Second,
    })

    cl.setupMiddlewares()
    cl.setupRoutes()

    cl.InfoLog.Println("Server listening on", *addr)
    err = cl.App.Listen(*addr)
    cl.ErrorLog.Println(err)
}
