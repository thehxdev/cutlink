package main

import (
    "os"
    "log"
    "fmt"
    "flag"
    "time"
    "cutlink/models"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html/v2"
    "github.com/gofiber/fiber/v2/middleware/session"
)


type cutlink struct {
    App         *fiber.App
    Conn        *models.Conn
    ErrorLog    *log.Logger
    InfoLog     *log.Logger
    Store       *session.Store
}


func main() {
    addr    := flag.String("addr", ":5000", "Listening Address")
    dbFile  := flag.String("db", "./database.db", "Path to database file")
    flag.Parse()

    infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
    errLog  := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)


    db, err := sqlx.Open("sqlite3", fmt.Sprintf("%s?parseTime=true", *dbFile))
    if err != nil {
        errLog.Fatal(err)
    }
    defer db.Close()


    engine := html.New("./ui/html", ".html")
    store := session.New(session.Config{
        Expiration: 12 * time.Hour,
        CookieHTTPOnly: true,
        CookieSecure: true,
    })
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    cl := &cutlink{
        App: app,
        ErrorLog: errLog,
        InfoLog: infoLog,
        Conn: &models.Conn{ DB: db, },
        Store: store,
    }

    cl.setupMiddlewares()
    cl.setupRoutes()

    cl.InfoLog.Println("Server listening on", *addr)
    err = cl.App.Listen(*addr)
    cl.ErrorLog.Println(err)
}
