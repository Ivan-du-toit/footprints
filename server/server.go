package main

import (
    "encoding/json"
    "github.com/emicklei/go-restful"
    "github.com/emicklei/go-restful/swagger"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "io/ioutil"
    "log"
    "net/http"
    "time"
)

type Page struct {
    Id         int       `db:"id"`
    Title      string    `db:"title"`
    Url        string    `db:"url"`
    Content    string    `db:"content"`
    UpdateTime time.Time `db:"updatetime"`
}

type Config struct {
    DBName     string
    DBUser     string
    DBPassword string
}

func (page Page) Validate() (bool) {
    return page.Url != "" && page.Content != ""
}

func loadConfig() Config {
    var config Config
    configFile, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal("opening config file. ", err.Error())
    }

    err = json.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatal("Error parsing config file. ", err.Error())
    }
    return config
}

func initSwagger() {
    swaggerConfig := swagger.Config{
        WebServices:    restful.RegisteredWebServices(), // you control what services are visible
        WebServicesUrl: "http://localhost:8888",
        ApiPath:        "/apidocs.json",
        // Optionally, specifiy where the UI is located
        SwaggerPath:     "/apidocs/",
        SwaggerFilePath: "swagger"}
    swagger.InstallSwaggerService(swaggerConfig)
}

func dbSetup(config Config) *sqlx.DB {
    db, err := sqlx.Open("postgres", "user="+config.DBUser+
        " dbname="+config.DBName+
        " password="+config.DBPassword+" host=localhost")
    if err != nil {
        log.Fatal(err)
    }
    //defer db.Close()

    return db
}

func main() {
    config := loadConfig()

    db := dbSetup(config)

    pageService := PageService{db}
    pageService.Register()

    initSwagger()

    log.Printf("start listening on localhost:8888")
    log.Fatal(http.ListenAndServe("192.168.1.102:8888", nil))
}
