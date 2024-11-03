package handlers

import (
	"github.com/OPTIC7409/tutor-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type StudentHandler struct {
	DB *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{DB: db}
}

func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	var student models.Student
	if err := c.BodyParser(&student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	result := h.DB.Create(&student)
	if err := result.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create student"})
	}

	return c.Status(fiber.StatusCreated).JSON(student)
}

func (h *StudentHandler) GetStudents(c *fiber.Ctx) error {
	var students []models.Student
	result := h.DB.Preload("User").Find(&students)
	if err := result.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch students"})
	}

	return c.JSON(students)
}

func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	var student models.Student
	result := h.DB.Preload("User").First(&student, id)
	if err := result.Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	return c.JSON(student)
}

func (h *StudentHandler) UpdateStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	var student models.Student
	if err := h.DB.First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	if err := c.BodyParser(&student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	h.DB.Save(&student)
	return c.JSON(student)
}

func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	result := h.DB.Delete(&models.Student{}, id)
	if err := result.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete student"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
