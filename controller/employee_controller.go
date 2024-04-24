package controllers

import (
	"context"
	"net/http"
	"time"

	configs "GoSIS/config"
	"GoSIS/models"
	"GoSIS/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		return c.Status(http.StatusNotFound).JSON(responses.EmployeeResponse{Status: http.StatusNotFound, Message: "error", Data: nil})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var employee models.Employee
		if err := cursor.Decode(&employee); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
		}
		employees = append(employees, employee)
	}

	return c.JSON(responses.EmployeeResponse{Status: http.StatusOK, Message: "success", Data: &employees})
}

func CreateEmployee(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var employee models.EmployeeNotID
	if err := c.BodyParser(&employee); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "invalid employee data", Data: nil})
	}

	// Generate a new UUID for EmployeeID
	employee.EmployeeId = uuid.New().String()

	// Set timestamps for creation and update
	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	// Insert the employee document into the database
	_, err := employeeCollection.InsertOne(ctx, employee)
	if err != nil {
		// Handle potential duplicate EmployeeID errors
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(http.StatusConflict).JSON(responses.EmployeeResponse{Status: http.StatusConflict, Message: err.Error(), Data: nil})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error creating employee", Data: nil})
	}

	return c.Status(http.StatusCreated).JSON(responses.CreateEmployeeResponse{Status: http.StatusCreated, Message: "employee created successfully", Data: &employee})
}
