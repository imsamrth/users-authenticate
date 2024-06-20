package controllers

import (
	"fmt"
	"instix_auth/constants"
	helper "instix_auth/helpers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/xuri/excelize/v2"
)

func UploadSheet() gin.HandlerFunc {

	return func(c *gin.Context) {

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		file, _ := c.FormFile("file")
		fp := constants.StaticFile

		err := os.Remove(fp)
		if err != nil {
			fmt.Println(err.Error())
			if err.Error() == fmt.Sprintf("remove %s: The system cannot find the file specified.", fp) {
				err = nil

			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		err = c.SaveUploadedFile(file, fp)
		if err != nil {
			fmt.Println("Error in saving Image :", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "File uploaded successfully"})
	}
}

func GetColumn() gin.HandlerFunc {

	return func(c *gin.Context) {

		sheet := c.Param("sheet")

		file, err := excelize.OpenFile(constants.StaticFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string][]string{}
		cols, err := file.GetCols(sheet)
		for _, col := range cols {
			if col[0] == "" {
				continue
			}
			response[col[0]] = col[1:]
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
