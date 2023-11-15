package internal

import (
	"andrewsaputra/go-rest-sample/api"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewApiHandler(service api.Service) *ApiHandler {
	var propsFields []string
	for _, field := range reflect.VisibleFields(reflect.TypeOf(api.AlbumPropertiesDTO{})) {
		propsFields = append(propsFields, field.Name)
	}

	return &ApiHandler{
		Service:          service,
		Validator:        validator.New(validator.WithRequiredStructEnabled()),
		AlbumPropsFields: propsFields,
	}
}

type ApiHandler struct {
	Service          api.Service
	Validator        *validator.Validate
	AlbumPropsFields []string
}

func (this *ApiHandler) GetAlbums(c *gin.Context) {
	resp := this.Service.GetAlbums()
	this.HandleResponse(c, resp)
}

func (this *ApiHandler) GetAlbumById(c *gin.Context) {
	id := c.Param("id")
	resp := this.Service.GetAlbumById(id)
	this.HandleResponse(c, resp)
}

func (this *ApiHandler) InsertAlbum(c *gin.Context) {
	var props api.AlbumPropertiesDTO
	if err := c.BindJSON(&props); err != nil {
		this.HandleResponse(c, api.HandlerResponse{Code: http.StatusBadRequest, Error: err})
		return
	}

	if err := this.Validator.Struct(props); err != nil {
		this.HandleResponse(c, api.HandlerResponse{Code: http.StatusBadRequest, Error: err})
		return
	}

	resp := this.Service.InsertAlbum(props)
	this.HandleResponse(c, resp)
}

func (this *ApiHandler) ReplaceAlbum(c *gin.Context) {
	id := c.Param("id")
	var props api.AlbumPropertiesDTO
	if err := c.BindJSON(&props); err != nil {
		this.HandleResponse(c, api.HandlerResponse{Code: http.StatusBadRequest, Error: err})
		return
	}

	if err := this.Validator.Struct(props); err != nil {
		this.HandleResponse(c, api.HandlerResponse{Code: http.StatusBadRequest, Error: err})
		return
	}

	resp := this.Service.ReplaceAlbum(id, props)
	this.HandleResponse(c, resp)
}

func (this *ApiHandler) UpdateAlbum(c *gin.Context) {
	id := c.Param("id")
	var updates api.AlbumUpdatesDTO
	if err := c.BindJSON(&updates); err != nil {
		this.HandleResponse(c, api.HandlerResponse{Code: http.StatusBadRequest, Error: err})
		return
	}

	if err := this.Validator.Struct(updates); err != nil {
		this.HandleResponse(c, api.HandlerResponse{Code: http.StatusBadRequest, Error: err})
		return
	}

	resp := this.Service.UpdateAlbum(id, updates)
	this.HandleResponse(c, resp)
}

func (this *ApiHandler) DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	resp := this.Service.DeleteAlbum(id)
	this.HandleResponse(c, resp)
}

func (this *ApiHandler) HandleResponse(c *gin.Context, resp api.HandlerResponse) {
	if resp.Error != nil {
		c.JSON(resp.Code, gin.H{"message": resp.Error.Error()})
		return
	}

	c.JSON(resp.Code, resp.Body)
}
