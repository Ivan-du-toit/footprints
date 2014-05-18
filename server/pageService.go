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
    var page []Page
    err := ps.db.Select(&page, "SELECT * FROM page WHERE id = $1", 1)
    if err != nil {
        log.Print("Error: ", err)
        response.WriteErrorString(http.StatusNotFound, "Page "+id+" could not be found.")
    }
    response.WriteEntity(page[0])
}
