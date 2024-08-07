package helper

import (
	"fmt"
	"instix_auth/constants"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	domName = os.Getenv("DomainName")
)

func GetImageURL(fh *multipart.FileHeader, id string, fp string, sp string, c *gin.Context) (ImageURL string) {

	fileExt := strings.Split(fh.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", id, fileExt)
	filepath := fmt.Sprintf("%s/%s", fp, image)

	err := os.Remove(filepath)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == fmt.Sprintf("remove %s: The system cannot find the file specified.", filepath) {
			err = nil
		}
	}
	err = c.SaveUploadedFile(fh, filepath)
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

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Hash is incorrect")
		check = false
	}
	return check, msg
}
