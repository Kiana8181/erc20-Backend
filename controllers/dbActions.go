package controllers

import (
	"energytoken/blockchain"
	"energytoken/models"
	"energytoken/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Migrate(c *gin.Context) {
	models.Migrate()
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func GetTransactions(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := models.GetUsernameByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	accountID, _ := blockchain.AccountID(username)
	Transactions := models.GetTransactions(username, accountID)
	// fmt.Println(Transactions)
	// fmt.Println(trxs)
	c.JSON(http.StatusOK, gin.H{"data": Transactions})

}
