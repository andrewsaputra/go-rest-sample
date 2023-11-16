package internal

import (
	"andrewsaputra/go-rest-sample/api"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

/*
CLI command for local table creation :
aws dynamodb create-table \
--endpoint-url http://localhost:8000 \
--table-name albums \
--billing-mode PAY_PER_REQUEST \
--attribute-definitions AttributeName=Id,AttributeType=S AttributeName=TimeCreated,AttributeType=N AttributeName=Artist,AttributeType=S \
--key-schema AttributeName=Id,KeyType=HASH \
--global-secondary-indexes '[{"IndexName":"gsi_id_timecreated","KeySchema":[{"AttributeName":"Id","KeyType":"HASH"},{"AttributeName":"TimeCreated","KeyType":"RANGE"}],"Projection":{"ProjectionType":"ALL"}},{"IndexName":"gsi_artist_timecreated","KeySchema":[{"AttributeName":"Artist","KeyType":"HASH"},{"AttributeName":"TimeCreated","KeyType":"RANGE"}],"Projection":{"ProjectionType":"ALL"}}]'

*/

func NewDynamoDbService(config api.DynamoDbConfig, idGen api.IdGenerator) (api.Service, error) {
	timeout := time.Duration(config.QueryTimeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(config.Region))
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if config.LocalEndpoint != "" {
			o.BaseEndpoint = aws.String(config.LocalEndpoint)
		}
	})

	return &DynamoDbService{
		IdGen:     idGen,
		Client:    client,
		TableName: config.TableName,
		Timeout:   timeout,
	}, nil
}

type DynamoDbService struct {
	IdGen     api.IdGenerator
	Client    *dynamodb.Client
	TableName string
	Timeout   time.Duration
}

func (this *DynamoDbService) GetAlbums() api.HandlerResponse {
	ctx, cancel := context.WithTimeout(context.Background(), this.Timeout)
	defer cancel()

	params := dynamodb.ScanInput{
		TableName:      aws.String(this.TableName),
		ConsistentRead: aws.Bool(false),
	}
	res, err := this.Client.Scan(ctx, &params)
	if err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	albums := []api.Album{}
	for _, v := range res.Items {
		var alb api.Album
		if err := attributevalue.UnmarshalMap(v, &alb); err != nil {
			return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
		}

		albums = append(albums, alb)
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: albums},
	}
}

func (this *DynamoDbService) GetAlbumById(id string) api.HandlerResponse {
	ctx, cancel := context.WithTimeout(context.Background(), this.Timeout)
	defer cancel()

	params := dynamodb.GetItemInput{
		TableName: aws.String(this.TableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	}

	res, err := this.Client.GetItem(ctx, &params)
	if err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	if len(res.Item) == 0 {
		return api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("album data not found")}
	}

	var alb api.Album
	if err := attributevalue.UnmarshalMap(res.Item, &alb); err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: alb},
	}
}

func (this *DynamoDbService) InsertAlbum(props api.AlbumPropertiesDTO) api.HandlerResponse {
	newData := api.Album{
		Id:          this.IdGen.NextId(),
		Title:       props.Title,
		Artist:      props.Artist,
		Price:       props.Price,
		TimeCreated: time.Now().UnixMilli(),
	}

	item, err := attributevalue.MarshalMap(newData)
	if err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	params := dynamodb.PutItemInput{
		TableName: aws.String(this.TableName),
		Item:      item,
	}
	if _, err := this.Client.PutItem(context.Background(), &params); err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{
			Data:    newData,
			Message: "new album data created",
		},
	}
}

func (this *DynamoDbService) ReplaceAlbum(id string, props api.AlbumPropertiesDTO) api.HandlerResponse {
	update := expression.
		Set(expression.Name("Title"), expression.Value(props.Title)).
		Set(expression.Name("Artist"), expression.Value(props.Artist)).
		Set(expression.Name("Price"), expression.Value(props.Price))

	expr, err := expression.NewBuilder().
		WithUpdate(update).
		WithCondition(expression.AttributeExists(expression.Name("Id"))).
		Build()
	if err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	params := dynamodb.UpdateItemInput{
		TableName: aws.String(this.TableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValueAllNew,
	}
	result, err := this.Client.UpdateItem(context.Background(), &params)
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("album data not found")}
		}

		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	var alb api.Album
	if err := attributevalue.UnmarshalMap(result.Attributes, &alb); err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: alb, Message: "album data replaced"},
	}
}

func (this *DynamoDbService) UpdateAlbum(id string, updates api.AlbumUpdatesDTO) api.HandlerResponse {
	var update expression.UpdateBuilder
	if updates.Title != "" {
		update = update.Set(expression.Name("Title"), expression.Value(updates.Title))
	}
	if updates.Artist != "" {
		update = update.Set(expression.Name("Artist"), expression.Value(updates.Artist))
	}
	if updates.Price > 0 {
		update = update.Set(expression.Name("Price"), expression.Value(updates.Price))
	}

	expr, err := expression.NewBuilder().
		WithUpdate(update).
		WithCondition(expression.AttributeExists(expression.Name("Id"))).
		Build()
	if err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	params := dynamodb.UpdateItemInput{
		TableName: aws.String(this.TableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValueAllNew,
	}
	result, err := this.Client.UpdateItem(context.Background(), &params)
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return api.HandlerResponse{Code: http.StatusNotFound, Error: errors.New("album data not found")}
		}

		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	var alb api.Album
	if err := attributevalue.UnmarshalMap(result.Attributes, &alb); err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Data: alb, Message: "album data updated"},
	}
}

func (this *DynamoDbService) DeleteAlbum(id string) api.HandlerResponse {
	params := dynamodb.DeleteItemInput{
		TableName: aws.String(this.TableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	}

	if _, err := this.Client.DeleteItem(context.Background(), &params); err != nil {
		return api.HandlerResponse{Code: http.StatusInternalServerError, Error: err}
	}

	return api.HandlerResponse{
		Code: http.StatusOK,
		Body: api.ResponseBody{Message: "album data removal processed"},
	}
}
