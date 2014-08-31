package main

import (
    "github.com/emicklei/go-restful"
    "github.com/jmoiron/sqlx"
    "log"
    "net/http"
)

type PageService struct {
    db *sqlx.DB
    //pages map[string]Page
}

func (page PageService) Register() {
    ws := new(restful.WebService)
    ws.
        Path("/pages").
        Consumes(restful.MIME_XML, restful.MIME_JSON).
        Produces(restful.MIME_JSON, restful.MIME_XML)

    ws.Route(ws.GET("/{page-id}").To(page.GetPage).
        Doc("Get a page").
        Operation("GetPage").
        Param(ws.PathParameter("page-id", "Identifier of the page").DataType("integer")).
        Writes(Page{}))

    ws.Route(ws.POST("/").To(page.AddPage).
        Doc("Add a page").
        Operation("AddPage").
        Reads(Page{}).
        Writes(Page{}))

    restful.Add(ws)
}

func (ps PageService) GetPage(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("page-id")
    var page []Page
    err := ps.db.Select(&page, "SELECT * FROM pages WHERE id = $1", id)
    if err != nil {
        log.Print("Error: ", err)
        response.WriteErrorString(http.StatusInternalServerError, "While loading page " + id + " we encountered an error.")
    }
    if (len(page) == 0) {
        response.WriteErrorString(http.StatusNotFound, "Page " + id + " could not be found.")
    } else {
        response.WriteEntity(page)
    }
}

func (ps PageService) AddPage(request *restful.Request, response *restful.Response) {
    page := new(Page)
    request.ReadEntity(page)
    if (!page.Validate()) {
        log.Print("Invalid page posted", page)
        //TODO: Add beter information about what exactly was missing or wrong
        response.WriteErrorString(http.StatusBadRequest, "Invalid page data")
        return
    }
    //If url exists update
    existingPage := ps.findPageByUrl(page.Url)
    var err error
    if existingPage.Id > 0 {
        //Update should rather be moved to another table and insert the new record.
        //Not sure if this should be done in code or in the DB...
        page.Id = existingPage.Id
        _, err = ps.db.NamedExec("UPDATE pages SET content = :content WHERE id = :id", page)
    } else {
        log.Print("Inserting new page: ", page)
        _, err = ps.db.NamedExec("INSERT INTO pages (content, url) VALUES (:content, :url)", page)
    }
    if err != nil {
        log.Print("Error: ", err)
        response.WriteErrorString(http.StatusInternalServerError, "Failed to save page.")
    } else {
        response.WriteHeader(http.StatusCreated)
        //TODO: Get the page because this page object does not have the ID
        response.WriteEntity(page)
    }
}

//Helper functions should be moved to a model class that handles all the persistance logic
func (ps PageService) findPageByUrl(url string) (page Page) {
    var pages []Page
    err := ps.db.Select(&pages, "SELECT * FROM pages WHERE url = $1", url)
    if err != nil {
        log.Print("Error: ", err)
        return Page{}
    }
    if (len(pages) == 0) {
        return Page{}
    }
    return pages[0]
}
