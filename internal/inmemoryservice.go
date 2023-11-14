package internal

import (
	"andrewsaputra/go-rest-sample/api"
	"errors"
	"net/http"
	"sync"
	"time"
)

func ConstructInMemoryService(idGen api.IdGenerator) *InMemoryService {
	return &InMemoryService{
		Albums: []api.Album{},
		IdGen:  idGen,
	}
}

type InMemoryService struct {
	Albums []api.Album
	IdGen  api.IdGenerator
	Lock   sync.RWMutex
}

func (this *InMemoryService) GetAlbums() api.HandlerResponse {
	this.Lock.RLock()
	defer this.Lock.RUnlock()

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: this.Albums},
	}
}

func (this *InMemoryService) GetAlbumById(id string) api.HandlerResponse {
	this.Lock.RLock()
	defer this.Lock.RUnlock()

	for _, v := range this.Albums {
		if v.Id == id {
			return api.HandlerResponse{
				Code: http.StatusOK,
				Body: api.ResponseBody{Data: v},
			}
		}
	}

	return api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("album data not found")}
}

func (this *InMemoryService) InsertAlbum(props api.AlbumPropertiesDTO) api.HandlerResponse {
	this.Lock.Lock()
	defer this.Lock.Unlock()

	newData := api.Album{
		Id:          this.IdGen.NextId(),
		Title:       props.Title,
		Artist:      props.Artist,
		Price:       props.Price,
		TimeCreated: time.Now().UnixMilli(),
	}
	this.Albums = append(this.Albums, newData)

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: newData, Message: "new album data created"},
	}
}

func (this *InMemoryService) ReplaceAlbum(id string, props api.AlbumPropertiesDTO) api.HandlerResponse {
	this.Lock.Lock()
	defer this.Lock.Unlock()

	for i, _ := range this.Albums {
		album := &this.Albums[i]
		if album.Id == id {
			album.Title = props.Title
			album.Artist = props.Artist
			album.Price = props.Price

			return api.HandlerResponse{
				Code: http.StatusOK,
				Body: api.ResponseBody{Data: *album, Message: "album data replaced"},
			}
		}
	}

	return api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("album data not found")}
}

func (this *InMemoryService) UpdateAlbum(id string, updates api.AlbumUpdatesDTO) api.HandlerResponse {
	this.Lock.Lock()
	defer this.Lock.Unlock()

	for i, _ := range this.Albums {
		album := &this.Albums[i]
		if album.Id == id {
			if updates.Title != "" {
				album.Title = updates.Title
			}
			if updates.Artist != "" {
				album.Artist = updates.Artist
			}
			if updates.Price > 0 {
				album.Price = updates.Price
			}

			return api.HandlerResponse{
				Code: http.StatusOK,
				Body: api.ResponseBody{Data: *album, Message: "album data updated"},
			}
		}
	}

	return api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("album data not found")}
}

func (this *InMemoryService) DeleteAlbum(id string) api.HandlerResponse {
	this.Lock.Lock()
	defer this.Lock.Unlock()

	for i, v := range this.Albums {
		if v.Id == id {
			this.Albums = append(this.Albums[:i], this.Albums[i+1:]...)

			return api.HandlerResponse{
				Code: http.StatusOK,
				Body: api.ResponseBody{Message: "album data removed"},
			}
		}
	}

	return api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("album data not found")}
}
