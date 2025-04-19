package controller

import (
	"time"

	"github.com/SaharKhamseh/cinema-backend/database"
	"github.com/SaharKhamseh/cinema-backend/models"
	"github.com/gofiber/fiber/v2"
)

// CreateMovie creates a new movie
func CreateMovie(c *fiber.Ctx) error {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	requiredFields := []string{"title", "duration", "language"}
	for _, field := range requiredFields {
		if data[field] == nil {
			return c.Status(400).JSON(fiber.Map{
				"message": field + " is required",
			})
		}
	}

	// Parse release date
	releaseDate, err := time.Parse("2006-01-02", data["release_date"].(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid release date format. Use YYYY-MM-DD",
		})
	}

	movie := models.Movie{
		Title:       data["title"].(string),
		Description: data["description"].(string),
		Duration:    int(data["duration"].(float64)),
		Genre:       data["genre"].(string),
		Language:    data["language"].(string),
		ReleaseDate: releaseDate,
		PosterURL:   data["poster_url"].(string),
	}

	if err := database.DB.Create(&movie).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create movie",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Movie created successfully",
		"movie":   movie,
	})
}

// GetMovies returns all movies
func GetMovies(c *fiber.Ctx) error {
	var movies []models.Movie

	if err := database.DB.Find(&movies).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch movies",
			"error":   err.Error(),
		})
	}

	return c.JSON(movies)
}

// GetMovie returns a specific movie
func GetMovie(c *fiber.Ctx) error {
	id := c.Params("id")
	var movie models.Movie

	if err := database.DB.First(&movie, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Movie not found",
		})
	}

	return c.JSON(movie)
}

// UpdateMovie updates a movie
func UpdateMovie(c *fiber.Ctx) error {
	id := c.Params("id")
	var movie models.Movie

	if err := database.DB.First(&movie, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Movie not found",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Update fields if they exist in the request
	if data["title"] != nil {
		movie.Title = data["title"].(string)
	}
	if data["description"] != nil {
		movie.Description = data["description"].(string)
	}
	if data["duration"] != nil {
		movie.Duration = int(data["duration"].(float64))
	}
	if data["genre"] != nil {
		movie.Genre = data["genre"].(string)
	}
	if data["language"] != nil {
		movie.Language = data["language"].(string)
	}
	if data["release_date"] != nil {
		releaseDate, err := time.Parse("2006-01-02", data["release_date"].(string))
		if err == nil {
			movie.ReleaseDate = releaseDate
		}
	}
	if data["poster_url"] != nil {
		movie.PosterURL = data["poster_url"].(string)
	}

	database.DB.Save(&movie)

	return c.JSON(fiber.Map{
		"message": "Movie updated successfully",
		"movie":   movie,
	})
}

// DeleteMovie deletes a movie
func DeleteMovie(c *fiber.Ctx) error {
	id := c.Params("id")
	var movie models.Movie

	if err := database.DB.First(&movie, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Movie not found",
		})
	}

	database.DB.Delete(&movie)

	return c.JSON(fiber.Map{
		"message": "Movie deleted successfully",
	})
}
