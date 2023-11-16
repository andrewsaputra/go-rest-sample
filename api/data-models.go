package api

type AppConfig struct {
	DbType         string
	MongoConfig    MongoConfig
	DynamoDbConfig DynamoDbConfig
}

type MongoConfig struct {
	Hosts               []string
	Database            string
	Collection          string
	QueryTimeoutSeconds int
}

type DynamoDbConfig struct {
	LocalEndpoint       string
	TableName           string
	Region              string
	QueryTimeoutSeconds int
}

type ResponseBody struct {
	Data    any    `json:",omitempty"`
	Message string `json:",omitempty"`
}

type HandlerResponse struct {
	Code  int
	Body  ResponseBody
	Error error
}

type AlbumPropertiesDTO struct {
	Title  string  `validate:"required"`
	Artist string  `validate:"required"`
	Price  float64 `validate:"required"`
}

type AlbumUpdatesDTO struct {
	Title  string  `validate:"required_without_all=Artist Price"`
	Artist string  `validate:"required_without_all=Title Price"`
	Price  float64 `validate:"required_without_all=Title Artist"`
}

type Album struct {
	Id          string `bson:"_id"`
	Title       string
	Artist      string
	Price       float64
	TimeCreated int64
}
