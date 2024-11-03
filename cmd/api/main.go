package main

import (
	"log"
	"os"

	"github.com/OPTIC7409/tutor-api/config"
	"github.com/OPTIC7409/tutor-api/internal/database"
	"github.com/OPTIC7409/tutor-api/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	app := fiber.New()

	app.Use(cors.New())

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	authHandler := handlers.NewAuthHandler(db)
	tutorHandler := handlers.NewTutorHandler(db)
	studentHandler := handlers.NewStudentHandler(db)
	chatHandler := handlers.NewChatHandler(db)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	tutors := api.Group("/tutors")
	tutors.Post("/", tutorHandler.CreateTutor)
	tutors.Get("/", tutorHandler.GetTutors)
	tutors.Get("/:id", tutorHandler.GetTutor)
	tutors.Put("/:id", tutorHandler.UpdateTutor)
	tutors.Delete("/:id", tutorHandler.DeleteTutor)

	students := api.Group("/students")
	students.Post("/", studentHandler.CreateStudent)
	students.Get("/", studentHandler.GetStudents)
	students.Get("/:id", studentHandler.GetStudent)
	students.Put("/:id", studentHandler.UpdateStudent)
	students.Delete("/:id", studentHandler.DeleteStudent)

	chats := api.Group("/chats")
	chats.Get("/", chatHandler.GetChats)
	chats.Get("/:id", chatHandler.GetChat)
	chats.Post("/", chatHandler.CreateChat)
	chats.Post("/:id/messages", chatHandler.SendMessage)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
