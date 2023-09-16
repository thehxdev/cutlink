package main

import (
    "os"
    "log"
    "fmt"
    "flag"
    "net/http"
    "cutlink/models"

    "github.com/julienschmidt/httprouter"
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
    "github.com/joho/godotenv"
)


type app struct {
    Urls        *models.Urls
    ErrorLog    *log.Logger
    InfoLog     *log.Logger
    AdminToken  string
}


func main() {
    tls     := flag.Bool("tls", false, "Enable TLS - Must used with -crt and -key")
    crt     := flag.String("crt", "", "Path to .cert file for TLS")
    key     := flag.String("key", "", "Path to .key file for TLS")
    addr    := flag.String("addr", ":5000", "Listening Address")
    envFile := flag.String("env", "admin.env", "Path to .env file for ADMIN_TOKEN")
    dbFile  := flag.String("db", "./database.db", "Path to database file")
    flag.Parse()

    errLog  := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
    infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)


    err := godotenv.Load(*envFile)
    if err != nil {
        errLog.Fatal("Cannot open", *envFile)
        infoLog.Println("Please make a .env file and set ADMIN_TOKEN in that")
    }

    db, err := sqlx.Open("sqlite3", fmt.Sprintf("%s?parseTime=true", *dbFile))
    if err != nil {
        errLog.Fatal(err)
    }
    defer db.Close()

    urls := &models.Urls{
        DB: db,
    }

    app := &app{
        ErrorLog: errLog,
        InfoLog: infoLog,
        Urls: urls,
        AdminToken: os.Getenv("ADMIN_TOKEN"),
    }

    router := httprouter.New()

    router.GET("/", app.Root)
    router.GET("/r/:hash", app.Redirector)
    router.GET("/view", app.ViewUrl)
    router.GET("/all", app.ViewAll)
    router.GET("/search", app.SearchUrl)
    router.POST("/add", app.AddUrl)
    router.DELETE("/delete/:hash", app.Delete)

    app.InfoLog.Println("Listening on", *addr)
    if *tls {
        app.InfoLog.Println("TLS enabled")
        err = http.ListenAndServeTLS(*addr, *crt, *key, router)
    } else {
        err = http.ListenAndServe(*addr, router)
    }
    app.ErrorLog.Println(err)
}
