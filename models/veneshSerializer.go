package models

type BillResponse struct {
	UserID uint     `json:"userID"`
	Data   BillData `json:"data"`
}

type BillData struct {
	Amount string `json:"Amount"`
	BillId string `json:"BillId"`
	PayId  string `json:"PayId"`
	Date   string `json:"Date"`
	Info   Info   `json:"info"`
}

type Info struct {
	CompanyName           string `json:"CompanyName"`
	CustomerName          string `json:"CustomerName"`
	CustomerFamily        string `json:"CustomerFamily"`
	CustomerType          string `json:"CustomerType"`
	Address               string `json:"Address"`
	PostalCode            string `json:"PostalCode"`
	FileNumber            string `json:"FileNumber"`
	ComputerPassword      string `json:"ComputerPassword"`
	IdentificationNumber  string `json:"IdentificationNumber"`
	TariffType            string `json:"TariffType"`
	Phase                 string `json:"Phase"`
	Amper                 string `json:"Amper"`
	VoltageType           string `json:"VoltageType"`
	ContractDemand        string `json:"ContractDemand"`
	Year                  string `json:"Year"`
	Period                string `json:"Period"`
	PreviousReadingDate   string `json:"PreviousReadingDate"`
	CurrentReadingDate    string `json:"CurrentReadingDate"`
	BillExportationDate   string `json:"BillExportationDate"`
	RejectDate            string `json:"RejectDate"`
	NormalConsumption     string `json:"NormalConsumption"`
	PeakConsumption       string `json:"PeakConsumption"`
	LowConsumption        string `json:"LowConsumption"`
	FridayConsumption     string `json:"FridayConsumption"`
	ReactiveConsumption   string `json:"ReactiveConsumption"`
	DemandRead            string `json:"DemandRead"`
	AverageConsumption    string `json:"AverageConsumption"`
	BillPayableAmount     string `json:"BillPayableAmount"`
	PeriodAmount          string `json:"PeriodAmount"`
	InsuranceAmount       string `json:"InsuranceAmount"`
	TaxAmount             string `json:"PaytollAmount"`
	ElectricityTaxAmount  string `json:"ElectricityTaxAmount"`
	PreviousDebt          string `json:"PreviousDebt"`
	EnergyAmount          string `json:"EnergyAmount"`
	ReactiveAmount        string `json:"ReactiveAmount"`
	DemandAmount          string `json:"DemandAmount"`
	SubscriptionAmount    string `json:"SubscriptionAmount"`
	SeasonAmount          string `json:"SeasonAmount"`
	FreeAmount            string `json:"FreeAmount"`
	GasDiscountAmount     string `json:"GasDiscountAmount"`
	DiscountAmount        string `json:"DiscountAmount"`
	WarmDaysCount         string `json:"WarmDaysCount"`
	ColdDaysCount         string `json:"ColdDaysCount"`
	TotalDaysCount        string `json:"TotalDaysCount"`
	ConsumptionDebtAmount string `json:"ConsumptionDebtAmount"`
	OtherDebtAmount       string `json:"OtherDebtAmount"`
	BranchDebtAmount      string `json:"BranchDebtAmount"`
}
