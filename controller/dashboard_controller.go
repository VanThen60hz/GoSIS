package controllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"GoSIS/models"
	"GoSIS/responses"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GenderRatio(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	personalsMap := make(map[string]models.Personal)
	employeesMap := make(map[string]models.Employee)

	rows, _ := sqlServerDB.QueryContext(ctx, "SELECT * FROM Personal")
	for rows.Next() {
		var p models.Personal
		rows.Scan(&p.EmployeeID, &p.FirstName, &p.LastName, &p.MiddleInitial, &p.Address1, &p.Address2, &p.City, &p.State, &p.Zip, &p.Email, &p.PhoneNumber, &p.SocialSecurityNumber, &p.DriversLicense, &p.MaritalStatus, &p.Gender, &p.ShareholderStatus, &p.BenefitPlans, &p.Ethnicity)
		personalsMap[p.FirstName+p.LastName] = p
	}
	rows.Close()

	cursor, _ := employeeCollection.Find(ctx, bson.M{})
	for cursor.Next(ctx) {
		var e models.Employee
		cursor.Decode(&e)
		employeesMap[e.FirstName+e.LastName] = e
	}
	cursor.Close(ctx)

	var mergedData []models.MergePerson
	for key, p := range personalsMap {
		if e, ok := employeesMap[key]; ok {
			mergedPerson := models.MergePerson{
				EmployeeID:           e.EmployeeID + strconv.Itoa(int(p.EmployeeID)),
				FirstName:            e.FirstName,
				LastName:             e.LastName,
				VacationDays:         e.VacationDays,
				PaidToDate:           e.PaidToDate,
				PaidLastYear:         e.PaidLastYear,
				PayRate:              e.PayRate,
				PayRateID:            e.PayRateID,
				MiddleInitial:        p.MiddleInitial,
				Address1:             p.Address1,
				Address2:             p.Address2,
				City:                 p.City,
				State:                p.State,
				Zip:                  p.Zip,
				Email:                p.Email,
				PhoneNumber:          p.PhoneNumber,
				SocialSecurityNumber: p.SocialSecurityNumber,
				DriversLicense:       p.DriversLicense,
				MaritalStatus:        p.MaritalStatus,
				Gender:               p.Gender,
				ShareholderStatus:    p.ShareholderStatus,
				BenefitPlans:         p.BenefitPlans,
				Ethnicity:            p.Ethnicity,
			}
			mergedData = append(mergedData, mergedPerson)
		} else {
			mergedPerson := models.MergePerson{
				EmployeeID:           strconv.Itoa(int(p.EmployeeID)),
				FirstName:            p.FirstName,
				LastName:             p.LastName,
				VacationDays:         0,
				PaidToDate:           0,
				PaidLastYear:         0,
				PayRate:              0,
				PayRateID:            0,
				MiddleInitial:        p.MiddleInitial,
				Address1:             p.Address1,
				Address2:             p.Address2,
				City:                 p.City,
				State:                p.State,
				Zip:                  p.Zip,
				Email:                p.Email,
				PhoneNumber:          p.PhoneNumber,
				SocialSecurityNumber: p.SocialSecurityNumber,
				DriversLicense:       p.DriversLicense,
				MaritalStatus:        p.MaritalStatus,
				Gender:               p.Gender,
				ShareholderStatus:    p.ShareholderStatus,
				BenefitPlans:         p.BenefitPlans,
				Ethnicity:            p.Ethnicity,
			}
			mergedData = append(mergedData, mergedPerson)
		}
	}

	for key, e := range employeesMap {
		if _, ok := personalsMap[key]; !ok {
			mergedPerson := models.MergePerson{
				EmployeeID:           e.EmployeeID,
				FirstName:            e.FirstName,
				LastName:             e.LastName,
				VacationDays:         e.VacationDays,
				PaidToDate:           e.PaidToDate,
				PaidLastYear:         e.PaidLastYear,
				PayRate:              e.PayRate,
				PayRateID:            e.PayRateID,
				MiddleInitial:        "",
				Address1:             "",
				Address2:             "",
				City:                 "",
				State:                "",
				Zip:                  0,
				Email:                "",
				PhoneNumber:          "",
				SocialSecurityNumber: "",
				DriversLicense:       "",
				MaritalStatus:        "",
				Gender:               false,
				ShareholderStatus:    false,
				BenefitPlans:         0,
				Ethnicity:            "",
			}
			mergedData = append(mergedData, mergedPerson)
		}
	}

	maleCount, femaleCount := 0.0, 0.0
	for _, person := range mergedData {
		if person.Gender {
			maleCount++
		} else {
			femaleCount++
		}
	}

	maleRatio := maleCount * 100 / (maleCount + femaleCount)
	femaleRatio := femaleCount * 100 / (maleCount + femaleCount)

	dataMap := fiber.Map{
		"total":  len(mergedData),
		"male":   math.Round(maleRatio),
		"female": math.Round(femaleRatio),
	}
	return c.JSON(responses.GenderRatioResponse{Status: http.StatusOK, Message: "success", Data: &dataMap})
}
