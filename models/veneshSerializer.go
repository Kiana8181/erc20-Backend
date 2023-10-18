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

func Serialize(bill BillResponse) (Bill, error) {
	var b Bill
	b.Amount = bill.Data.Amount
	b.BillId = bill.Data.BillId
	b.PayId = bill.Data.PayId
	b.Date = bill.Data.Date
	b.CompanyName = bill.Data.Info.CompanyName
	b.CustomerName = bill.Data.Info.CustomerName
	b.CustomerFamily = bill.Data.Info.CustomerFamily
	b.CustomerType = bill.Data.Info.CustomerType
	b.Address = bill.Data.Info.Address
	b.PostalCode = bill.Data.Info.PostalCode
	b.FileNumber = bill.Data.Info.FileNumber
	b.ComputerPassword = bill.Data.Info.ComputerPassword
	b.IdentificationNumber = bill.Data.Info.IdentificationNumber
	b.TariffType = bill.Data.Info.TariffType
	b.Phase = bill.Data.Info.Phase
	b.Amper = bill.Data.Info.Amper
	b.VoltageType = bill.Data.Info.VoltageType
	b.ContractDemand = bill.Data.Info.ContractDemand
	b.Year = bill.Data.Info.Year
	b.Period = bill.Data.Info.Period
	b.PreviousReadingDate = bill.Data.Info.PreviousReadingDate
	b.CurrentReadingDate = bill.Data.Info.CurrentReadingDate
	b.BillExportationDate = bill.Data.Info.BillExportationDate
	b.RejectDate = bill.Data.Info.RejectDate
	b.NormalConsumption = bill.Data.Info.NormalConsumption
	b.PeakConsumption = bill.Data.Info.PeakConsumption
	b.LowConsumption = bill.Data.Info.LowConsumption
	b.FridayConsumption = bill.Data.Info.FridayConsumption
	b.ReactiveConsumption = bill.Data.Info.ReactiveConsumption
	b.DemandRead = bill.Data.Info.DemandRead
	b.AverageConsumption = bill.Data.Info.AverageConsumption
	b.BillPayableAmount = bill.Data.Info.BillPayableAmount
	b.PeriodAmount = bill.Data.Info.PeriodAmount
	b.InsuranceAmount = bill.Data.Info.InsuranceAmount
	b.TaxAmount = bill.Data.Info.TaxAmount
	b.ElectricityTaxAmount = bill.Data.Info.ElectricityTaxAmount
	b.PreviousDebt = bill.Data.Info.PreviousDebt
	b.EnergyAmount = bill.Data.Info.EnergyAmount
	b.ReactiveAmount = bill.Data.Info.ReactiveAmount
	b.DemandAmount = bill.Data.Info.DemandAmount
	b.SubscriptionAmount = bill.Data.Info.SubscriptionAmount
	b.SeasonAmount = bill.Data.Info.SeasonAmount
	b.FreeAmount = bill.Data.Info.FreeAmount
	b.GasDiscountAmount = bill.Data.Info.GasDiscountAmount
	b.DiscountAmount = bill.Data.Info.DiscountAmount
	b.WarmDaysCount = bill.Data.Info.WarmDaysCount
	b.ColdDaysCount = bill.Data.Info.ColdDaysCount
	b.TotalDaysCount = bill.Data.Info.TotalDaysCount
	b.ConsumptionDebtAmount = bill.Data.Info.ConsumptionDebtAmount
	b.OtherDebtAmount = bill.Data.Info.OtherDebtAmount
	b.BranchDebtAmount = bill.Data.Info.BranchDebtAmount

	// fmt.Printf("%+v\n", b)

	return b, nil
}
