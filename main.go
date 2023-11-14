package main

import (
	"andrewsaputra/go-rest-sample/api"
	"andrewsaputra/go-rest-sample/internal"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime time.Time = time.Now()

func main() {
	idGenerator := internal.ConstructXidGenerator()
	service := internal.ConstructInMemoryService(idGenerator)
	handler := internal.ConstructHandler(service)
	router := initRouter(handler)

	router.Run(":8080")
}

func initRouter(handler api.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/status", statusCheck)

	router.GET("/albums", handler.GetAlbums)
	router.GET("/albums/:id", handler.GetAlbumById)
	router.POST("/albums", handler.InsertAlbum)
	router.PUT("/albums/:id", handler.ReplaceAlbum)
	router.PATCH("/albums/:id", handler.UpdateAlbum)
	router.DELETE("/albums/:id", handler.DeleteAlbum)

	return router
}

func statusCheck(c *gin.Context) {
	response := map[string]any{}
	response["status"] = "Healthy"
	response["startedAt"] = startTime.Format(time.RFC1123Z)

	c.JSON(http.StatusOK, response)
}
