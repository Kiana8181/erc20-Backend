package blockchain

//API to hyperledger fabric API

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(username string, password string) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	url := os.Getenv("HYPERLEDGER_API") + "/register"
	token := os.Getenv("HYPERLEDGER_TOKEN")

	var payload RegisterReq
	payload.Username = username
	payload.Password = password

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("unable to create wallet")
	}

	return nil
}

type MintReq struct {
	Username string `json:"username"`
	Value    int    `json:"value"`
}

func Mint(username string, value int) error {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	url := os.Getenv("HYPERLEDGER_API") + "/mint"
	token := os.Getenv("HYPERLEDGER_TOKEN")

	var payload MintReq
	payload.Username = username
	payload.Value = value

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("unable to mint")
	}

	return nil
}

type BalanceInput struct {
	Username string `json:"username"`
}

type BalanceReq struct {
	Username string `json:"username"`
}

type BalanceRes struct {
	Message string `json:"message"`
	Value   string `json:"value"`
}

func Balance(username string) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	url := os.Getenv("HYPERLEDGER_API") + "/balance"
	token := os.Getenv("HYPERLEDGER_TOKEN")

	var payload BalanceReq
	payload.Username = username

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
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	if err != nil {
		return "0", nil // error will be occuared when balance is 0
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", errors.New("cannot get a response from the network")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result BalanceRes
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Value, nil
}

type AccountIDReq struct {
	Username string `json:"username"`
}

type AccountIDRes struct {
	ID string `json:"id"`
}

func AccountID(username string) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	url := os.Getenv("HYPERLEDGER_API") + "/accountid"
	token := os.Getenv("HYPERLEDGER_TOKEN")

	var payload AccountIDReq
	payload.Username = username

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
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", errors.New("unable to get wallet id")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result AccountIDRes
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.ID, nil
}

type TransferReq struct {
	Username string `json:"username"`
	Receiver string `json:"receiver"`
	Value    int    `json:"value"`
}

func Transfer(username string, receiver string, value int) error {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	url := os.Getenv("HYPERLEDGER_API") + "/transfer"
	token := os.Getenv("HYPERLEDGER_TOKEN")

	var payload TransferReq
	payload.Username = username
	payload.Receiver = receiver
	payload.Value = value

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("unable to transfer")
	}

	return nil
}
