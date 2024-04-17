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

	personalsMap, totalPersonal, err := fetchPersonals(ctx)
	if err != nil {
		return c.JSON(responses.MergeResponse{Status: http.StatusInternalServerError, Message: "Error fetching personals", Data: nil})
	}

	employeesMap := make(map[string]models.Employee)

	sort := bson.M{"FirstName": 1}
	cursor, err := employeeCollection.Find(ctx, bson.M{}, options.Find().SetSort(sort))
	if err != nil {
		return c.JSON(responses.MergeResponse{Status: http.StatusInternalServerError, Message: "Error fetching employees", Data: nil})
	}
	defer cursor.Close(ctx)

	totalEmployee, _ := employeeCollection.CountDocuments(context.TODO(), bson.D{})

	for cursor.Next(ctx) {
		var e models.Employee
		err := cursor.Decode(&e)
		if err != nil {
			return c.JSON(responses.MergeResponse{Status: http.StatusInternalServerError, Message: "Error decoding employee", Data: nil})
		}
		employeesMap[e.FirstName+e.LastName] = e
	}

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
				Zip:                  0,
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

	// Sort the mergedData slice by FirstName
	sortByFirstName(mergedData)

	// Đọc pageNumber từ query string của URL
	pageNumber, err := strconv.Atoi(c.Query("pageNumber"))
	if err != nil {
		return c.JSON(responses.MergeResponse{Status: http.StatusBadRequest, Message: "Error in param pageNumber", Data: nil})
	}

	// Đọc pageNumber từ query string của URL
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		return c.JSON(responses.MergeResponse{Status: http.StatusBadRequest, Message: "Error in param pageSize", Data: nil})
	}

	mergeRes := mergedData[pageNumber-1 : pageNumber-1+pageSize]

	return c.JSON(responses.MergeResponse{Status: http.StatusOK, Message: "success", Data: &mergeRes, TotalSize: totalPersonal + int(totalEmployee)})
}

func sortByFirstName(mergedData []models.MergePerson) {
	for i := 0; i < len(mergedData)-1; i++ {
		for j := i + 1; j < len(mergedData); j++ {
			if mergedData[i].FirstName > mergedData[j].FirstName {
				mergedData[i], mergedData[j] = mergedData[j], mergedData[i]
			}
		}
	}
}
