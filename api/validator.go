package api

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/techschool/simplebank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

func fileTypeValidator(file *multipart.FileHeader) error {
	supportedFileTypes := []string{"png", "jpg", "jpeg"}
	filename := strings.Split(file.Filename, ".")
	fileType := filename[1]
	for _, currentType := range supportedFileTypes {
		if currentType == fileType {
			return nil
		}
	}
	return fmt.Errorf("please select the supported file type i.e png, jpg or jpeg")
}
