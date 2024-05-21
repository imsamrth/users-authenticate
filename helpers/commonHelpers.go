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

func GetImageURL(fh *multipart.FileHeader, id string, fp string, sp string, c *gin.Context) (ImageURL string) {

	fileExt := strings.Split(fh.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", id, fileExt)

	err := c.SaveUploadedFile(fh, fmt.Sprintf("%s/%s", fp, image))
	if err != nil {
		fmt.Println("Error in saving Image :", err)
		return constants.IMAGE_NOT_UPLOADED
	}

	ImageURL = fmt.Sprintf("%s%s/%s", domName, sp, image)
	return ImageURL
}

func RemoveImage(imgURL string, preURL string, preDir string, c *gin.Context) (err error) {

	err = nil
	//prefixURL := domName + constants.ProductImageURL + "/" + pid
	//prefixURLLength := len(prefixURL + "/")
	//prefixProductsDir := constants.ProductImageDir + "/" + pid + "/"

	preURL = fmt.Sprintf("%s%s/", domName, preURL)
	fmt.Printf("PRE_URL : %s\n PRE_DIR : %s\n IMG_URL : %s\n", preURL, preDir, imgURL)
	filename := imgURL[len(preURL):]
	filepath := preDir + "/" + filename
	err = os.Remove(filepath)

	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == fmt.Sprintf("remove %s: The system cannot find the file specified.", filepath) {
			err = nil
		}
	}

	return err
}
