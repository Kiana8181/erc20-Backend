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
	Username  string         `gorm:"size:225;not null;unique" json:"username"`
	Password  string         `gorm:"size:225;not null" json:"password"`
	WalletID  string         `gorm:"size:255" json:"walletId"`
	Email     string         `gorm:"size:225" json:"email"`
	Firstname string         `gorm:"size:225" json:"firstname"`
	Lastname  string         `gorm:"size:255" json:"lastname"`
	IsActive  bool           `gorm:"size:255" json:"is_active"`
	IsSuper   bool           `gorm:"size:255" json:"is_super"`
	History   []Transactions `gorm:"ForeignKey:ID" json:"mytransactions"`
}

type Transactions struct {
	// gorm.Model
	ID    uint
	Key   uint   `gorm:"primary_key"`
	From  string `gorm:"size:255" json:"from"`
	To    string `gorm:"size:255" json:"to"`
	Date  string `gorm:"size:255" json:"date"`
	Value string `gorm:"size:255" json:"value"`
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

	DB.Preload("History").First(&u, uid)
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

func (b *Transactions) SaveTransaction(id uint) (*Transactions, error) {
	b.ID = id
	var err error = DB.Create(&b).Error
	if err != nil {
		return &Transactions{}, err
	}
	return b, nil
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

func InsertTransaction(transaction Transactions, user_id uint) {
	// fmt.Printf("%+v\n", bill)
	transaction.SaveTransaction(user_id)
}
