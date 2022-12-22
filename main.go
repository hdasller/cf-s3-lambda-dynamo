package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type File struct {
	Content       string
	FileName      string
	FileExtension string
}

type Register struct {
	UUID      string `json:"uuid"`
	S3URL     string `json:"s3_url"`
	CreatedAt string `json:"created_at"`
}

const (
	BUCKET_NAME   = "cf-hd-s3-upload"
	TABLE_NAME    = "upload-registers"
	JSON_EXT      = ".json"
	ErrUploadFile = "Erro ao realizar upload"
)

var (
	HEADERS = map[string]string{"Content-Type": "application/json"}
)

func HandleLambdaEvent(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	response := events.APIGatewayProxyResponse{
		Headers:    HEADERS,
		StatusCode: http.StatusCreated,
		Body:       "Registro persistido com sucesso",
	}

	err := Upload(sess, BUCKET_NAME, event.Body)
	if err != nil {
		response.Body = ErrUploadFile
		response.StatusCode = http.StatusInternalServerError
	}

	return response, err
}

func Upload(sess *session.Session, bucketName, data string) error {

	u := uuid.New()
	key := u.String() + JSON_EXT

	url, err := Save(sess, bucketName, key, data)
	if err != nil {
		return err
	}

	register := &Register{
		UUID:      uuid.NewString(),
		S3URL:     url,
		CreatedAt: time.Now().UTC().String(),
	}

	err = PersistUploadInfo(sess, TABLE_NAME, register)
	if err != nil {
		return err
	}

	return nil
}

func Save(sess *session.Session, bucketName, key string, data string) (string, error) {
	uploader := s3manager.NewUploader(sess)
	s3Response, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    &key,
		Body:   bytes.NewReader([]byte(data)),
	})

	if err != nil {
		return "", err
	}

	return s3Response.Location, nil
}

func PersistUploadInfo(sess *session.Session, tableName string, register *Register) error {
	av, err := dynamodbattribute.MarshalMap(register)
	if err != nil {
		log.Fatalf("Erro ao realizar parse de dados da tabela: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	svc := dynamodb.New(sess)
	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Erro ao inserir os registros na tablea: %s", err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
