package destbuffer

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

type S3BufferConfig struct {
	Region        string
	Bucket        string
	PrefixFormat  string
	KeyFormat     string
	ObjType       string
	ItemSeparator string
}

func (cfg S3BufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("S3 Buffering")

	if len(entries) == 0 {
		return nil
	}

	var combinedBody bytes.Buffer

	for _, entry := range entries {

		combinedBody.Write(entry.Body)
		combinedBody.WriteString(cfg.ItemSeparator)
	}

	objKey := cfg.generateObjKey()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		logging.Logger.Error(err.Error())
		return err
	}

	s3Client := s3.New(sess)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(objKey),
		Body:   bytes.NewReader(combinedBody.Bytes()),
	})
	if err != nil {
		logging.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (cfg S3BufferConfig) generateObjKey() string {
	currentTime := time.Now()

	prefix := utils.FormatConfig(cfg.PrefixFormat, currentTime)

	key := utils.FormatConfig(cfg.KeyFormat, currentTime)

	return (prefix + "/" + key + "." + cfg.ObjType)
}

func (cfg S3BufferConfig) String() string {
	obj := cfg.generateObjKey()
	return fmt.Sprintf(
		"s3://%s/%s "+
			"Region: %s",
		cfg.Bucket, obj, cfg.Region)
}
