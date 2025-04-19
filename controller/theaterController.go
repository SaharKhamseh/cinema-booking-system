package controller

import (
	"github.com/SaharKhamseh/cinema-backend/database"
	"github.com/SaharKhamseh/cinema-backend/models"
	"github.com/gofiber/fiber/v2"
)

// CreateTheater creates a new theater with screens
func CreateTheater(c *fiber.Ctx) error {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if data["name"] == nil || data["capacity"] == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Name and capacity are required",
		})
	}

	theater := models.Theater{
		Name:     data["name"].(string),
		Capacity: int(data["capacity"].(float64)),
	}

	if err := database.DB.Create(&theater).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create theater",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Theater created successfully",
		"theater": theater,
	})
}

// CreateScreen adds a new screen to a theater
func CreateScreen(c *fiber.Ctx) error {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	requiredFields := []string{"theater_id", "name", "capacity"}
	for _, field := range requiredFields {
		if data[field] == nil {
			return c.Status(400).JSON(fiber.Map{
				"message": field + " is required",
			})
		}
	}

	// Check if theater exists
	var theater models.Theater
	if err := database.DB.First(&theater, data["theater_id"]).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Theater not found",
		})
	}

	screen := models.Screen{
		TheaterID: uint(data["theater_id"].(float64)),
		Name:      data["name"].(string),
		Capacity:  int(data["capacity"].(float64)),
	}

	if err := database.DB.Create(&screen).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create screen",
			"error":   err.Error(),
		})
	}

	// Create seats for the screen
	if err := createSeatsForScreen(&screen); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create seats",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Screen created successfully",
		"screen":  screen,
	})
}

// Helper function to create seats for a screen
func createSeatsForScreen(screen *models.Screen) error {
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	seatsPerRow := screen.Capacity / len(rows)

	for _, row := range rows {
		for i := 1; i <= seatsPerRow; i++ {
			seat := models.Seat{
				ScreenID: screen.ID,
				Row:      row,
				Number:   i,
				Category: "standard", // You can modify this based on row position
			}
			if err := database.DB.Create(&seat).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// GetTheaters returns all theaters with their screens
func GetTheaters(c *fiber.Ctx) error {
	var theaters []models.Theater

	if err := database.DB.Preload("Screens").Find(&theaters).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch theaters",
			"error":   err.Error(),
		})
	}

	return c.JSON(theaters)
}

// GetTheater returns a specific theater with its screens
func GetTheater(c *fiber.Ctx) error {
	id := c.Params("id")
	var theater models.Theater

	if err := database.DB.Preload("Screens").First(&theater, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Theater not found",
		})
	}

	return c.JSON(theater)
}

// GetScreenSeats returns all seats for a specific screen
func GetScreenSeats(c *fiber.Ctx) error {
	screenID := c.Params("id")
	var seats []models.Seat

	if err := database.DB.Where("screen_id = ?", screenID).Find(&seats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch seats",
			"error":   err.Error(),
		})
	}

	return c.JSON(seats)
}
