package helper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	iitbDomain = "iitb.ac.in"
)

func ValidateEmail(email string) (statusCode int, err error) {
	at := strings.LastIndex(email, "@")
	if at >= 0 {
		username, domain := email[:at], email[at+1:]
		fmt.Printf("Username: %s, Domain: %s\n", username, domain)
	} else {
		fmt.Printf("Error: %s is an invalid email address\n", email)
	}
	// check if disposable
	if email[at+1:] != iitbDomain {
		err = errors.New("sorry, we do not accept non-IITB email addresses")
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unaithorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("Unaithorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
