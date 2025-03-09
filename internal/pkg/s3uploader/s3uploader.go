package s3uploader

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadToS3(file multipart.File, filename string) (string, error) {
	accessKey := os.Getenv("S3ACCESSKEY")
	secretKey := os.Getenv("S3SECRETKEY")
	endpoint := os.Getenv("S3URL") // Selectel S3 URL
	region := os.Getenv("REGION")
	name := os.Getenv("S3NAME")
	imgUrl := os.Getenv("IMGURL")

	// Задаем конфигурацию с явными ключами
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})),
	)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Обязательно для Selectel
	})

	filename = primitive.NewObjectID().Hex() + filename

	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(name),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		return "", err
	}
	return imgUrl + "/" + filename, nil

}
