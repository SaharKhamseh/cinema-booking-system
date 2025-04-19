package controller

import (
	"time"

	"github.com/SaharKhamseh/cinema-backend/database"
	"github.com/SaharKhamseh/cinema-backend/models"
	"github.com/gofiber/fiber/v2"
)

// CreateShowTime creates a new show timing
func CreateShowTime(c *fiber.Ctx) error {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	requiredFields := []string{"movie_id", "screen_id", "start_time", "price"}
	for _, field := range requiredFields {
		if data[field] == nil {
			return c.Status(400).JSON(fiber.Map{
				"message": field + " is required",
			})
		}
	}

	// Parse start time
	startTime, err := time.Parse("2006-01-02 15:04:05", data["start_time"].(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid start time format. Use YYYY-MM-DD HH:MM:SS",
		})
	}

	// Verify movie exists
	var movie models.Movie
	if err := database.DB.First(&movie, data["movie_id"]).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Movie not found",
		})
	}

	// Calculate end time based on movie duration
	endTime := startTime.Add(time.Minute * time.Duration(movie.Duration))

	// Check for time conflicts
	var conflictingShows int64
	database.DB.Model(&models.ShowTime{}).
		Where("screen_id = ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?))",
			data["screen_id"], startTime, endTime, startTime, endTime).
		Count(&conflictingShows)

	if conflictingShows > 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Time slot conflicts with existing show",
		})
	}

	showTime := models.ShowTime{
		MovieID:   uint(data["movie_id"].(float64)),
		ScreenID:  uint(data["screen_id"].(float64)),
		StartTime: startTime,
		EndTime:   endTime,
		Price:     data["price"].(float64),
	}

	if err := database.DB.Create(&showTime).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create show time",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":  "Show time created successfully",
		"showtime": showTime,
	})
}

// GetShowTimes returns all show times for a specific date
func GetShowTimes(c *fiber.Ctx) error {
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	startOfDay, _ := time.Parse("2006-01-02", date)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var showTimes []models.ShowTime
	if err := database.DB.
		Preload("Movie").
		Preload("Screen").
		Where("start_time BETWEEN ? AND ?", startOfDay, endOfDay).
		Find(&showTimes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch show times",
			"error":   err.Error(),
		})
	}

	return c.JSON(showTimes)
}

// GetShowTime returns a specific show time
func GetShowTime(c *fiber.Ctx) error {
	id := c.Params("id")
	var showTime models.ShowTime

	if err := database.DB.
		Preload("Movie").
		Preload("Screen").
		First(&showTime, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Show time not found",
		})
	}

	return c.JSON(showTime)
}

// UpdateShowTime updates a show time
func UpdateShowTime(c *fiber.Ctx) error {
	id := c.Params("id")
	var showTime models.ShowTime

	if err := database.DB.First(&showTime, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Show time not found",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if data["start_time"] != nil {
		startTime, err := time.Parse("2006-01-02 15:04:05", data["start_time"].(string))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid start time format",
			})
		}
		showTime.StartTime = startTime

		// Recalculate end time
		var movie models.Movie
		database.DB.First(&movie, showTime.MovieID)
		showTime.EndTime = startTime.Add(time.Minute * time.Duration(movie.Duration))
	}

	if data["price"] != nil {
		showTime.Price = data["price"].(float64)
	}

	if err := database.DB.Save(&showTime).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update show time",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Show time updated successfully",
		"showtime": showTime,
	})
}

// DeleteShowTime deletes a show time
func DeleteShowTime(c *fiber.Ctx) error {
	id := c.Params("id")
	var showTime models.ShowTime

	if err := database.DB.First(&showTime, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Show time not found",
		})
	}

	if err := database.DB.Delete(&showTime).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to delete show time",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Show time deleted successfully",
	})
}
