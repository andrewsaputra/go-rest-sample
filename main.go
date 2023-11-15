package main

import (
	"andrewsaputra/go-rest-sample/api"
	"andrewsaputra/go-rest-sample/internal"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime time.Time = time.Now()

func main() {
	config, err := GetAppConfig("configs/appconfig.json")
	if err != nil {
		log.Fatal(err)
	}

	idGenerator := internal.NewXidGenerator()
	service, err := InitService(*config, idGenerator)
	if err != nil {
		log.Fatal(err)
	}

	handler := internal.NewApiHandler(service)
	router := InitRouter(handler)

	router.Run(":8080")
}

func GetAppConfig(path string) (*api.AppConfig, error) {
	rawConfig, _ := os.ReadFile(path)
	var config *api.AppConfig
	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func InitService(config api.AppConfig, idGenerator api.IdGenerator) (api.Service, error) {
	var service api.Service
	var err error

	switch config.DbType {
	case "inmemory":
		service, err = internal.NewInMemoryService(idGenerator)
	case "mongodb":
		service, err = internal.NewMongoDBService(config, idGenerator)
	default:
		return nil, errors.ErrUnsupported
	}

	if err != nil {
		return nil, err
	}

	return service, nil
}
func InitRouter(handler api.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/status", StatusCheck)

	router.GET("/albums", handler.GetAlbums)
	router.GET("/albums/:id", handler.GetAlbumById)
	router.POST("/albums", handler.InsertAlbum)
	router.PUT("/albums/:id", handler.ReplaceAlbum)
	router.PATCH("/albums/:id", handler.UpdateAlbum)
	router.DELETE("/albums/:id", handler.DeleteAlbum)

	return router
}

func StatusCheck(c *gin.Context) {
	response := map[string]any{}
	response["status"] = "Healthy"
	response["startedAt"] = startTime.Format(time.RFC1123Z)

	c.JSON(http.StatusOK, response)
}
