package controllers

import (
	"context"
	"net/http"
	"time"

	configs "GoSIS/config"
	"GoSIS/models"
	"GoSIS/responses"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var employeeCollection *mongo.Collection = configs.GetCollection(configs.MongoDB, "employees")

func GetAllEmployees(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var employees []models.Employee

	cursor, err := employeeCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var employee models.Employee
		if err := cursor.Decode(&employee); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
		}
		employees = append(employees, employee)
	}

	// Tạo một fiber.Map mới và điền dữ liệu vào đó
	dataMap := fiber.Map{"data": employees}

	return c.JSON(responses.EmployeeResponse{Status: http.StatusOK, Message: "success", Data: &dataMap})
}
