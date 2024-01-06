package controllers

import (
	"energytoken/blockchain"
	"energytoken/models"
	"energytoken/utils/token"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	current_time := time.Now()
	var trxLog models.Transactions
	trxLog.QQFrom = username
	trxLog.QQTo = input.Receiver
	trxLog.Value = input.Value
	trxLog.Date = fmt.Sprint(current_time.Format("2006-01-02 15:04:05"))
	trxLog.SaveTransaction()

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
