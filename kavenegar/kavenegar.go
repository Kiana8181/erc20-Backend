package kavenegar

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateCode() string {
	low := 100000
	high := 999999
	rand.Seed(time.Now().UnixNano())
	code := strconv.Itoa(low + rand.Intn(high-low))
	return code
}

func GenerateOTP(phoneNumber string) (string, error) {

	url := os.Getenv("KAVENEGAR_API")
	method := "POST"

	OTPCode := GenerateCode()
	payloadString := fmt.Sprintf("receptor=" + phoneNumber + "&token=" + OTPCode + "&template=" + os.Getenv("KAVENEGAR_TEMPLATE"))
	payload := strings.NewReader(payloadString)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", errors.New("unable to create otp message")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", errors.New("unable to send request to otp server")
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("unable to convert response from otp server")
	}
	// fmt.Println(string(body))

	return OTPCode, nil
}
