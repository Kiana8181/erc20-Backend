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
	Username    string `gorm:"size:225;not null;unique" json:"username"`
	Password    string `gorm:"size:225;not null" json:"password"`
	WalletID    string `gorm:"size:255" json:"walletId"`
	NationalId  string `gorm:"size:225" json:"nationalId"`
	PostalCode  string `gorm:"size:225" json:"postalCode"`
	MeterNumber string `gorm:"size:255" json:"meter_number"`
	Ownership   string `gorm:"size:255" json:"ownership"`
	Meterage    string `gorm:"size:255" json:"meterage"`
	IsActive    bool   `gorm:"size:255" json:"is_active"`
	IsSuper     bool   `gorm:"size:255" json:"is_super"`
	BillList    []Bill `gorm:"ForeignKey:ID" json:"mybills"`
}

type Bill struct {
	// gorm.Model
	ID                    uint
	Key                   uint   `gorm:"primary_key"`
	Amount                string `gorm:"size:255" json:"Amount"`
	BillId                string `gorm:"size:255" json:"BillId"`
	PayId                 string `gorm:"size:255" json:"PayId"`
	Date                  string `gorm:"size:255" json:"Date"`
	CompanyName           string `gorm:"size:255" json:"CompanyName"`
	CustomerName          string `gorm:"size:255" json:"CustomerName"`
	CustomerFamily        string `gorm:"size:255" json:"CustomerFamily"`
	CustomerType          string `gorm:"size:255" json:"CustomerType"`
	Address               string `gorm:"size:255" json:"Address"`
	PostalCode            string `gorm:"size:255" json:"PostalCode"`
	FileNumber            string `gorm:"size:255" json:"FileNumber"`
	ComputerPassword      string `gorm:"size:255" json:"ComputerPassword"`
	IdentificationNumber  string `gorm:"size:255" json:"IdentificationNumber"`
	TariffType            string `gorm:"size:255" json:"TariffType"`
	Phase                 string `gorm:"size:255" json:"Phase"`
	Amper                 string `gorm:"size:255" json:"Amper"`
	VoltageType           string `gorm:"size:255" json:"VoltageType"`
	ContractDemand        string `gorm:"size:255" json:"ContractDemand"`
	Year                  string `gorm:"size:255" json:"Year"`
	Period                string `gorm:"size:255" json:"Period"`
	PreviousReadingDate   string `gorm:"size:255" json:"PreviousReadingDate"`
	CurrentReadingDate    string `gorm:"size:255" json:"CurrentReadingDate"`
	BillExportationDate   string `gorm:"size:255" json:"BillExportationDate"`
	RejectDate            string `gorm:"size:255" json:"RejectDate"`
	NormalConsumption     string `gorm:"size:255" json:"NormalConsumption"`
	PeakConsumption       string `gorm:"size:255" json:"PeakConsumption"`
	LowConsumption        string `gorm:"size:255" json:"LowConsumption"`
	FridayConsumption     string `gorm:"size:255" json:"FridayConsumption"`
	ReactiveConsumption   string `gorm:"size:255" json:"ReactiveConsumption"`
	DemandRead            string `gorm:"size:255" json:"DemandRead"`
	AverageConsumption    string `gorm:"size:255" json:"AverageConsumption"`
	BillPayableAmount     string `gorm:"size:255" json:"BillPayableAmount"`
	PeriodAmount          string `gorm:"size:255" json:"PeriodAmount"`
	InsuranceAmount       string `gorm:"size:255" json:"InsuranceAmount"`
	TaxAmount             string `gorm:"size:255" json:"PaytollAmount"`
	ElectricityTaxAmount  string `gorm:"size:255" json:"ElectricityTaxAmount"`
	PreviousDebt          string `gorm:"size:255" json:"PreviousDebt"`
	EnergyAmount          string `gorm:"size:255" json:"EnergyAmount"`
	ReactiveAmount        string `gorm:"size:255" json:"ReactiveAmount"`
	DemandAmount          string `gorm:"size:255" json:"DemandAmount"`
	SubscriptionAmount    string `gorm:"size:255" json:"SubscriptionAmount"`
	SeasonAmount          string `gorm:"size:255" json:"SeasonAmount"`
	FreeAmount            string `gorm:"size:255" json:"FreeAmount"`
	GasDiscountAmount     string `gorm:"size:255" json:"GasDiscountAmount"`
	DiscountAmount        string `gorm:"size:255" json:"DiscountAmount"`
	WarmDaysCount         string `gorm:"size:255" json:"WarmDaysCount"`
	ColdDaysCount         string `gorm:"size:255" json:"ColdDaysCount"`
	TotalDaysCount        string `gorm:"size:255" json:"TotalDaysCount"`
	ConsumptionDebtAmount string `gorm:"size:255" json:"ConsumptionDebtAmount"`
	OtherDebtAmount       string `gorm:"size:255" json:"OtherDebtAmount"`
	BranchDebtAmount      string `gorm:"size:255" json:"BranchDebtAmount"`
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
	DB.AutoMigrate(&User{}, &Bill{}, &ResidentialProposal{})
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUserByID(uid uint) (User, error) {
	var u User

	// Retreive Users data
	DB.Preload("BillList").First(&u, uid)
	// if err := DB.First(&u, uid).Error; err != nil {
	// 	return u, errors.New("user not found")
	// }

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

	token, err := token.GenerateToken(u.ID, u.NationalId, u.WalletID, u.IsSuper)
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

func (b *Bill) SaveBill(id uint) (*Bill, error) {
	b.ID = id
	var err error = DB.Create(&b).Error
	if err != nil {
		return &Bill{}, err
	}
	return b, nil
}

func (p *ResidentialProposal) SaveResidentialProposal() (*ResidentialProposal, error) {
	var err error = DB.Create(&p).Error
	if err != nil {
		return &ResidentialProposal{}, err
	}
	return p, nil
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

func SetResidential(username string, postalCode string, ownership string, meterNumber string, meterage string) error {
	var u User
	err := DB.Model(&u).Where("username = ?", username).Updates(map[string]interface{}{"postalCode": postalCode, "ownership": ownership, "meter_number": meterNumber, "meterage": meterage}).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateProposal(ID uint, postalCode string, ownership string, meterNumber string, meterage string, status uint) error {
	var p ResidentialProposal
	err := DB.Model(&p).Where("id = ?", ID).Update(map[string]interface{}{"postalCode": postalCode, "ownership": ownership, "meter_number": meterNumber, "meterage": meterage, "status": status}).Error
	if err != nil {
		return err
	}
	return nil
}

func RejectProposal(ID uint) error {
	var p ResidentialProposal
	err := DB.Model(&p).Where("id = ?", ID).Update(map[string]interface{}{"status": 2}).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertBill(bill Bill, user_id uint) {
	// fmt.Printf("%+v\n", bill)
	bill.SaveBill(user_id)
}

func GetAllResidentialProposal() []ResidentialProposal {
	var p []ResidentialProposal
	DB.Find(&p)
	return p
}

func GetMetrage(user_id uint) (string, error) {
	var u User
	if err := DB.First(&u, user_id).Error; err != nil {
		return "", errors.New("user not found")
	}
	return u.Meterage, nil
}
