package models

type Personal struct {
	SQLEmployeeId        int64  `json:"SQL_Employee_ID"`
	FirstName            string `json:"First_Name"`
	LastName             string `json:"Last_Name"`
	MiddleInitial        string `json:"Middle_Initial"`
	Address1             string `json:"Address1"`
	Address2             string `json:"Address2"`
	City                 string `json:"City"`
	State                string `json:"State"`
	Zip                  int64  `json:"Zip"`
	Email                string `json:"Email"`
	PhoneNumber          string `json:"PhoneNumber"`
	SocialSecurityNumber string `json:"SocialSecurityNumber"`
	DriversLicense       string `json:"DriversLicense"`
	MaritalStatus        string `json:"Marital_Status"`
	Gender               bool   `json:"Gender"`
	ShareholderStatus    bool   `json:"Shareholder_Status"`
	BenefitPlans         int64  `json:"Benefit_Plans"`
	Ethnicity            string `json:"Ethnicity"`
}
