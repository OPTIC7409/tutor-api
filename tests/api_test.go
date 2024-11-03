package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/OPTIC7409/tutor-api/config"
	"github.com/OPTIC7409/tutor-api/internal/database"
	"gorm.io/gorm"
)

const (
	greenColor = "\033[32m"
	redColor   = "\033[31m"
	resetColor = "\033[0m"
)

type TestCase struct {
	Name     string
	Method   string
	URL      string
	Body     interface{}
	Expected int
}

var db *gorm.DB

func TestMain(m *testing.M) {
	fmt.Println("Setting up test environment...")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	db, err = database.InitDatabase(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	err = SeedDatabase(db)
	if err != nil {
		fmt.Printf("Failed to seed database: %v\n", err)
	} else {
		fmt.Println("Database seeded successfully")
	}

	fmt.Println("Running API tests...")
	code := m.Run()

	fmt.Println("Cleaning up test environment...")
	cleanupDatabase(db)
	sqlDB, _ := db.DB()
	sqlDB.Close()

	os.Exit(code)
}

func cleanupDatabase(db *gorm.DB) {
	db.Exec("DROP SCHEMA public CASCADE")
	db.Exec("CREATE SCHEMA public")
	db.Exec("GRANT ALL ON SCHEMA public TO postgres")
	db.Exec("GRANT ALL ON SCHEMA public TO public")
}

func findRootDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", fmt.Errorf("could not find project root")
		}
		currentDir = parentDir
	}
}

func TestAPI(t *testing.T) {
	baseURL := fmt.Sprintf("http://localhost:%s/api", os.Getenv("SERVER_PORT"))

	testCases := []TestCase{
		{
			Name:   "Register User",
			Method: "POST",
			URL:    baseURL + "/auth/register",
			Body: map[string]interface{}{
				"name":     "Test User",
				"email":    "johna@example.com",
				"password": "password123",
				"userType": "student",
			},
			Expected: http.StatusCreated,
		},
		{
			Name:   "Register Duplicate User",
			Method: "POST",
			URL:    baseURL + "/auth/register",
			Body: map[string]interface{}{
				"name":     "Test User",
				"email":    "john@example.com",
				"password": "password",
				"userType": "student",
			},
			Expected: http.StatusConflict,
		},
		{
			Name:   "Login User",
			Method: "POST",
			URL:    baseURL + "/auth/login",
			Body: map[string]interface{}{
				"email":    "johna@example.com",
				"password": "password",
			},
			Expected: http.StatusOK,
		},
		{
			Name:     "Get Tutors",
			Method:   "GET",
			URL:      baseURL + "/tutors",
			Expected: http.StatusOK,
		},
		{
			Name:     "Get Chats",
			Method:   "GET",
			URL:      baseURL + "/chats",
			Expected: http.StatusOK,
		},
		{
			Name:     "Get Chat",
			Method:   "GET",
			URL:      baseURL + "/chats/1",
			Expected: http.StatusOK,
		},
		{
			Name:   "Send Message",
			Method: "POST",
			URL:    baseURL + "/chats/1/messages",
			Body: map[string]interface{}{
				"senderID": 1,
				"content":  "Hello, this is a test message.",
			},
			Expected: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var resp *http.Response
			var err error

			switch tc.Method {
			case "GET":
				resp, err = http.Get(tc.URL)
			case "POST":
				jsonBody, _ := json.Marshal(tc.Body)
				resp, err = http.Post(tc.URL, "application/json", bytes.NewBuffer(jsonBody))
			default:
				t.Fatalf("Unsupported HTTP method: %s", tc.Method)
			}

			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			statusColor := redColor
			if resp.StatusCode == tc.Expected {
				statusColor = greenColor
			}

			fmt.Printf("%s■%s %s: %d\n", statusColor, resetColor, tc.URL, resp.StatusCode)

			if resp.StatusCode != tc.Expected {
				t.Errorf("Expected status %d, got %d. Response body: %s", tc.Expected, resp.StatusCode, string(body))
			}
		})
	}
}
