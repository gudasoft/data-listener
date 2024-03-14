package deststream

import (
	"buffer-handler/logging"
	"buffer-handler/models"
	"buffer-handler/utils"
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3StreamConfig struct {
	Region       string
	Bucket       string
	PrefixFormat string
	KeyFormat    string
	ObjType      string
}

func (cfg S3StreamConfig) Notify(data models.EntryData) error {
	logging.Logger.Debug("S3 Streaming")

	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	}))

	svc := s3.New(session)

	objKey := cfg.generateObjKey()

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(objKey),
		Body:   bytes.NewReader(data.Body),
	})
	if err != nil {
		logging.Logger.Sugar().Errorf("Error uploading object:", err)
		return err
	}

	logging.Logger.Debug("Object uploaded successfully!")
	return nil
}

func (cfg S3StreamConfig) generateObjKey() string {
	currentTime := time.Now()

	prefix := utils.FormatConfig(cfg.PrefixFormat, currentTime)

	key := utils.FormatConfig(cfg.KeyFormat, currentTime)

	return (prefix + "/" + key + "." + cfg.ObjType)
}

func (cfg S3StreamConfig) String() string {
	obj := cfg.generateObjKey()
	return fmt.Sprintf(
		"s3://%s/%s "+
			"Region: %s",
		cfg.Bucket, obj, cfg.Region)
}
