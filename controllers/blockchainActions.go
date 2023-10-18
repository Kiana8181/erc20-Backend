package controllers

import (
	"energytoken/blockchain"
	"energytoken/models"
	"energytoken/utils/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransferReq struct {
	Value    string `json:"value"`
	Receiver string `json:"receiver"`
}

func Transfer(c *gin.Context) {
	var input TransferReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := models.GetUsernameByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	convertedValue, err := strconv.Atoi(input.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = blockchain.Transfer(username, input.Receiver, convertedValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func Balance(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := models.GetUsernameByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	value, err := blockchain.Balance(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}
