package responses

import (
	"GoSIS/models"
)

type EmployeeResponse struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    *[]models.Employee `json:"data"`
}

type PersonalResponse struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    *[]models.Personal `json:"data"`
}

type CreatePersonalReponse struct {
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
