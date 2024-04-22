package models

type MergePerson struct {
	SQLEmployeeId        *int64   `json:"SQL_Employee_ID"`
	MongoDBEmployeeID    *string  `json:"mongoDBEmployeeId"`
	FirstName            *string  `json:"firstName"`
	LastName             *string  `json:"lastName"`
	VacationDays         *int64   `json:"vacationDays"`
	PaidToDate           *int64   `json:"paidToDate"`
	PaidLastYear         *int64   `json:"paidLastYear"`
	PayRate              *float64 `json:"payRate"`
	PayRateID            *int64   `json:"payRateId"`
	MiddleInitial        *string  `json:"MiddleInitial"`
	Address1             *string  `json:"Address1"`
	Address2             *string  `json:"Address2"`
	City                 *string  `json:"City"`
	State                *string  `json:"State"`
	Zip                  *int64   `json:"Zip"`
	Email                *string  `json:"Email"`
	PhoneNumber          *string  `json:"PhoneNumber"`
	SocialSecurityNumber *string  `json:"SocialSecurityNumber"`
	DriversLicense       *string  `json:"DriversLicense"`
	MaritalStatus        *string  `json:"MaritalStatus"`
	Gender               *bool    `json:"Gender"`
	ShareholderStatus    *bool    `json:"Shareholder_Status"`
	BenefitPlans         *int64   `json:"BenefitPlans"`
	Ethnicity            *string  `json:"Ethnicity"`
}
