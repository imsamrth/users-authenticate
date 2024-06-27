package helper

import (
	"fmt"
	"instix_auth/constants"
	"instix_auth/database"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	domName       = os.Getenv("DomainName")
	bucketHandle  = database.BucketHandle
	bucket        = os.Getenv("BUCKET")
	storageDomain = os.Getenv("StorageDomain")
)

func GetImageURL(fh *multipart.FileHeader, id string, fp string, sp string, c *gin.Context) (ImageURL string) {

	fmt.Println("Filename is ", fh.Filename)
	var fileExt string

	filenameArr := strings.Split(fh.Filename, ".")
	if len(filenameArr) > 1 {
		fileExt = filenameArr[1]
	} else {
		return constants.IMAGE_NOT_UPLOADED
	}
	image := fmt.Sprintf("%s.%s", id, fileExt)
	filepath := fmt.Sprintf("%s/%s", fp, image)

	err := os.Remove(filepath)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == fmt.Sprintf("remove %s: The system cannot find the file specified.", filepath) {
			err = nil
		}
	}
	ImageURL, err = UploadFile(fh, filepath, c)
	if err != nil {
		fmt.Println("Error in saving Image :", err)
		return constants.IMAGE_NOT_UPLOADED
	}

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

func UploadFile(fh *multipart.FileHeader, filepath string, ctx *gin.Context) (imageurl string, err error) {
	err = nil
	sw := bucketHandle.Object(filepath).NewWriter(ctx)
	read, err := fh.Open()
	if _, err = io.Copy(sw, read); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	if err = sw.Close(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	u, err := url.Parse(storageDomain + "/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	return u.String(), err
}

var objectIDFromHex = func(hex string, c gin.Context) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return objectID
}
