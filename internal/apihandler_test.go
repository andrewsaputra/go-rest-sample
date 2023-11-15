package internal

import (
	"andrewsaputra/go-rest-sample/api"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerGetAlbums_ServiceReturnOK_ReturnOK(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	expectedResponse := api.HandlerResponse{Code: http.StatusOK, Body: api.ResponseBody{Message: "message"}}
	service.On("GetAlbums").Return(expectedResponse)

	handler.GetAlbums(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "GetAlbums", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Body.Message, respBody.Message)
}

func TestHandlerGetAlbums_ServiceReturnError_ReturnError(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	expectedResponse := api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("sample error")}
	service.On("GetAlbums").Return(expectedResponse)

	handler.GetAlbums(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "GetAlbums", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Error.Error(), respBody.Message)
}

func TestHandlerGetAlbumById_ServiceReturnOK_ReturnOK(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	expectedResponse := api.HandlerResponse{Code: http.StatusOK, Body: api.ResponseBody{Message: "message"}}
	service.On("GetAlbumById", mock.Anything).Return(expectedResponse)

	handler.GetAlbumById(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "GetAlbumById", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Body.Message, respBody.Message)
}

func TestHandlerGetAlbumById_ServiceReturnError_ReturnError(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	expectedResponse := api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("sample error")}
	service.On("GetAlbumById", mock.Anything).Return(expectedResponse)

	handler.GetAlbumById(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "GetAlbumById", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Error.Error(), respBody.Message)

}

func TestHandlerInsertAlbum_ServiceReturnOK_ReturnOK(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	requestDto, _ := json.Marshal(api.AlbumPropertiesDTO{Title: "title", Artist: "artist", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestDto)))

	expectedResponse := api.HandlerResponse{Code: http.StatusOK, Body: api.ResponseBody{Message: "message"}}
	service.On("InsertAlbum", mock.Anything).Return(expectedResponse)

	handler.InsertAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "InsertAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Body.Message, respBody.Message)
}

func TestHandlerInsertAlbum_ServiceReturnError_ReturnError(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	requestDto, _ := json.Marshal(api.AlbumPropertiesDTO{Title: "title", Artist: "artist", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestDto)))

	expectedResponse := api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("sample error")}
	service.On("InsertAlbum", mock.Anything).Return(expectedResponse)

	handler.InsertAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "InsertAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Error.Error(), respBody.Message)
}

func TestHandlerInsertAlbum_InvalidRequests_ReturnBadRequest(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()

	//nil payload
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	handler.InsertAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//non json payload
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte("invalid json"))))
	handler.InsertAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//empty json payload
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte("{}"))))
	handler.InsertAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 1
	requestDto, _ := json.Marshal(api.AlbumPropertiesDTO{Title: "", Artist: "", Price: 0})
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.InsertAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 2
	requestDto, _ = json.Marshal(api.AlbumPropertiesDTO{Title: "value", Artist: "", Price: 0})
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.InsertAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 3
	requestDto, _ = json.Marshal(api.AlbumPropertiesDTO{Title: "", Artist: "value", Price: 0})
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.InsertAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 4
	requestDto, _ = json.Marshal(api.AlbumPropertiesDTO{Title: "", Artist: "", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.InsertAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	service.AssertNumberOfCalls(t, "InsertAlbum", 0)
}

func TestHandlerReplaceAlbum_ServiceReturnOK_ReturnOK(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	requestDto, _ := json.Marshal(api.AlbumPropertiesDTO{Title: "title", Artist: "artist", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader(requestDto)))

	expectedResponse := api.HandlerResponse{Code: http.StatusOK, Body: api.ResponseBody{Message: "message"}}
	service.On("ReplaceAlbum", mock.Anything, mock.Anything).Return(expectedResponse)

	handler.ReplaceAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "ReplaceAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Body.Message, respBody.Message)
}

func TestHandlerReplaceAlbum_ServiceReturnError_ReturnError(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	requestDto, _ := json.Marshal(api.AlbumPropertiesDTO{Title: "title", Artist: "artist", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader(requestDto)))

	expectedResponse := api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("sample error")}
	service.On("ReplaceAlbum", mock.Anything, mock.Anything).Return(expectedResponse)

	handler.ReplaceAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "ReplaceAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Error.Error(), respBody.Message)
}

