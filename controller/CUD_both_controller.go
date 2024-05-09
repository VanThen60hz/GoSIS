package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	configs "GoSIS/config"
	"GoSIS/models"
	"GoSIS/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pusher/pusher-http-go/v5"
	"go.mongodb.org/mongo-driver/bson"
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

	session, err := configs.MongoDB.StartSession()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EmployeeResponse{
			Status: http.StatusInternalServerError,

			Message: "error starting MongoDB session",

			Data: nil,
		})
	}

	tx, err := sqlServerDB.BeginTx(ctx, nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{
			Status: http.StatusInternalServerError,

			Message: "error starting SQL Server transaction",

			Data: nil,
		})
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	defer session.EndSession(ctx)

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

	_, err = employeeCollection.InsertOne(ctx, e)
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

func UpdateBoth(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mergePerson models.MergePerson
	if err := c.BodyParser(&mergePerson); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CreateBothResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid merge person data",
			Data:    nil,
		})
	}

	if mergePerson.MongoDBEmployeeID == nil {
		fmt.Println("tạo mới mongo db")
		e := models.EmployeeNotID{
			FirstName:    *mergePerson.FirstName,
			LastName:     *mergePerson.LastName,
			VacationDays: *mergePerson.VacationDays,
			PaidToDate:   *mergePerson.PaidToDate,
			PaidLastYear: *mergePerson.PaidLastYear,
			PayRate:      *mergePerson.PayRate,
			PayRateID:    *mergePerson.PayRateID,
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

		fmt.Println(e.EmployeeId)
		mergePerson.MongoDBEmployeeID = &e.EmployeeId

	} else {
		// Update MongoDB
		filter := bson.M{"employeeId": mergePerson.MongoDBEmployeeID}
		update := bson.M{
			"$set": bson.M{
				"firstName":    mergePerson.FirstName,
				"lastName":     mergePerson.LastName,
				"vacationDays": mergePerson.VacationDays,
				"paidToDate":   mergePerson.PaidToDate,
				"paidLastYear": mergePerson.PaidLastYear,
				"payRate":      mergePerson.PayRate,
				"payRateId":    mergePerson.PayRateID,
				"updatedAt":    time.Now(),
			},
		}
		_, err := employeeCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UpdateEmployeeResponse{
				Status:  http.StatusInternalServerError,
				Message: "error updating employee in MongoDB",
				Data:    nil,
			})
		}
	}

	if mergePerson.SQLEmployeeId == nil {
		fmt.Println("tạo ms sql")
		p := models.Personal{
			FirstName:            *mergePerson.FirstName,
			LastName:             *mergePerson.LastName,
			MiddleInitial:        *mergePerson.MiddleInitial,
			Address1:             *mergePerson.Address1,
			Address2:             *mergePerson.Address2,
			City:                 *mergePerson.City,
			State:                *mergePerson.State,
			Zip:                  *mergePerson.Zip,
			Email:                *mergePerson.Email,
			PhoneNumber:          *mergePerson.PhoneNumber,
			SocialSecurityNumber: *mergePerson.SocialSecurityNumber,
			DriversLicense:       *mergePerson.DriversLicense,
			MaritalStatus:        *mergePerson.MaritalStatus,
			Gender:               *mergePerson.Gender,
			ShareholderStatus:    *mergePerson.ShareholderStatus,
			BenefitPlans:         *mergePerson.BenefitPlans,
			Ethnicity:            *mergePerson.Ethnicity,
		}

		var lastInsertID int64
		err := sqlServerDB.QueryRowContext(ctx, "SELECT TOP 1 Employee_ID FROM Personal ORDER BY Employee_ID DESC").Scan(&lastInsertID)
		if err != nil {
			// Handle SQL query error
			return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}
		newEmployeeID := lastInsertID + 1
		p.SQLEmployeeId = newEmployeeID

		// Insert new record into SQL database
		_, err = sqlServerDB.ExecContext(ctx, "INSERT INTO Personal (Employee_ID, First_Name, Last_Name, Middle_Initial, Address1, Address2, City, State, Zip, Email, Phone_Number, Social_Security_Number, Drivers_License, Marital_Status, Gender, Shareholder_Status, Benefit_Plans, Ethnicity) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18)",
			sql.Named("p1", newEmployeeID), sql.Named("p2", p.FirstName), sql.Named("p3", p.LastName), sql.Named("p4", p.MiddleInitial), sql.Named("p5", p.Address1), sql.Named("p6", p.Address2), sql.Named("p7", p.City), sql.Named("p8", p.State), sql.Named("p9", p.Zip), sql.Named("p10", p.Email), sql.Named("p11", p.PhoneNumber), sql.Named("p12", p.SocialSecurityNumber), sql.Named("p13", p.DriversLicense), sql.Named("p14", p.MaritalStatus), sql.Named("p15", p.Gender), sql.Named("p16", p.ShareholderStatus), sql.Named("p17", p.BenefitPlans), sql.Named("p18", p.Ethnicity))
		if err != nil {
			// Handle SQL insert error
			return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		// Update mergePerson.SQLEmployeeId with the newEmployeeID
		mergePerson.SQLEmployeeId = &newEmployeeID

	} else {
		p := &models.Personal{
			SQLEmployeeId:        *mergePerson.SQLEmployeeId,
			FirstName:            *mergePerson.FirstName,
			LastName:             *mergePerson.LastName,
			MiddleInitial:        *mergePerson.MiddleInitial,
			Address1:             *mergePerson.Address1,
			Address2:             *mergePerson.Address2,
			City:                 *mergePerson.City,
			State:                *mergePerson.State,
			Zip:                  *mergePerson.Zip,
			Email:                *mergePerson.Email,
			PhoneNumber:          *mergePerson.PhoneNumber,
			SocialSecurityNumber: *mergePerson.SocialSecurityNumber,
			DriversLicense:       *mergePerson.DriversLicense,
			MaritalStatus:        *mergePerson.MaritalStatus,
			Gender:               *mergePerson.Gender,
			ShareholderStatus:    *mergePerson.ShareholderStatus,
			BenefitPlans:         *mergePerson.BenefitPlans,
			Ethnicity:            *mergePerson.Ethnicity,
		}

		_, err := sqlServerDB.ExecContext(ctx, `
	UPDATE Personal 
	SET First_Name = @p2, Last_Name = @p3, Middle_Initial = @p4, Address1 = @p5, Address2 = @p6, 
		City = @p7, State = @p8, Zip = @p9, Email = @p10, Phone_Number = @p11, 
		Social_Security_Number = @p12, Drivers_License = @p13, Marital_Status = @p14, 
		Gender = @p15, Shareholder_Status = @p16, Benefit_Plans = @p17, Ethnicity = @p18
	WHERE Employee_ID = @p1`,
			sql.Named("p1", p.SQLEmployeeId),
			sql.Named("p2", p.FirstName),
			sql.Named("p3", p.LastName),
			sql.Named("p4", p.MiddleInitial),
			sql.Named("p5", p.Address1),
			sql.Named("p6", p.Address2),
			sql.Named("p7", p.City),
			sql.Named("p8", p.State),
			sql.Named("p9", p.Zip),
			sql.Named("p10", p.Email),
			sql.Named("p11", p.PhoneNumber),
			sql.Named("p12", p.SocialSecurityNumber),
			sql.Named("p13", p.DriversLicense),
			sql.Named("p14", p.MaritalStatus),
			sql.Named("p15", p.Gender),
			sql.Named("p16", p.ShareholderStatus),
			sql.Named("p17", p.BenefitPlans),
			sql.Named("p18", p.Ethnicity),
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

	}

	err := pusherClient.Trigger("GoSIS", "both-edited", mergePerson)
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.Status(http.StatusCreated).JSON(responses.UpdateBothResponse{
		Status:  http.StatusOK,
		Message: "edit both successfully",
		Data:    &mergePerson,
	})
}

func DeleteBoth(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type DeleteRequest struct {
		SQLEmployeeID     *int64  `json:"SQL_Employee_ID"`
		MongoDBEmployeeID *string `json:"mongoDBEmployeeId"`
	}

	var request DeleteRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
		})
	}

	if request.SQLEmployeeID == nil && request.MongoDBEmployeeID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing delete info in request body",
		})
	}

	session, err := configs.MongoDB.StartSession()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to start MongoDB session",
		})
	}
	defer session.EndSession(ctx)

	tx, err := sqlServerDB.BeginTx(ctx, nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to start SQL Server transaction",
		})
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			session.AbortTransaction(ctx)
		}
	}()

	defer tx.Rollback()

	if request.MongoDBEmployeeID != nil {
		mongoDBEmployeeID := *request.MongoDBEmployeeID
		filter := bson.M{"employeeId": mongoDBEmployeeID}

		result, err := employeeCollection.DeleteOne(context.Background(), filter)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		if result.DeletedCount == 0 {
			tx.Rollback()
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Employee not found in MongoDB or SQL server",
			})
		}
	}

	if request.SQLEmployeeID != nil {
		sqlEmployeeID := *request.SQLEmployeeID

		result, err := tx.ExecContext(ctx, "DELETE FROM Personal WHERE Employee_ID = @p1", sql.Named("p1", sqlEmployeeID))
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		rowCount, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get rows affected",
			})
		}

		if rowCount == 0 {
			tx.Rollback()
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Employee not found in SQL database",
			})
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to commit SQL Server transaction",
		})
	}

	err = pusherClient.Trigger("GoSIS", "both-deleted", request)
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Delete successful",
	})
}
