package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"GoSIS/models"
	"GoSIS/responses"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MergeData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Đọc pageNumber từ query string của URL
	pageNumber, err := strconv.Atoi(c.Query("pageNumber"))
	if err != nil {
		c.JSON(responses.MergeResponse{Status: http.StatusBadRequest, Message: "Error in param pageNumber", Data: nil})
	}

	pageSize := 10

	personalsMap, totalPersonal, _ := fetchPersonals(ctx, pageNumber, pageSize)
	employeesMap := make(map[string]models.Employee)

	sort := bson.M{"FirstName": 1}
	cursor, _ := employeeCollection.Find(ctx, bson.M{}, options.Find().SetSort(sort).SetSkip(int64(pageNumber)).SetLimit(int64(pageSize)))
	totalEmployee, _ := employeeCollection.CountDocuments(context.TODO(), bson.D{})
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

	return c.JSON(responses.MergeResponse{Status: http.StatusOK, Message: "success", Data: &mergedData, TotalSize: totalPersonal + int(totalEmployee)})
}
