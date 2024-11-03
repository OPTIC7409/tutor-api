package tests

import (
	"github.com/OPTIC7409/tutor-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}, &models.Tutor{}, &models.Student{}, &models.Chat{}, &models.Message{}); err != nil {
		return err
	}

	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Password: "password123", UserType: "student"},
		{Name: "Jane Smith", Email: "jane@example.com", Password: "password456", UserType: "tutor"},
	}

	for i := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		users[i].Password = string(hashedPassword)

		if err := db.FirstOrCreate(&users[i], models.User{Email: users[i].Email}).Error; err != nil {
			return err
		}
	}

	var tutorUser models.User
	if err := db.Where("email = ?", "jane@example.com").First(&tutorUser).Error; err != nil {
		return err
	}

	tutor := models.Tutor{UserID: tutorUser.ID, Subject: "Mathematics", YearsExperience: 5, HourlyRate: 50, Location: "New York"}
	if err := db.FirstOrCreate(&tutor, models.Tutor{UserID: tutorUser.ID}).Error; err != nil {
		return err
	}

	chat := models.Chat{}
	if err := db.FirstOrCreate(&chat).Error; err != nil {
		return err
	}

	for _, user := range users {
		if err := db.Model(&chat).Association("Participants").Append(&user); err != nil {
			return err
		}
	}

	messages := []models.Message{
		{ChatID: chat.ID, SenderID: users[0].ID, Content: "Hello, I need help with math."},
		{ChatID: chat.ID, SenderID: users[1].ID, Content: "Sure, I'd be happy to help. What topic are you struggling with?"},
	}

	for _, message := range messages {
		if err := db.FirstOrCreate(&message, models.Message{ChatID: message.ChatID, SenderID: message.SenderID, Content: message.Content}).Error; err != nil {
			return err
		}
	}

	return nil
}
