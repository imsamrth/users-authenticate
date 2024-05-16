package helper

import (
	"fmt"
	"instix_auth/constants"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	domName = os.Getenv("DomainName")
)

func GetImageURL(fieldName string, pid string, c *gin.Context) (ImagesURL []string, uploadedImages int) {

	form, _ := c.MultipartForm()
	files := form.File[fieldName]

	for i, file := range files {
		fmt.Println(file.Filename)
		fileExt := strings.Split(file.Filename, ".")[1]
		image := fmt.Sprintf("%d.%s", i, fileExt)
		// Upload the file to specific dst.

		err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s/%s", constants.ProductImageDir, pid, image))
		if err != nil {
			fmt.Println("Error in saving Image :", err)
			return ImagesURL, i
		}
		ImagesURL = append(ImagesURL, fmt.Sprintf("%s%s/%s/%s", domName, constants.ProductImageURL, pid, image))
	}

	return ImagesURL, 10
}
