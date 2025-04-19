package controller

import (
	"strconv"
	"time"

	"github.com/SaharKhamseh/cinema-backend/database"
	"github.com/SaharKhamseh/cinema-backend/models"
	"github.com/SaharKhamseh/cinema-backend/util"
	"github.com/gofiber/fiber/v2"
)

// CreateBooking handles new ticket bookings
func CreateBooking(c *fiber.Ctx) error {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	requiredFields := []string{"show_time_id", "seat_ids"}
	for _, field := range requiredFields {
		if data[field] == nil {
			return c.Status(400).JSON(fiber.Map{
				"message": field + " is required",
			})
		}
	}

	// Get user ID from JWT token
	cookie := c.Cookies("jwt")
	userIdStr, err := util.Parsejwt(cookie)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Convert string userId to uint
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	// Verify showtime exists and is in the future
	var showTime models.ShowTime
	if err := database.DB.First(&showTime, data["show_time_id"]).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Show time not found",
		})
	}

	if showTime.StartTime.Before(time.Now()) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Cannot book tickets for past shows",
		})
	}

	// Convert seat_ids from interface{} to []uint
	seatIDsInterface := data["seat_ids"].([]interface{})
	seatIDs := make([]uint, len(seatIDsInterface))
	for i, id := range seatIDsInterface {
		seatIDs[i] = uint(id.(float64))
	}

	// Check if seats are available
	var count int64
	database.DB.Model(&models.Booking{}).
		Joins("JOIN booking_seats ON bookings.id = booking_seats.booking_id").
		Where("bookings.show_time_id = ? AND booking_seats.seat_id IN ? AND bookings.status != ?",
			showTime.ID, seatIDs, "cancelled").
		Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "One or more selected seats are already booked",
		})
	}

	// Calculate total price
	totalPrice := showTime.Price * float64(len(seatIDs))

	// Create booking
	booking := models.Booking{
		UserID:     uint(userId),
		ShowTimeID: showTime.ID,
		TotalPrice: totalPrice,
		Status:     "pending",
		BookedAt:   time.Now(),
	}

	// Start transaction
	tx := database.DB.Begin()

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create booking",
			"error":   err.Error(),
		})
	}

	// Add seats to booking
	for _, seatID := range seatIDs {
		if err := tx.Exec("INSERT INTO booking_seats (booking_id, seat_id) VALUES (?, ?)",
			booking.ID, seatID).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to assign seats",
				"error":   err.Error(),
			})
		}
	}

	tx.Commit()

	// Load the complete booking with relationships
	database.DB.Preload("User").Preload("ShowTime").Preload("Seats").First(&booking, booking.ID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Booking created successfully",
		"booking": booking,
	})
}

// GetUserBookings returns all bookings for the logged-in user
func GetUserBookings(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	userIdStr, err := util.Parsejwt(cookie)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Convert string userId to uint
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	var bookings []models.Booking
	if err := database.DB.
		Preload("ShowTime.Movie").
		Preload("ShowTime.Screen").
		Preload("Seats").
		Where("user_id = ?", uint(userId)).
		Order("booked_at DESC").
		Find(&bookings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch bookings",
			"error":   err.Error(),
		})
	}

	return c.JSON(bookings)
}

// GetBooking returns a specific booking
func GetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	var booking models.Booking

	if err := database.DB.
		Preload("ShowTime.Movie").
		Preload("ShowTime.Screen").
		Preload("Seats").
		First(&booking, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Booking not found",
		})
	}

	// Verify user owns this booking
	cookie := c.Cookies("jwt")
	userIdStr, _ := util.Parsejwt(cookie)
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	if booking.UserID != uint(userId) {
		return c.Status(403).JSON(fiber.Map{
			"message": "Unauthorized to view this booking",
		})
	}

	return c.JSON(booking)
}

// CancelBooking cancels a booking
func CancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	var booking models.Booking

	if err := database.DB.First(&booking, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Booking not found",
		})
	}

	// Verify user owns this booking
	cookie := c.Cookies("jwt")
	userIdStr, _ := util.Parsejwt(cookie)
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	if booking.UserID != uint(userId) {
		return c.Status(403).JSON(fiber.Map{
			"message": "Unauthorized to view this booking",
		})
	}

	// Check if show hasn't started yet
	if booking.ShowTime.StartTime.Before(time.Now()) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Cannot cancel booking for past shows",
		})
	}

	booking.Status = "cancelled"
	database.DB.Save(&booking)

	return c.JSON(fiber.Map{
		"message": "Booking cancelled successfully",
	})
}
