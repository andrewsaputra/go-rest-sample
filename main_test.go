package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var router *gin.Engine

func TestInitRouter_RegisterRoutes_HandlerFunctionsCalled(t *testing.T) {
	handler := new(MockHandler)
	handler.On("GetAlbums", mock.Anything).Return()
	handler.On("GetAlbumById", mock.Anything).Return()
	handler.On("InsertAlbum", mock.Anything).Return()
	handler.On("ReplaceAlbum", mock.Anything).Return()
	handler.On("UpdateAlbum", mock.Anything).Return()
	handler.On("DeleteAlbum", mock.Anything).Return()

	router := initRouter(handler)

	request, _ := http.NewRequest(http.MethodGet, "/albums", nil)
	router.ServeHTTP(httptest.NewRecorder(), request)

	request, _ = http.NewRequest(http.MethodGet, "/albums/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), request)

	request, _ = http.NewRequest(http.MethodPost, "/albums", nil)
	router.ServeHTTP(httptest.NewRecorder(), request)

	request, _ = http.NewRequest(http.MethodPut, "/albums/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), request)

	request, _ = http.NewRequest(http.MethodPatch, "/albums/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), request)

	request, _ = http.NewRequest(http.MethodDelete, "/albums/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), request)

	handler.AssertNumberOfCalls(t, "GetAlbums", 1)
	handler.AssertNumberOfCalls(t, "GetAlbumById", 1)
	handler.AssertNumberOfCalls(t, "InsertAlbum", 1)
	handler.AssertNumberOfCalls(t, "ReplaceAlbum", 1)
	handler.AssertNumberOfCalls(t, "UpdateAlbum", 1)
	handler.AssertNumberOfCalls(t, "DeleteAlbum", 1)
}

func TestStatusCheck_StatusCheckSuccess(t *testing.T) {
	router := gin.Default()
	router.GET("/status", statusCheck)

	request, _ := http.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var responseBody map[string]any
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "Healthy", responseBody["status"])
}

type MockHandler struct {
	mock.Mock
}

func (this *MockHandler) GetAlbums(c *gin.Context) {
	this.Called(c)
}

func (this *MockHandler) GetAlbumById(c *gin.Context) {
	this.Called(c)
}

func (this *MockHandler) InsertAlbum(c *gin.Context) {
	this.Called(c)
}

func (this *MockHandler) ReplaceAlbum(c *gin.Context) {
	this.Called(c)
}

func (this *MockHandler) UpdateAlbum(c *gin.Context) {
	this.Called(c)
}

func (this *MockHandler) DeleteAlbum(c *gin.Context) {
	this.Called(c)
}
