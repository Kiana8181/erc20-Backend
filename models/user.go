package models

import (
	"errors"

	// "fmt"
	"html"
	"strings"

	"energytoken/utils/token"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:225;not null;unique" json:"username"`
	Password  string `gorm:"size:225;not null" json:"password"`
	WalletID  string `gorm:"size:255" json:"walletId"`
	Email     string `gorm:"size:225" json:"email"`
	Firstname string `gorm:"size:225" json:"firstname"`
	Lastname  string `gorm:"size:255" json:"lastname"`
	IsActive  bool   `gorm:"size:255" json:"is_active"`
	IsSuper   bool   `gorm:"size:255" json:"is_super"`
	// History   []Transactions `gorm:"ForeignKey:ID" json:"mytransactions"`
}

type Transactions struct {
	gorm.Model
	// ID    uint
	// Key   uint   `gorm:"primary_key"`
	QQFrom string `gorm:"size:255" json:"qqfrom"`
	QQTo   string `gorm:"size:255" json:"qqto"`
	Date   string `gorm:"size:255" json:"date"`
	Value  string `gorm:"size:255" json:"value"`
}

type ResidentialProposal struct {
	gorm.Model
	NationalID  string `gorm:"size:225;not null" json:"nationalId"`
	PostalCode  string `gorm:"size:225" json:"postalCode"`
	MeterNumber string `gorm:"size:255" json:"meter_number"`
	Ownership   string `gorm:"size:255" json:"ownership"`
	PhoneNumber string `gorm:"size:255" json:"phoneNumber"`
	Meterage    string `gorm:"size:255" json:"meterage"`
	Status      uint   `gorm:"size:255" json:"status"`
}

func Migrate() {
	DB.AutoMigrate(&User{}, &Transactions{}, &ResidentialProposal{})
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUserByID(uid uint) (User, error) {
	var u User

	DB.First(&u, uid)
	// DB.Preload("History").First(&u, uid)
	u.PrepareGive()

	return u, nil
}

func GetUsernameByID(uid uint) (string, error) {
	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return "", errors.New("user not found")
	}

	return u.Username, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type GetTransactionsResponse struct {
	From  string `gorm:"size:255" json:"from"`
	To    string `gorm:"size:255" json:"to"`
	Date  string `gorm:"size:255" json:"date"`
	Value string `gorm:"size:255" json:"value"`
	Type  string `gorm:"size:255" json:"type"`
}

func GetTransactions(username string, accountID string) []GetTransactionsResponse {
	t1 := []Transactions{}
	t2 := []Transactions{}
	t3 := []Transactions{}

	DB.Find(&t1, "qq_to = ?", username)
	DB.Find(&t2, "qq_from = ?", username)
	DB.Find(&t3, "qq_to = ?", accountID)

	var JsonForm []GetTransactionsResponse
	var JsonElement GetTransactionsResponse

	t := append(t1, t2...)
	t = append(t, t3...)

	for i := 0; i < len(t); i++ {
		JsonElement.Date = t[i].Date
		JsonElement.From = t[i].QQFrom
		JsonElement.To = t[i].QQTo
		JsonElement.Value = t[i].Value

		if t[i].QQFrom == "System" {
			JsonElement.Type = "0"
		}
		if t[i].QQFrom == username {
			JsonElement.Type = "1"
		}
		if t[i].QQTo == accountID {
			JsonElement.Type = "2"
		}

		JsonForm = append(JsonForm, JsonElement)
	}

	return JsonForm
}

func LoginCheck(username string, password string) (string, error) {

	var err error
	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID, u.Email, u.WalletID, u.IsSuper)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *User) SaveUser() (*User, error) {
	var err error = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (t *Transactions) SaveTransaction() (*Transactions, error) {
	// b.ID = id
	var err error = DB.Create(&t).Error
	if err != nil {
		return &Transactions{}, err
	}
	return t, nil
}

func (u *User) BeforeSave() error {

	// Turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	// remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

// func InsertTransaction(transaction Transactions, user_id uint) {
// 	// fmt.Printf("%+v\n", bill)
// 	transaction.SaveTransaction(user_id)
// }
