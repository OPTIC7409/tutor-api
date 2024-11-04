package handlers

import (
	"github.com/OPTIC7409/tutor-api/internal/models"
	"github.com/OPTIC7409/tutor-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetDashboardData(c *fiber.Ctx) error {
	userID, err := utils.ExtractUserIDFromToken(c, h.DB)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user data",
		})
	}

	var dashboardData models.DashboardData
	dashboardData.ID = int(user.ID)
	dashboardData.Name = user.Name
	dashboardData.Email = user.Email
	dashboardData.UserType = user.UserType

	if user.UserType == "tutor" {
		tutorData, err := h.GetTutorDashboardData(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch tutor data",
			})
		}
		dashboardData.Stats = tutorData.Stats
		dashboardData.Requests = tutorData.Requests
	} else {
		studentData, err := h.GetStudentDashboardData(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch student data",
			})
		}
		dashboardData.UpcomingSessions = studentData.UpcomingSessions
	}

	chats, err := h.GetUserChats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user chats",
		})
	}
	dashboardData.Chats = chats

	return c.JSON(dashboardData)
}

func (h *UserHandler) GetTutorDashboardData(tutorID int) (*models.DashboardData, error) {
	var tutorData models.DashboardData

	// Fetch tutor stats
	var stats models.TutorStats
	if err := h.DB.Table("sessions").
		Select("COUNT(DISTINCT student_id) as active_students, COUNT(*) as upcoming_sessions, SUM(price) as earnings_this_month").
		Where("tutor_id = ? AND start_time > NOW() AND start_time < DATE_ADD(NOW(), INTERVAL 1 MONTH)", tutorID).
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	tutorData.Stats = &stats

	// Fetch tutor requests
	var requests []models.Request
	if err := h.DB.Table("requests").
		Select("requests.id, users.name as student, subjects.name as subject, requests.budget, users.avatar").
		Joins("JOIN users ON requests.student_id = users.id").
		Joins("JOIN subjects ON requests.subject_id = subjects.id").
		Where("requests.tutor_id = ? AND requests.status = ?", tutorID, "pending").
		Scan(&requests).Error; err != nil {
		return nil, err
	}
	tutorData.Requests = requests

	return &tutorData, nil
}

func (h *UserHandler) GetStudentDashboardData(studentID int) (*models.DashboardData, error) {
	var studentData models.DashboardData

	// Fetch upcoming sessions
	var upcomingSessions []models.UpcomingSession
	if err := h.DB.Table("sessions").
		Select("users.name as tutor, subjects.name as subject, sessions.start_time as datetime").
		Joins("JOIN users ON sessions.tutor_id = users.id").
		Joins("JOIN subjects ON sessions.subject_id = subjects.id").
		Where("sessions.student_id = ? AND sessions.start_time > NOW()", studentID).
		Order("sessions.start_time ASC").
		Limit(5).
		Scan(&upcomingSessions).Error; err != nil {
		return nil, err
	}
	studentData.UpcomingSessions = upcomingSessions

	return &studentData, nil
}

func (h *UserHandler) GetUserChats(userID int) ([]models.Chat, error) {
	var chats []models.Chat

	if err := h.DB.Table("chats").
		Select("chats.id, users.name, messages.content as last_message, users.avatar").
		Joins("JOIN users ON CASE WHEN chats.user1_id = ? THEN chats.user2_id = users.id ELSE chats.user1_id = users.id END", userID).
		Joins("LEFT JOIN messages ON chats.id = messages.chat_id AND messages.id = (SELECT MAX(id) FROM messages WHERE chat_id = chats.id)").
		Where("chats.user1_id = ? OR chats.user2_id = ?", userID, userID).
		Scan(&chats).Error; err != nil {
		return nil, err
	}

	return chats, nil
}
