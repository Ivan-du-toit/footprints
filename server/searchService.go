package main

import (
    "github.com/emicklei/go-restful"
    "github.com/jmoiron/sqlx"
    "net/http"
    "strings"
)

type SearchService struct {
    db *sqlx.DB
    //pages map[string]Page
}

func (service SearchService) Register() {
    ws := new(restful.WebService)
    ws.
        Path("/search").
        Consumes(restful.MIME_XML, restful.MIME_JSON).
        //Produces(restful.MIME_JSON, restful.MIME_XML)
        Produces(restful.MIME_JSON)

    ws.Route(ws.GET("/{query}").To(service.Search).
        Doc("Search pages").
        Operation("Search").
        Param(ws.PathParameter("query", "Query to search on.").DataType("string")).
        Writes(PageList{}))

    restful.Add(ws)
}

func splitQuery(queryTerms string) string {
    return "%" + strings.Join(strings.Split(queryTerms, " "), "%") + "%"
}

func (service SearchService) Search(request *restful.Request, response *restful.Response) {
    query := request.PathParameter("query")
    //size := request.QueryParameter("size")
    //start := request.QueryParameter("start")
    Info.Print("Searching for phrases: ", query)
    var page []Page
    err := service.db.Select(&page, "SELECT id, content title, url FROM pages WHERE content LIKE $1 OR url LIKE $1 LIMIT 10", splitQuery(query))
    if err != nil {
        Error.Print(err)
        response.WriteErrorString(http.StatusInternalServerError, "While searching for " + query + " we encountered an error.")
        return
    }
    Info.Print("Count of pages found: ", len(page))
    if (len(page) == 0) {
        response.WriteErrorString(http.StatusNotFound, "No pages match " + query)
    } else {
        result := PageList{page, 0, len(page)}
        response.WriteEntity(result)
    }
}
