package aws

import (
	"context"

	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials"
)

type AWS struct {
	S3 *S3Client
}

func NewAWS(server *server.Server) (*AWS, error) {
	awsConfig := server.Config.AWS

	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(awsConfig.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsConfig.AccessKeyID,
			awsConfig.SecretAccessKey,
			"",
		)),
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return nil, err
	}

	s3Client := NewS3Client(server, cfg)

	return &AWS{
		S3: s3Client,
	}, nil
}
