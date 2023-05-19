package linode

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/techschool/simplebank/util"
)

func (lin *Linode) UploadToBucket(ctx *gin.Context, file *multipart.FileHeader, config util.Config) (*string, error) {
	fileToUpload, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer fileToUpload.Close()
	fileBuffer := make([]byte, file.Size)
	_, err = fileToUpload.Read(fileBuffer)
	if err != nil {
		return nil, err
	}
	params := &s3.PutObjectInput{
		Bucket:      aws.String(config.BucketName),
		Key:         aws.String(file.Filename),
		Body:        bytes.NewReader(fileBuffer),
		ContentType: aws.String(http.DetectContentType(fileBuffer)),
		ACL:         aws.String(config.BucketACL),
	}
	_, err = s3.New(lin.LinodeSession).PutObject(params)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%v.%v/%v", config.BucketName, config.BucketEndpoint, file.Filename)
	return &url, nil
}
