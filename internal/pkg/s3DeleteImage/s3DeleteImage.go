package s3DeleteImage

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DeleteFromS3 удаляет файл по URL
func DeleteFromS3(fileURL string) error {
	accessKey := os.Getenv("S3ACCESSKEY")
	secretKey := os.Getenv("S3SECRETKEY")
	endpoint := os.Getenv("S3URL")
	bucketName := os.Getenv("S3NAME")

	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("ошибка разбора URL: %w", err)
	}

	fileKey := parsedURL.Path[1:]

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ru-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: endpoint, SigningRegion: "ru-1"}, nil
		})),
	)
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Обязательно для Selectel
	})

	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return fmt.Errorf("ошибка удаления файла из S3: %w", err)
	}
	return nil
}
