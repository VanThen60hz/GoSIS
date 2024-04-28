package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	EmployeeId   string             `json:"employeeId"`
	FirstName    string             `json:"firstName"`
	LastName     string             `json:"lastName"`
	VacationDays int64              `json:"vacationDays"`
	PaidToDate   int64              `json:"paidToDate"`
	PaidLastYear int64              `json:"paidLastYear"`
	PayRate      float64            `json:"payRate"`
	PayRateID    int64              `json:"payRateId"`
	CreatedAt    time.Time          `json:"createdAt,omitempty"`
	UpdatedAt    time.Time          `json:"updatedAt,omitempty"`
}

type EmployeeNotID struct {
	EmployeeId   string    `bson:"employeeId" form:"employeeId"`
	FirstName    string    `bson:"firstName" form:"firstName"`
	LastName     string    `bson:"lastName" form:"lastName"`
	VacationDays int64     `bson:"vacationDays" form:"vacationDays"`
	PaidToDate   int64     `bson:"paidToDate" form:"paidToDate"`
	PaidLastYear int64     `bson:"paidLastYear" form:"paidLastYear"`
	PayRate      float64   `bson:"payRate" form:"payRate"`
	PayRateID    int64     `bson:"payRateId" form:"payRateId"`
	CreatedAt    time.Time `bson:"createdAt,omitempty"`
	UpdatedAt    time.Time `bson:"updatedAt,omitempty"`
}
