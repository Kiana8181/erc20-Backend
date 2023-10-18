package controllers

import (
	"energytoken/models"
	"energytoken/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Migrate(c *gin.Context) {
	models.Migrate()
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

type ResidentialInput struct {
	PostalCode  string `json:"postalCode"`
	Ownership   string `json:"ownership"`
	MeterNumber string `json:"meter_number"`
}

func ResidentialProposal(c *gin.Context) {
	nationalID, err := token.ExtractTokenNationalID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input ResidentialInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var p = models.ResidentialProposal{}
	p.NationalID = nationalID
	p.PostalCode = input.PostalCode
	p.MeterNumber = input.MeterNumber
	p.Ownership = input.Ownership
	p.PhoneNumber = u.Username
	p.Meterage = u.Meterage
	p.Status = 0

	_, err = p.SaveResidentialProposal()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// func UpdateResidential(c *gin.Context) {

// 	user_id, err := token.ExtractTokenID(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	username, err := models.GetUsernameByID(user_id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}

// 	var input ResidentialInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	models.SetResidential(username, input.PostalCode, input.Ownership, input.MeterNumber)

// 	c.JSON(http.StatusOK, gin.H{"message": "OK"})

// }

type BillInput struct {
	Parameter string `json:"Parameter"`
}

type VeneshBill struct {
	Type      string `json:"Type"`
	Parameter string `json:"Parameter"`
	Info      string `json:"Info"`
}

func GetAllResidentialProposal(c *gin.Context) {
	proposals := models.GetAllResidentialProposal()
	c.JSON(http.StatusOK, gin.H{"data": proposals})
}

type AcceptPorposalInput struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	PostalCode  string `json:"postalCode"`
	Ownership   string `json:"ownership"`
	MeterNumber string `json:"meter_number"`
	Meterage    string `json:"meterage"`
}

func AcceptProposal(c *gin.Context) {
	var input AcceptPorposalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := input.PhoneNumber

	err := models.SetResidential(username, input.PostalCode, input.Ownership, input.MeterNumber, input.Meterage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	err = models.UpdateProposal(input.ID, input.PostalCode, input.Ownership, input.MeterNumber, input.Meterage, 1) //accept code for a proposal is 1
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

type RejectPorposalInput struct {
	ID uint `json:"id"`
}

func RejectProposal(c *gin.Context) {
	var input RejectPorposalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.RejectProposal(input.ID) //accept code for a proposal is 2
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
