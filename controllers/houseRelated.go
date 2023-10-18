package controllers

import (
	"energytoken/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Migrate(c *gin.Context) {
	models.Migrate()
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
