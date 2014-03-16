package main

import (
    "database/sql"
    "fmt"
    "github.com/codegangsta/martini"
    //"github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "log"
    "os"
)

type Page struct {
    url     string
    content string
}

func setup() {
    //TODO:Swap the loggin out for go-logging or log4go, log4go might be ther best choice
    log.SetOutput(os.Stdout)
    log.Println("Start logging")

}

//TODO: After experimentation follow a BDD approach
func main() {
    db, err := sql.Open("postgres", "user=psql dbname=footprints password=Anilihst1 host=localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    //TODO:Implement the sqlx struct mapping
    rows, err := db.Query("SELECT * FROM pages")
    if err != nil {
        log.Fatal(err)
    }
    var pageId int
    var url string
    rows.Next()
    rows.Scan(&pageId, &url)
    fmt.Println("Rows:", rows)
    fmt.Println("Fields:", pageId)
    fmt.Println("Url: ", url)

    //config := jsonConfig.LoadConfig()
    //fmt.Println("Config: ", config)

    m := martini.Classic()
    m.Get("/", func() string {
        return "Hello there rows found"
    })
    m.Run()
}
