package controllers

import (
	"context"
	"database/sql"
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
		err := rows.Scan(&p.EmployeeID, &p.FirstName, &p.LastName, &p.MiddleInitial, &p.Address1, &p.Address2, &p.City, &p.State, &p.Zip, &p.Email, &p.PhoneNumber, &p.SocialSecurityNumber, &p.DriversLicense, &p.MaritalStatus, &p.Gender, &p.ShareholderStatus, &p.BenefitPlans, &p.Ethnicity)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
		}
		personals = append(personals, p)
	}

	// Tạo một fiber.Map mới và điền dữ liệu vào đó
	dataMap := fiber.Map{"data": personals}

	return c.JSON(responses.PersonalResponse{Status: http.StatusOK, Message: "success", Data: &dataMap})
}

func CreatePersonal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize a new Personal object from the data sent from the client
	p := new(models.Personal)
	if err := c.BodyParser(p); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PersonalResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: nil})
	}

	// Add the data to the database
	_, err := sqlServerDB.ExecContext(ctx, "INSERT INTO Personal (Employee_ID, First_Name, Last_Name, Middle_Initial, Address1, Address2, City, State, Zip, Email, Phone_Number, Social_Security_Number, Drivers_License, Marital_Status, Gender, Shareholder_Status, Benefit_Plans, Ethnicity) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18)",
		sql.Named("p1", p.EmployeeID), sql.Named("p2", p.FirstName), sql.Named("p3", p.LastName), sql.Named("p4", p.MiddleInitial), sql.Named("p5", p.Address1), sql.Named("p6", p.Address2), sql.Named("p7", p.City), sql.Named("p8", p.State), sql.Named("p9", p.Zip), sql.Named("p10", p.Email), sql.Named("p11", p.PhoneNumber), sql.Named("p12", p.SocialSecurityNumber), sql.Named("p13", p.DriversLicense), sql.Named("p14", p.MaritalStatus), sql.Named("p15", p.Gender), sql.Named("p16", p.ShareholderStatus), sql.Named("p17", p.BenefitPlans), sql.Named("p18", p.Ethnicity))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil})
	}

	// Create a new fiber.Map and fill it with the data
	dataMap := fiber.Map{"personal": p}

	// Return the data that was added to the database
	return c.JSON(responses.PersonalResponse{Status: http.StatusOK, Message: "success", Data: &dataMap})
}
