package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	EmployeeNotID
}

type EmployeeNotID struct {
	EmployeeId   string    `bson:"employeeId" json:"employeeId"`
	FirstName    string    `bson:"firstName" json:"firstName"`
	LastName     string    `bson:"lastName" json:"lastName"`
	VacationDays int64     `bson:"vacationDays" json:"vacationDays"`
	PaidToDate   int64     `bson:"paidToDate" json:"paidToDate"`
	PaidLastYear int64     `bson:"paidLastYear" json:"paidLastYear"`
	PayRate      float64   `bson:"payRate" json:"payRate"`
	PayRateID    int64     `bson:"payRateId" json:"payRateId"`
	CreatedAt    time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
