package api

import "github.com/gin-gonic/gin"

type IdGenerator interface {
	NextId() string
}

type Handler interface {
	GetAlbums(c *gin.Context)
	GetAlbumById(c *gin.Context)
	InsertAlbum(c *gin.Context)
	ReplaceAlbum(c *gin.Context)
	UpdateAlbum(c *gin.Context)
	DeleteAlbum(c *gin.Context)
}

type Service interface {
	GetAlbums() HandlerResponse
	GetAlbumById(id string) HandlerResponse
	InsertAlbum(props AlbumPropertiesDTO) HandlerResponse
	ReplaceAlbum(id string, props AlbumPropertiesDTO) HandlerResponse
	UpdateAlbum(id string, updates AlbumUpdatesDTO) HandlerResponse
	DeleteAlbum(id string) HandlerResponse
}
