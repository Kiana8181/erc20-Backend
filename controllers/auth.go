package controllers

import (
	"bytes"
	"encoding/json"
	"energytoken/blockchain"
	"energytoken/kavenegar"
	"energytoken/models"
	cacheService "energytoken/redis"
	"energytoken/utils/token"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

type VeneshGenerateTokenReq struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Grant_type    string `json:"grant_type"`
	Client_id     string `json:"client_id"`
	Client_secret string `json:"client_secret"`
}

type VeneshGenerateTokenRes struct {
	Token_type    string `json:"token_type"`
	Expires_in    int    `json:"expires_in"`
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

func VeneshGenerateToken() (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("VENESH_USERNAME")
	password := os.Getenv("VENESH_PASSWORD")
	grant_type := os.Getenv("VENESH_GRANT_TYPE")
	client_id := os.Getenv("VENESH_CLIENT_ID")
	client_secret := os.Getenv("VENESH_CLIENT_SECRET")
	url := os.Getenv("VENESH_OAUTH")

	var payload VeneshGenerateTokenReq
	payload.Username = username
	payload.Password = password
	payload.Grant_type = grant_type
	payload.Client_id = client_id
	payload.Client_secret = client_secret

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var ShahkarInput VeneshGenerateTokenRes
	err = json.Unmarshal(body, &ShahkarInput)
	if err != nil {
		return "", err
	}

	return ShahkarInput.Access_token, nil

}

type ValidatePhoneNumberInput struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	NationalId  string `json:"nationalId" binding:"required"`
}

type ShahkarReq struct {
	NationalId string `json:"NationalId"`
	Mobile     string `json:"Mobile"`
}

type ShahkarRes struct {
	Status  int    `json:"Status"`
	Message string `json:"Message"`
}

func ValidatePhoneAndID(c *gin.Context) {

	var input ValidatePhoneNumberInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := os.Getenv("VENESH_API") + "/Law/Users/Shahkar"

	var payload ShahkarReq
	payload.NationalId = input.NationalId
	payload.Mobile = input.PhoneNumber

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Shahkar Request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := VeneshGenerateToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token = "Bearer " + token

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ShahkarInput ShahkarRes
	err = json.Unmarshal(body, &ShahkarInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Shahkar status on 0 means the phone number and national id matched
	if ShahkarInput.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NationalID or PhoneNumber is not correct."})
		return
	}

	err = cacheService.RegisterCacheInsert(input.PhoneNumber, input.NationalId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

type RegisterInput struct {
	NationalId  string `json:"nationalId" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cacheService.RegisterCacheCheck(input.PhoneNumber, input.NationalId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = blockchain.Register(input.PhoneNumber, "pas123") //password is hardcoded for now
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.NationalId = input.NationalId
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

	//Generate OTP Token
	otpCode, err := kavenegar.GenerateOTP(u.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = cacheService.LoginTokenInsert(u.Username, otpCode, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "an error occurred during authentication process."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

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