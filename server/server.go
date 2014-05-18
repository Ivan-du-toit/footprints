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
)

type Page struct {
    Id      int
    Title   string `db:"title"`
    Url     string
    Content string
}

type Config struct {
    DBName     string
    DBUser     string
    DBPassword string
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
        WebServicesUrl: "http://localhost:8080",
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
    initSwagger()

    db := dbSetup(config)

    pageService := PageService{db}
    pageService.Register()

    log.Printf("start listening on localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
