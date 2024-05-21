package helper

import (
	"fmt"
	"instix_auth/constants"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	domName = os.Getenv("DomainName")
)

func GetImageURL(fh *multipart.FileHeader, id string, fp string, c *gin.Context) (ImagesURL string) {

	fileExt := strings.Split(fh.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", id, fileExt)

	err := c.SaveUploadedFile(fh, fmt.Sprintf("%s/%s", fp, image))
	if err != nil {
		fmt.Println("Error in saving Image :", err)
		return constants.IMAGE_NOT_UPLOADED
	}
	return ImagesURL
}
