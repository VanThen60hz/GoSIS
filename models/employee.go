package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	MongoDBObjectId primitive.ObjectID `json:"id,omitempty"`
	EmployeeId      string             `json:"mongoDBEmployeeId"`
	FirstName       string             `json:"firstName"`
	LastName        string             `json:"lastName"`
	VacationDays    int64              `json:"vacationDays"`
	PaidToDate      int64              `json:"paidToDate"`
	PaidLastYear    int64              `json:"paidLastYear"`
	PayRate         float64            `json:"payRate"`
	PayRateID       int64              `json:"payRateId"`
}
