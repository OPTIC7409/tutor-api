package handlers

import (
	"github.com/OPTIC7409/tutor-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TutorHandler struct {
	DB *gorm.DB
}

func NewTutorHandler(db *gorm.DB) *TutorHandler {
	return &TutorHandler{DB: db}
}

func (h *TutorHandler) CreateTutor(c *fiber.Ctx) error {
	var tutor models.Tutor
	if err := c.BodyParser(&tutor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	result := h.DB.Create(&tutor)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tutor"})
	}

	return c.Status(fiber.StatusCreated).JSON(tutor)
}

func (h *TutorHandler) GetTutors(c *fiber.Ctx) error {
	var tutors []models.Tutor
	result := h.DB.Preload("User").Find(&tutors)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tutors"})
	}

	return c.JSON(tutors)
}

func (h *TutorHandler) GetTutor(c *fiber.Ctx) error {
	id := c.Params("id")
	var tutor models.Tutor
	result := h.DB.Preload("User").First(&tutor, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tutor not found"})
	}

	return c.JSON(tutor)
}

func (h *TutorHandler) UpdateTutor(c *fiber.Ctx) error {
	id := c.Params("id")
	var tutor models.Tutor
	if err := h.DB.First(&tutor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tutor not found"})
	}

	if err := c.BodyParser(&tutor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	h.DB.Save(&tutor)
	return c.JSON(tutor)
}

func (h *TutorHandler) DeleteTutor(c *fiber.Ctx) error {
	id := c.Params("id")
	result := h.DB.Delete(&models.Tutor{}, id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete tutor"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
