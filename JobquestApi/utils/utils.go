package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)


func HashPassword() {

}

func VerifyPassword() {}

func MatchTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("role")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unathourised to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("role")
	err = nil
	if userType != role {
		err = errors.New("unathourised to access this resource")
		return err
	}
	return err
}