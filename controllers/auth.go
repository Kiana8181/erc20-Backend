package controllers

import (
	"energytoken/blockchain"
	"energytoken/models"
	cacheService "energytoken/redis"
	"energytoken/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

type RegisterInput struct {
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := blockchain.Register(input.PhoneNumber, "pas123") //password is hardcoded for now ;)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Email = input.Email
	u.Username = input.PhoneNumber
	u.Password = input.Password
	u.IsSuper = false

	accountID, err := blockchain.AccountID(input.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.WalletID = accountID

	_, err = u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})

}

type LoginInput struct {
	Username string `json:"phoneNumber" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	// Generate Token
	token, err := models.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	// err = cacheService.LoginTokenInsert(u.Username, token)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "an error occurred during authentication process."})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"token": token})

}

type VerifyInput struct {
	Username string `json:"phoneNumber" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func VerifyLogin(c *gin.Context) {
	var input VerifyInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := cacheService.LoginTokenFetch(input.Username, input.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
