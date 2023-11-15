package internal

import (
	"andrewsaputra/go-rest-sample/api"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func InitServiceWithMocks() api.Service {
	idGen := NewXidGenerator()
	service := NewInMemoryService(idGen)
	return service
}

func TestServiceInsertAlbum_InsertSuccess_ReturnData(t *testing.T) {
	service := InitServiceWithMocks()
	props := api.AlbumPropertiesDTO{Title: "title 1", Artist: "artist 1", Price: 1.11}
	response := service.InsertAlbum(props)
	albumResp := response.Body.Data.(api.Album)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotEmpty(t, albumResp.Id)
	assert.Equal(t, props.Title, albumResp.Title)
	assert.Equal(t, props.Artist, albumResp.Artist)
	assert.Equal(t, props.Price, albumResp.Price)
	assert.NotEmpty(t, albumResp.TimeCreated)
	assert.Nil(t, response.Error)
}

func TestServiceGetAlbums_NoData_ReturnEmpty(t *testing.T) {
	service := InitServiceWithMocks()
	response := service.GetAlbums()

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, []api.Album{}, response.Body.Data)
	assert.Nil(t, response.Error)
}

func TestServiceGetAlbums_HasData_ReturnData(t *testing.T) {
	service := InitServiceWithMocks()

	albumProps := []api.AlbumPropertiesDTO{
		{Title: "title 1", Artist: "artist 1", Price: 1.11},
		{Title: "title 2", Artist: "artist 2", Price: 2.22},
	}
	for _, props := range albumProps {
		service.InsertAlbum(props)
	}

	response := service.GetAlbums()

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Nil(t, response.Error)
	for i, data := range response.Body.Data.([]api.Album) {
		expected := albumProps[i]

		assert.NotEmpty(t, data.Id)
		assert.Equal(t, expected.Title, data.Title)
		assert.Equal(t, expected.Artist, data.Artist)
		assert.Equal(t, expected.Price, data.Price)
		assert.NotEmpty(t, data.TimeCreated)
	}
}

func TestServiceGetAlbumById_NoData_ReturnErrorNotFound(t *testing.T) {
	service := InitServiceWithMocks()
	response := service.GetAlbumById("id")

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.NotNil(t, response.Error)
}

func TestServiceGetAlbumById_HasData_ReturnData(t *testing.T) {
	service := InitServiceWithMocks()

	props := api.AlbumPropertiesDTO{Title: "title 1", Artist: "artist 1", Price: 1.11}
	insertResp := service.InsertAlbum(props)
	albumResp := insertResp.Body.Data.(api.Album)

	response := service.GetAlbumById(albumResp.Id)
	respData := response.Body.Data.(api.Album)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, albumResp.Id, respData.Id)
	assert.Equal(t, props.Title, respData.Title)
	assert.Equal(t, props.Artist, respData.Artist)
	assert.Equal(t, props.Price, respData.Price)
	assert.Equal(t, albumResp.TimeCreated, respData.TimeCreated)
	assert.Nil(t, response.Error)
}

func TestServiceReplaceAlbum_NoData_ReturnErrorNotFound(t *testing.T) {
	service := InitServiceWithMocks()

	replacementProps := api.AlbumPropertiesDTO{Title: "title 2", Artist: "artist 2", Price: 2.22}
	response := service.ReplaceAlbum("id", replacementProps)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.NotNil(t, response.Error)
}

func TestServiceReplaceAlbum_HasData_ReplaceAndReturnData(t *testing.T) {
	service := InitServiceWithMocks()

	props := api.AlbumPropertiesDTO{Title: "title 1", Artist: "artist 1", Price: 1.11}
	insertResp := service.InsertAlbum(props)
	albumResp := insertResp.Body.Data.(api.Album)

	replacementProps := api.AlbumPropertiesDTO{Title: "title 2", Artist: "artist 2", Price: 2.22}
	response := service.ReplaceAlbum(albumResp.Id, replacementProps)
	respData := response.Body.Data.(api.Album)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, albumResp.Id, respData.Id)
	assert.Equal(t, replacementProps.Title, respData.Title)
	assert.Equal(t, replacementProps.Artist, respData.Artist)
	assert.Equal(t, replacementProps.Price, respData.Price)
	assert.Equal(t, albumResp.TimeCreated, respData.TimeCreated)
	assert.Nil(t, response.Error)
}

func TestServiceUpdateAlbum_NoData_ReturnErrorNotFound(t *testing.T) {
	service := InitServiceWithMocks()

	replacementProps := api.AlbumUpdatesDTO{Title: "title 2", Artist: "artist 2", Price: 2.22}
	response := service.UpdateAlbum("id", replacementProps)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.NotNil(t, response.Error)
}

func TestServiceUpdateAlbum_HasData_ReplaceAndReturnData(t *testing.T) {
	service := InitServiceWithMocks()

	props := api.AlbumPropertiesDTO{Title: "title 1", Artist: "artist 1", Price: 1.11}
	insertResp := service.InsertAlbum(props)
	albumResp := insertResp.Body.Data.(api.Album)

	replacementProps := api.AlbumUpdatesDTO{Title: "title 2", Artist: "artist 2", Price: 2.22}
	response := service.UpdateAlbum(albumResp.Id, replacementProps)
	respData := response.Body.Data.(api.Album)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, albumResp.Id, respData.Id)
	assert.Equal(t, replacementProps.Title, respData.Title)
	assert.Equal(t, replacementProps.Artist, respData.Artist)
	assert.Equal(t, replacementProps.Price, respData.Price)
	assert.Equal(t, albumResp.TimeCreated, respData.TimeCreated)
	assert.Nil(t, response.Error)
}

func TestServiceDeleteAlbum_NoData_ReturnErrorNotFound(t *testing.T) {
	service := InitServiceWithMocks()

	response := service.DeleteAlbum("id")

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.NotNil(t, response.Error)
}

func TestServiceDeleteAlbum_HasData_ReturnOK(t *testing.T) {
	service := InitServiceWithMocks()

	props := api.AlbumPropertiesDTO{Title: "title 1", Artist: "artist 1", Price: 1.11}
	insertResp := service.InsertAlbum(props)
	albumResp := insertResp.Body.Data.(api.Album)

	response := service.DeleteAlbum(albumResp.Id)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Nil(t, response.Error)

	response = service.GetAlbumById(albumResp.Id)
	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.NotNil(t, response.Error)
}
