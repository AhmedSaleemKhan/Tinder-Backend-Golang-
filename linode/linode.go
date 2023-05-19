package linode

import (
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/techschool/simplebank/util"
)

type LinodeImpl interface {
	UploadToBucket(ctx *gin.Context, file *multipart.FileHeader, config util.Config) (*string, error)
}

type Linode struct {
	LinodeSession *session.Session
}

func NewLinodeSession(config util.Config) *Linode {
	cfg := &aws.Config{
		Credentials: credentials.NewStaticCredentials(config.BucketAccessKey, config.BucketSecretKey, ""),
		Endpoint:    aws.String(config.BucketEndpoint),
		Region:      aws.String(config.LinodeRegion),
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		log.Fatalf("failed to create linode session %v", err.Error())
	}
	return &Linode{
		LinodeSession: sess,
	}
}
