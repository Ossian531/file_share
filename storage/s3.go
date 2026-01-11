package storage


import (
	"context"
	"os"
    "time"
    "fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"file_share/app_config"

)




func loadMinioConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           app_config.AppConfig.S3_endpoint,
					SigningRegion: app_config.AppConfig.Bucket_region,
				}, nil
			}),
		),
	)
}

func generateS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := loadMinioConfig(ctx)

	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = app_config.AppConfig.UsePathStyle
	})
	return client, nil
}



func GeneratePresignedUploadURL(ctx context.Context, filename string) (string, error) {
    cfg, err := loadMinioConfig(ctx)

    if err != nil {
		fmt.Println(err)
        return "", err
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.UsePathStyle = app_config.AppConfig.UsePathStyle
    })
    presigner := s3.NewPresignClient(client)

    req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(app_config.AppConfig.Bucket_name),
        Key:    aws.String(filename),
    }, s3.WithPresignExpires(5*time.Minute))

    if err != nil {
		fmt.Println(err)
        return "", err
    }

    return req.URL, nil
}



func GeneratePresignedDownloadURL(ctx context.Context, filename string) (string, error) {
    cfg, err := loadMinioConfig(ctx)

    if err != nil {
		fmt.Println(err)
        return "", err
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.UsePathStyle = app_config.AppConfig.UsePathStyle
    })
    presigner := s3.NewPresignClient(client)

    req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(app_config.AppConfig.Bucket_name),
        Key:    aws.String(filename),
    }, s3.WithPresignExpires(5*time.Minute))

    if err != nil {
		fmt.Println(err)
        return "", err
    }

    return req.URL, nil
}

func ListObjects(ctx context.Context, prefix string) ([]string, error) {
	client, err := generateS3Client(ctx)
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}

    paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
        Bucket: aws.String(app_config.AppConfig.Bucket_name),
        Prefix: aws.String(prefix),
    })

    var keys []string

    for paginator.HasMorePages() {
        page, err := paginator.NextPage(ctx)

        if err != nil {
			fmt.Println(err)
            return nil, err
        }

        for _, obj := range page.Contents {
            keys = append(keys, *obj.Key)
        }
    }

    return keys, nil
}

