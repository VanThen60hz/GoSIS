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
)

var sqlServerDB *sql.DB = configs.ConnectSqlServerDB()

func GetAllPersonals(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var personals []models.Personal

	rows, err := sqlServerDB.QueryContext(ctx, "SELECT * FROM Personal")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Personal
		err := rows.Scan(&p.SQLEmployeeId, &p.FirstName, &p.LastName, &p.MiddleInitial, &p.Address1, &p.Address2, &p.City, &p.State, &p.Zip, &p.Email, &p.PhoneNumber, &p.SocialSecurityNumber, &p.DriversLicense, &p.MaritalStatus, &p.Gender, &p.ShareholderStatus, &p.BenefitPlans, &p.Ethnicity)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
		}
		personals = append(personals, p)
	}

	return c.JSON(responses.PersonalResponse{Status: http.StatusOK, Message: "success", Data: &personals})
}

func CreatePersonal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p := new(models.Personal)
	if err := c.BodyParser(p); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PersonalResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: nil})
	}

	var lastInsertID int64
	err := sqlServerDB.QueryRowContext(ctx, "SELECT TOP 1 Employee_ID FROM Personal ORDER BY Employee_ID DESC").Scan(&lastInsertID)
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

	return c.JSON(responses.CreatePersonalResponse{Status: http.StatusOK, Message: "Create personal successfully", Data: p})
}

func UpdatePersonal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p := new(models.Personal)
	if err := c.BodyParser(p); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PersonalResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: nil})
	}

	_, err := sqlServerDB.ExecContext(ctx, `
        UPDATE Personal 
        SET First_Name = @p2, Last_Name = @p3, Middle_Initial = @p4, Address1 = @p5, Address2 = @p6, 
            City = @p7, State = @p8, Zip = @p9, Email = @p10, Phone_Number = @p11, 
            Social_Security_Number = @p12, Drivers_License = @p13, Marital_Status = @p14, 
            Gender = @p15, Shareholder_Status = @p16, Benefit_Plans = @p17, Ethnicity = @p18
        WHERE Employee_ID = @p1`,
		sql.Named("p1", p.SQLEmployeeId), sql.Named("p2", p.FirstName), sql.Named("p3", p.LastName), sql.Named("p4", p.MiddleInitial),
		sql.Named("p5", p.Address1), sql.Named("p6", p.Address2), sql.Named("p7", p.City), sql.Named("p8", p.State),
		sql.Named("p9", p.Zip), sql.Named("p10", p.Email), sql.Named("p11", p.PhoneNumber), sql.Named("p12", p.SocialSecurityNumber),
		sql.Named("p13", p.DriversLicense), sql.Named("p14", p.MaritalStatus), sql.Named("p15", p.Gender),
		sql.Named("p16", p.ShareholderStatus), sql.Named("p17", p.BenefitPlans), sql.Named("p18", p.Ethnicity))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil})
	}

	return c.JSON(responses.UpdatePersonalResponse{Status: http.StatusOK, Message: "Update personal successfully", Data: p})
}

func fetchPersonals(ctx context.Context) (map[string]models.Personal, int, error) {
	rows, err := sqlServerDB.QueryContext(ctx, "SELECT * FROM Personal ORDER BY First_Name")
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching Personal data: %w", err)
	}
	defer rows.Close()

	personalsMap := make(map[string]models.Personal)
	var totalCount int

	for rows.Next() {
		var p models.Personal
		rows.Scan(&p.SQLEmployeeId, &p.FirstName, &p.LastName, &p.MiddleInitial, &p.Address1, &p.Address2, &p.City, &p.State, &p.Zip, &p.Email, &p.PhoneNumber, &p.SocialSecurityNumber, &p.DriversLicense, &p.MaritalStatus, &p.Gender, &p.ShareholderStatus, &p.BenefitPlans, &p.Ethnicity)
		personalsMap[p.FirstName+p.LastName] = p
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error fetching Personal data: %w", err)
	}

	err = sqlServerDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM Personal").Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching Personal data count: %w", err)
	}

	return personalsMap, totalCount, nil
}

func deletePersonal(personalID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := sqlServerDB.ExecContext(ctx, "DELETE FROM Personal WHERE Employee_ID = @p1", sql.Named("p1", personalID))
	if err != nil {
		return fmt.Errorf("error deleting Personal: %w", err)
	}

	return nil
}
