package controllers

import (
	"bytes"
	"encoding/json"
	"energytoken/blockchain"
	"energytoken/models"
	"energytoken/utils/token"
	"net/http"
	"os"
	"strconv"

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

func GetBill(c *gin.Context) {

	//Extract Input
	var input BillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := os.Getenv("VENESH_API") + "/BillInquiry"

	var payload VeneshBill
	payload.Type = "Electricity"
	payload.Parameter = input.Parameter
	payload.Info = "1"

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	veneshToken, err := VeneshGenerateToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	veneshToken = "Bearer " + veneshToken

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", veneshToken)
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	// Extract Bill Data
	billraw := models.BillResponse{}
	// json.Unmarshal([]byte(veneshResponse), &bill)

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&billraw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract Users ID from Token
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bill, err := models.Serialize(billraw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	models.InsertBill(bill, user_id)

	// Calculate Points
	metrage, err := models.GetMetrage(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	conv_metrage, err := strconv.ParseFloat(metrage, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	conv_LowConsumption, err := strconv.ParseFloat(bill.LowConsumption, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	conv_PeakConsumption, err := strconv.ParseFloat(bill.PeakConsumption, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var lowPrice, highPrice, lowCons, highCons, C1, C2, A1, A3 float64
	lowPrice = 3.71
	highPrice = 5.99
	lowCons = 85 - 21*lowPrice + 0.4*conv_metrage
	highCons = 52 - 7*highPrice + 0.3*conv_metrage

	if (conv_LowConsumption > lowCons) || (conv_PeakConsumption > highCons) {
		c.JSON(http.StatusOK, gin.H{"message": "user cannot get any points."})
		return
	}

	C1 = 5.95 - lowPrice
	C2 = 11.71 - highPrice

	A1 = C1 * lowCons * 0.5
	A3 = C2 * highCons * 0.5

	var points int
	points = int(A3 - A1)
	if points <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "user cannot get any points."})
		return
	}

	// Insert point for the user
	username, err := models.GetUsernameByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = blockchain.Mint(username, points)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": points})
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
