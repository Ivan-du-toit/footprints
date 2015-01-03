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
    "io"
    "os"
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

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func (page Page) Validate() (bool) {
    return page.Url != "" && page.Content != ""
}

func loadConfig() Config {
    var config Config
    configFile, err := ioutil.ReadFile("config.json")
    if err != nil {
        Error.Fatal("opening config file. ", err.Error())
    }

    err = json.Unmarshal(configFile, &config)
    if err != nil {
        Error.Fatal("Error parsing config file. ", err.Error())
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

func initLogs(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {

    Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Global Filter
func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    Info.Printf("[global-filter (logger)] %s,%s\n", req.Request.Method, req.Request.URL)
    chain.ProcessFilter(req, resp)
}

func main() {
    initLogs(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    config := loadConfig()

    db := dbSetup(config)

    restful.Filter(globalLogging)
    restful.DefaultContainer.EnableContentEncoding(true)

    restful.DefaultResponseContentType(restful.MIME_JSON)
    restful.DefaultRequestContentType(restful.MIME_JSON)

    pageService := PageService{db}
    pageService.Register()

    searchService := SearchService{db}
    searchService.Register()

    initSwagger()

    Info.Printf("start listening on localhost:8888")
    Info.Fatal(http.ListenAndServe("192.168.1.102:8888", nil))
}
