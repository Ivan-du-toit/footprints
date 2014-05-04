package main

import (
    "encoding/json"
    "fmt"
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

type PageService struct {
    pages map[string]Page
}

func (page PageService) Register() {
    ws := new(restful.WebService)
    ws.
        Path("/pages").
        Consumes(restful.MIME_XML, restful.MIME_JSON).
        Produces(restful.MIME_JSON, restful.MIME_XML)

    ws.Route(ws.GET("/{page-id}").To(page.FindPage).
        //This is just docs
        Doc("get a page").
        Operation("findPage").
        Param(ws.PathParameter("page-id", "identifier of the page").DataType("integer")).
        Writes(Page{}))

    restful.Add(ws)
}

func (ps PageService) FindPage(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("page-id")
    page := ps.pages[id]

    if page.Id == 0 {
        response.WriteErrorString(http.StatusNotFound, "Page could not be found.")
    } else {
        response.WriteEntity(page)
    }
}

func dbSetup(config Config, pageService PageService) {
    db, err := sqlx.Open("postgres", "user="+config.DBUser+
        " dbname="+config.DBName+
        " password="+config.DBPassword+" host=localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    pages := []Page{}
    err = db.Select(&pages, "SELECT * FROM page")
    if err != nil {
        fmt.Printf("err", err)
        return
    }

    pageService.pages["1"] = pages[0]
    fmt.Print("%#v\n", pageService.pages["1"])
}

func main() {
    pageService := PageService{map[string]Page{}}
    pageService.Register()

    var config Config
    configFile, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal("opening config file. ", err.Error())
    }

    err = json.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatal("Error parsing config file. ", err.Error())
    }

    swaggerConfig := swagger.Config{
        WebServices:    restful.RegisteredWebServices(), // you control what services are visible
        WebServicesUrl: "http://localhost:8080",
        ApiPath:        "/apidocs.json",
        // Optionally, specifiy where the UI is located
        SwaggerPath:     "/apidocs/",
        SwaggerFilePath: "swagger"}
    swagger.InstallSwaggerService(swaggerConfig)

    dbSetup(config, pageService)

    log.Printf("start listening on localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
