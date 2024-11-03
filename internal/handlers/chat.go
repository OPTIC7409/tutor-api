package handlers

import (
	"github.com/OPTIC7409/tutor-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChatHandler struct {
	DB *gorm.DB
}

func NewChatHandler(db *gorm.DB) *ChatHandler {
	return &ChatHandler{DB: db}
}

func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	var chats []models.Chat
	if err := h.DB.Preload("Participants").Find(&chats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch chats",
		})
	}
	return c.JSON(chats)
}

func (h *ChatHandler) GetChat(c *fiber.Ctx) error {
	id := c.Params("id")
	var chat models.Chat
	if err := h.DB.Preload("Participants").Preload("Messages.Sender").First(&chat, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Chat not found",
		})
	}
	return c.JSON(chat)
}

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	var input struct {
		Participants []uint `json:"participants"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if len(input.Participants) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least two participants are required",
		})
	}

	var users []models.User
	if err := h.DB.Where("id IN ?", input.Participants).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	if len(users) != len(input.Participants) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "One or more users not found",
		})
	}

	chat := models.Chat{
		Participants: users,
	}

	if err := h.DB.Create(&chat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create chat",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(chat)
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	chatID := c.Params("id")
	var input struct {
		SenderID uint   `json:"senderID"`
		Content  string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	var chat models.Chat
	if err := h.DB.First(&chat, chatID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Chat not found",
		})
	}

	message := models.Message{
		ChatID:   chat.ID,
		SenderID: input.SenderID,
		Content:  input.Content,
	}

	if err := h.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send message",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(message)
}
