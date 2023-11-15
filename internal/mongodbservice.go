package internal

import (
	"andrewsaputra/go-rest-sample/api"
	"context"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBService(config api.AppConfig, idGen api.IdGenerator) (*MongoDBService, error) {
	timeout := time.Duration(config.MongoConfig.QueryTimeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().
			SetHosts(config.MongoConfig.Hosts),
	)
	if err != nil {
		return nil, err
	}

	collection := client.
		Database(config.MongoConfig.Database).
		Collection(config.MongoConfig.Collection)

	return &MongoDBService{
		IdGen:      idGen,
		Collection: collection,
		Timeout:    timeout,
	}, nil
}

type MongoDBService struct {
	IdGen      api.IdGenerator
	Collection *mongo.Collection
	Timeout    time.Duration
}

func (this *MongoDBService) GetAlbums() api.HandlerResponse {
	ctx, cancel := context.WithTimeout(context.Background(), this.Timeout)
	defer cancel()

	findOpts := options.Find().SetSort(bson.M{"timecreated": 1})
	cursor, err := this.Collection.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		return api.HandlerResponse{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	albums := []api.Album{}
	for cursor.Next(nil) {
		var alb api.Album
		if err := cursor.Decode(&alb); err != nil {
			return api.HandlerResponse{
				Code:  http.StatusInternalServerError,
				Error: err,
			}
		}

		albums = append(albums, alb)
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{
			Data: albums,
		},
	}
}

func (this *MongoDBService) GetAlbumById(id string) api.HandlerResponse {
	ctx, cancel := context.WithTimeout(context.Background(), this.Timeout)
	defer cancel()

	filter := bson.M{"_id": id}
	result := this.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		var code int
		switch err {
		case mongo.ErrNoDocuments:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		return api.HandlerResponse{
			Code:  code,
			Error: err,
		}
	}

	var alb api.Album
	if err := result.Decode(&alb); err != nil {
		return api.HandlerResponse{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: alb},
	}
}

func (this *MongoDBService) InsertAlbum(props api.AlbumPropertiesDTO) api.HandlerResponse {
	newData := api.Album{
		Id:          this.IdGen.NextId(),
		Title:       props.Title,
		Artist:      props.Artist,
		Price:       props.Price,
		TimeCreated: time.Now().UnixMilli(),
	}

	_, err := this.Collection.InsertOne(nil, newData)
	if err != nil {
		return api.HandlerResponse{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: newData, Message: "new album data created"},
	}
}

func (this *MongoDBService) ReplaceAlbum(id string, props api.AlbumPropertiesDTO) api.HandlerResponse {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"title":  props.Title,
		"artist": props.Artist,
		"price":  props.Price,
	}}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	result := this.Collection.FindOneAndUpdate(nil, filter, update, opts)
	if err := result.Err(); err != nil {
		var code int
		switch err {
		case mongo.ErrNoDocuments:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		return api.HandlerResponse{
			Code:  code,
			Error: err,
		}
	}

	var alb api.Album
	if err := result.Decode(&alb); err != nil {
		return api.HandlerResponse{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: alb, Message: "album data replaced"},
	}
}

func (this *MongoDBService) UpdateAlbum(id string, updates api.AlbumUpdatesDTO) api.HandlerResponse {
	updateMap := bson.M{}
	if updates.Title != "" {
		updateMap["title"] = updates.Title
	}
	if updates.Artist != "" {
		updateMap["artist"] = updates.Artist
	}
	if updates.Price > 0 {
		updateMap["price"] = updates.Price
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateMap}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	result := this.Collection.FindOneAndUpdate(nil, filter, update, opts)
	if err := result.Err(); err != nil {
		var code int
		switch err {
		case mongo.ErrNoDocuments:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		return api.HandlerResponse{
			Code:  code,
			Error: err,
		}
	}

	var alb api.Album
	if err := result.Decode(&alb); err != nil {
		return api.HandlerResponse{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: alb, Message: "album data updated"},
	}
}

func (this *MongoDBService) DeleteAlbum(id string) api.HandlerResponse {
	filter := bson.M{"_id": id}
	result, err := this.Collection.DeleteOne(nil, filter)
	if err != nil {
		return api.HandlerResponse{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	if result.DeletedCount == 0 {
		return api.HandlerResponse{
			Code:  http.StatusNotFound,
			Error: errors.New("album data not found"),
		}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Message: "album data removed"},
	}
}
