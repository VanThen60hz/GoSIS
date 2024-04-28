package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"GoSIS/models"
	"GoSIS/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pusher/pusher-http-go/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

var pusherClient = pusher.Client{
	AppID:   "1790030",
	Key:     "a359a59a30b4ddb07bb5",
	Secret:  "7010eecf86f9469246bf",
	Cluster: "ap1",
	Secure:  true,
}

func CreateBoth(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mergePersonWithoutId models.MergePersonWithoutId
	if err := c.BodyParser(&mergePersonWithoutId); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CreateBothResponse{Status: http.StatusBadRequest, Message: "invalid merge person data", Data: nil})
	}

	e := models.EmployeeNotID{
		FirstName:    *mergePersonWithoutId.FirstName,
		LastName:     *mergePersonWithoutId.LastName,
		VacationDays: *mergePersonWithoutId.VacationDays,
		PaidToDate:   *mergePersonWithoutId.PaidToDate,
		PaidLastYear: *mergePersonWithoutId.PaidLastYear,
		PayRate:      *mergePersonWithoutId.PayRate,
		PayRateID:    *mergePersonWithoutId.PayRateID,
	}

	e.EmployeeId = uuid.New().String()

	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	_, err := employeeCollection.InsertOne(ctx, e)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(http.StatusConflict).JSON(responses.EmployeeResponse{Status: http.StatusConflict, Message: err.Error(), Data: nil})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error creating employee", Data: nil})
	}

	p := models.Personal{
		FirstName:            *mergePersonWithoutId.FirstName,
		LastName:             *mergePersonWithoutId.LastName,
		MiddleInitial:        *mergePersonWithoutId.MiddleInitial,
		Address1:             *mergePersonWithoutId.Address1,
		Address2:             *mergePersonWithoutId.Address2,
		City:                 *mergePersonWithoutId.City,
		State:                *mergePersonWithoutId.State,
		Zip:                  *mergePersonWithoutId.Zip,
		Email:                *mergePersonWithoutId.Email,
		PhoneNumber:          *mergePersonWithoutId.PhoneNumber,
		SocialSecurityNumber: *mergePersonWithoutId.SocialSecurityNumber,
		DriversLicense:       *mergePersonWithoutId.DriversLicense,
		MaritalStatus:        *mergePersonWithoutId.MaritalStatus,
		Gender:               *mergePersonWithoutId.Gender,
		ShareholderStatus:    *mergePersonWithoutId.ShareholderStatus,
		BenefitPlans:         *mergePersonWithoutId.BenefitPlans,
		Ethnicity:            *mergePersonWithoutId.Ethnicity,
	}

	var lastInsertID int64
	err = sqlServerDB.QueryRowContext(ctx, "SELECT TOP 1 Employee_ID FROM Personal ORDER BY Employee_ID DESC").Scan(&lastInsertID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil})
	}
	newEmployeeID := lastInsertID + 1
	p.SQLEmployeeId = newEmployeeID

	_, err = sqlServerDB.ExecContext(ctx, "INSERT INTO Personal (Employee_ID, First_Name, Last_Name, Middle_Initial, Address1, Address2, City, State, Zip, Email, Phone_Number, Social_Security_Number, Drivers_License, Marital_Status, Gender, Shareholder_Status, Benefit_Plans, Ethnicity) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18)",
		sql.Named("p1", newEmployeeID), sql.Named("p2", p.FirstName), sql.Named("p3", p.LastName), sql.Named("p4", p.MiddleInitial), sql.Named("p5", p.Address1), sql.Named("p6", p.Address2), sql.Named("p7", p.City), sql.Named("p8", p.State), sql.Named("p9", p.Zip), sql.Named("p10", p.Email), sql.Named("p11", p.PhoneNumber), sql.Named("p12", p.SocialSecurityNumber), sql.Named("p13", p.DriversLicense), sql.Named("p14", p.MaritalStatus), sql.Named("p15", p.Gender), sql.Named("p16", p.ShareholderStatus), sql.Named("p17", p.BenefitPlans), sql.Named("p18", p.Ethnicity))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil})
	}

	createBothData := map[string]interface{}{
		"mongoDBemployee":   e,
		"sqlServerPersonal": p,
	}

	err = pusherClient.Trigger("GoSIS", "both-created", createBothData)
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.Status(http.StatusCreated).JSON(responses.CreateBothResponse{Status: http.StatusCreated, Message: "employee created successfully", Data: &mergePersonWithoutId})
}
