package helper

import (
	"fmt"
	constants "instix_auth/constants"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetImagesURL(fieldName string, pid string, c *gin.Context) (ImagesURL []string, uploadedImages int) {

	form, _ := c.MultipartForm()
	files := form.File[fieldName]

	for i, file := range files {
		if i > constants.MaxItemImages {
			break
		}

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

	return ImagesURL, constants.MaxItemImages + 1
}

func UpdateImagesURL(removed []string, images []*multipart.FileHeader, pid string, c *gin.Context) (ImagesURL []string) {

	prefixURL := domName + constants.ProductImageURL + "/" + pid
	prefixURLLength := len(prefixURL + "/")
	prefixProductsDir := constants.ProductImageDir + "/" + pid + "/"

	for _, image := range removed {
		filename := image[prefixURLLength:]
		filepath := prefixProductsDir + filename
		err := os.Remove(filepath)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for i, image := range images {
		uniqueId := uuid.New()
		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := strings.Split(image.Filename, ".")[1]
		file := filename + "." + fileExt
		err := c.SaveUploadedFile(image, fmt.Sprintf("%s/%s/%s", constants.ProductImageDir, pid, file))
		if err != nil {
			fmt.Println("Error in saving Image :", i+1)
			break
		}
	}

	productDirectory, err := os.Open(constants.ProductImageDir + "/" + pid)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer productDirectory.Close()
	files, err := productDirectory.Readdir(-1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		ImagesURL = append(ImagesURL, prefixURL+"/"+file.Name())
	}

	return ImagesURL
}