func TestHandlerReplaceAlbum_InvalidRequests_ReturnBadRequest(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()

	//nil payload
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", nil)
	handler.ReplaceAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//non json payload
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader([]byte("invalid json"))))
	handler.ReplaceAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//empty json payload
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader([]byte("{}"))))
	handler.ReplaceAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 1
	requestDto, _ := json.Marshal(api.AlbumPropertiesDTO{Title: "", Artist: "", Price: 0})
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.ReplaceAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 2
	requestDto, _ = json.Marshal(api.AlbumPropertiesDTO{Title: "value", Artist: "", Price: 0})
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.ReplaceAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 3
	requestDto, _ = json.Marshal(api.AlbumPropertiesDTO{Title: "", Artist: "value", Price: 0})
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.ReplaceAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload - 4
	requestDto, _ = json.Marshal(api.AlbumPropertiesDTO{Title: "", Artist: "", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPut, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.ReplaceAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	service.AssertNumberOfCalls(t, "ReplaceAlbum", 0)
}

func TestHandlerUpdateAlbum_ServiceReturnOK_ReturnOK(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	requestDto, _ := json.Marshal(api.AlbumUpdatesDTO{Title: "title", Artist: "artist", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPatch, "/", io.NopCloser(bytes.NewReader(requestDto)))

	expectedResponse := api.HandlerResponse{Code: http.StatusOK, Body: api.ResponseBody{Message: "message"}}
	service.On("UpdateAlbum", mock.Anything, mock.Anything).Return(expectedResponse)

	handler.UpdateAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "UpdateAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Body.Message, respBody.Message)
}

func TestHandlerUpdateAlbum_ServiceReturnError_ReturnError(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	requestDto, _ := json.Marshal(api.AlbumUpdatesDTO{Title: "title", Artist: "artist", Price: 9.99})
	ginContext.Request, _ = http.NewRequest(http.MethodPatch, "/", io.NopCloser(bytes.NewReader(requestDto)))

	expectedResponse := api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("sample error")}
	service.On("UpdateAlbum", mock.Anything, mock.Anything).Return(expectedResponse)

	handler.UpdateAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "UpdateAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Error.Error(), respBody.Message)
}

func TestHandlerUpdateAlbum_InvalidRequests_ReturnBadRequest(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()

	//nil payload
	ginContext.Request, _ = http.NewRequest(http.MethodPatch, "/", nil)
	handler.UpdateAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//non json payload
	ginContext.Request, _ = http.NewRequest(http.MethodPatch, "/", io.NopCloser(bytes.NewReader([]byte("invalid json"))))
	handler.UpdateAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//empty json payload
	ginContext.Request, _ = http.NewRequest(http.MethodPatch, "/", io.NopCloser(bytes.NewReader([]byte("{}"))))
	handler.UpdateAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	//incomplete DTO payload
	requestDto, _ := json.Marshal(api.AlbumUpdatesDTO{Title: "", Artist: "", Price: 0})
	ginContext.Request, _ = http.NewRequest(http.MethodPatch, "/", io.NopCloser(bytes.NewReader(requestDto)))
	handler.UpdateAlbum(ginContext)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	service.AssertNumberOfCalls(t, "UpdateAlbum", 0)
}

func TestHandlerDeleteAlbum_ServiceReturnOK_ReturnOK(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	ginContext.Request, _ = http.NewRequest(http.MethodDelete, "/", nil)

	expectedResponse := api.HandlerResponse{Code: http.StatusOK, Body: api.ResponseBody{Message: "message"}}
	service.On("DeleteAlbum", mock.Anything, mock.Anything).Return(expectedResponse)

	handler.DeleteAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "DeleteAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Body.Message, respBody.Message)
}

func TestHandlerDeleteAlbum_ServiceReturnError_ReturnError(t *testing.T) {
	handler, service, ginContext, respWriter := InitHandlerWithMocks()
	ginContext.Request, _ = http.NewRequest(http.MethodDelete, "/", nil)

	expectedResponse := api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("sample error")}
	service.On("DeleteAlbum", mock.Anything, mock.Anything).Return(expectedResponse)

	handler.DeleteAlbum(ginContext)

	var respBody api.ResponseBody
	json.Unmarshal(respWriter.Body.Bytes(), &respBody)

	service.AssertNumberOfCalls(t, "DeleteAlbum", 1)
	assert.Equal(t, expectedResponse.Code, respWriter.Code)
	assert.Equal(t, expectedResponse.Error.Error(), respBody.Message)
}

func InitHandlerWithMocks() (api.Handler, *MockService, *gin.Context, *httptest.ResponseRecorder) {
	respWriter := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(respWriter)
	service := new(MockService)
	handler := NewApiHandler(service)
	return handler, service, context, respWriter
}

type MockService struct {
	mock.Mock
}

func (t *MockService) GetAlbums() api.HandlerResponse {
	args := t.Called()
	return args.Get(0).(api.HandlerResponse)
}

func (t *MockService) GetAlbumById(id string) api.HandlerResponse {
	args := t.Called(id)
	return args.Get(0).(api.HandlerResponse)
}

func (t *MockService) InsertAlbum(props api.AlbumPropertiesDTO) api.HandlerResponse {
	args := t.Called(props)
	return args.Get(0).(api.HandlerResponse)
}

func (t *MockService) ReplaceAlbum(id string, props api.AlbumPropertiesDTO) api.HandlerResponse {
	args := t.Called(id, props)
	return args.Get(0).(api.HandlerResponse)
}

func (t *MockService) UpdateAlbum(id string, updates api.AlbumUpdatesDTO) api.HandlerResponse {
	args := t.Called(id, updates)
	return args.Get(0).(api.HandlerResponse)
}

func (t *MockService) DeleteAlbum(id string) api.HandlerResponse {
	args := t.Called(id)
	return args.Get(0).(api.HandlerResponse)
}
