package responses

import (
	"GoSIS/models"

	"github.com/gofiber/fiber/v2"
)

type EmployeeResponse struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    *[]models.Employee `json:"data"`
}
type CreateEmployeeResponse struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    *models.EmployeeNotID `json:"data"`
}

type PersonalResponse struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    *[]models.Personal `json:"data"`
}

type CreatePersonalResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    *models.Personal `json:"data"`
}

type MergeResponse struct {
	Status    int                   `json:"status"`
	Message   string                `json:"message"`
	Data      *[]models.MergePerson `json:"data"`
	TotalSize int                   `json:"total_size"`
}

type GenderRatioResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

type CreateBothResponse struct {
	Status  int                          `json:"status"`
	Message string                       `json:"message"`
	Data    *models.MergePersonWithoutId `json:"data"`
}

type UpdateEmployeeResponse struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    *models.EmployeeNotID `json:"data"`
}

type UpdatePersonalResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    *models.Personal `json:"data"`
}

type UpdateBothResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    *models.MergePerson `json:"data"`
}

type DeleteBothResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
