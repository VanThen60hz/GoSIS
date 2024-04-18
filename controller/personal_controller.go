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
		err := rows.Scan(&p.EmployeeID, &p.FirstName, &p.LastName, &p.MiddleInitial, &p.Address1, &p.Address2, &p.City, &p.State, &p.Zip, &p.Email, &p.PhoneNumber, &p.SocialSecurityNumber, &p.DriversLicense, &p.MaritalStatus, &p.Gender, &p.ShareholderStatus, &p.BenefitPlans, &p.Ethnicity)
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

	_, err := sqlServerDB.ExecContext(ctx, "INSERT INTO Personal (Employee_ID, First_Name, Last_Name, Middle_Initial, Address1, Address2, City, State, Zip, Email, Phone_Number, Social_Security_Number, Drivers_License, Marital_Status, Gender, Shareholder_Status, Benefit_Plans, Ethnicity) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18)",
		sql.Named("p1", p.EmployeeID), sql.Named("p2", p.FirstName), sql.Named("p3", p.LastName), sql.Named("p4", p.MiddleInitial), sql.Named("p5", p.Address1), sql.Named("p6", p.Address2), sql.Named("p7", p.City), sql.Named("p8", p.State), sql.Named("p9", p.Zip), sql.Named("p10", p.Email), sql.Named("p11", p.PhoneNumber), sql.Named("p12", p.SocialSecurityNumber), sql.Named("p13", p.DriversLicense), sql.Named("p14", p.MaritalStatus), sql.Named("p15", p.Gender), sql.Named("p16", p.ShareholderStatus), sql.Named("p17", p.BenefitPlans), sql.Named("p18", p.Ethnicity))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PersonalResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil})
	}

	// // Tạo một map[string]interface{} để chứa thông tin cá nhân
	// personalData := map[string]interface{}{
	// 	"employeeID":           p.EmployeeID,
	// 	"firstName":            p.FirstName,
	// 	"lastName":             p.LastName,
	// 	"middleInitial":        p.MiddleInitial,
	// 	"address1":             p.Address1,
	// 	"address2":             p.Address2,
	// 	"city":                 p.City,
	// 	"state":                p.State,
	// 	"zip":                  p.Zip,
	// 	"email":                p.Email,
	// 	"phoneNumber":          p.PhoneNumber,
	// 	"socialSecurityNumber": p.SocialSecurityNumber,
	// 	"driversLicense":       p.DriversLicense,
	// 	"maritalStatus":        p.MaritalStatus,
	// 	"gender":               p.Gender,
	// 	"shareholderStatus":    p.ShareholderStatus,
	// 	"benefitPlans":         p.BenefitPlans,
	// 	"ethnicity":            p.Ethnicity,
	// }

	// // Sử dụng Pusher để gửi dữ liệu
	// pusherClient := pusher.Client{
	// 	AppID:   "1790030",
	// 	Key:     "a359a59a30b4ddb07bb5",
	// 	Secret:  "7010eecf86f9469246bf",
	// 	Cluster: "ap1",
	// 	Secure:  true,
	// }
	// err = pusherClient.Trigger("GoSIS", "personal-created", personalData)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	return c.JSON(responses.CreatePersonalReponse{Status: http.StatusOK, Message: "Create personal successfully", Data: p})
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
		rows.Scan(&p.EmployeeID, &p.FirstName, &p.LastName, &p.MiddleInitial, &p.Address1, &p.Address2, &p.City, &p.State, &p.Zip, &p.Email, &p.PhoneNumber, &p.SocialSecurityNumber, &p.DriversLicense, &p.MaritalStatus, &p.Gender, &p.ShareholderStatus, &p.BenefitPlans, &p.Ethnicity)
		personalsMap[p.FirstName+p.LastName] = p
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error fetching Personal data: %w", err)
	}

	// Get total count of Personal records
	err = sqlServerDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM Personal").Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching Personal data count: %w", err)
	}

	return personalsMap, totalCount, nil
}
