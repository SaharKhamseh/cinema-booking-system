package controller

import (
	"fmt"
	// "log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SaharKhamseh/cinema-backend/database"
	"github.com/SaharKhamseh/cinema-backend/models"
	"github.com/SaharKhamseh/cinema-backend/util"
	"github.com/gofiber/fiber/v2"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}

	// Parse request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Check if all required fields are present and not empty
	requiredFields := []string{"first_name", "last_name", "email", "password", "phone"}
	for _, field := range requiredFields {
		if data[field] == nil {
			return c.Status(400).JSON(fiber.Map{
				"message": fmt.Sprintf("%s is required", strings.ReplaceAll(field, "_", " ")),
			})
		}

		// Convert to string and check if empty after trimming
		if str, ok := data[field].(string); !ok || strings.TrimSpace(str) == "" {
			return c.Status(400).JSON(fiber.Map{
				"message": fmt.Sprintf("%s cannot be empty", strings.ReplaceAll(field, "_", " ")),
			})
		}
	}

	// Check if password is less than 6 characters
	if len(strings.TrimSpace(data["password"].(string))) <= 6 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password must be greater than 6 characters",
		})
	}

	email := strings.TrimSpace(data["email"].(string))

	// Validate email format
	if !validateEmail(email) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid email address",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Create new user with trimmed values
	user := models.User{
		FirstName: strings.TrimSpace(data["first_name"].(string)),
		LastName:  strings.TrimSpace(data["last_name"].(string)),
		Phone:     strings.TrimSpace(data["phone"].(string)),
		Email:     email,
		Role:      "user",
	}

	user.SetPassword(strings.TrimSpace(data["password"].(string)))

	// Create user in database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error creating user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Email Address doesn't exist, create an account",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrrect password",
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "You have successfully login",
		"user":    user,
	})
}
