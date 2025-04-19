package routes

import (
	"github.com/SaharKhamseh/cinema-backend/controller"
	"github.com/SaharKhamseh/cinema-backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {
	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Replace with your frontend URL
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true, // Important for cookies
	}))

	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	app.Use(middleware.IsAuthentication)

	// Movie routes
	app.Post("/api/movies", middleware.IsAdmin, controller.CreateMovie)
	app.Get("/api/movies", controller.GetMovies)
	app.Get("/api/movies/:id", controller.GetMovie)
	app.Put("/api/movies/:id", middleware.IsAdmin, controller.UpdateMovie)
	app.Delete("/api/movies/:id", middleware.IsAdmin, controller.DeleteMovie)

	// Theater routes
	app.Post("/api/theaters", middleware.IsAdmin, controller.CreateTheater)
	app.Get("/api/theaters", controller.GetTheaters)
	app.Get("/api/theaters/:id", controller.GetTheater)

	// Screen routes
	app.Post("/api/screens", middleware.IsAdmin, controller.CreateScreen)
	app.Get("/api/screens/:id/seats", controller.GetScreenSeats)

	// ShowTime routes
	app.Post("/api/showtimes", middleware.IsAdmin, controller.CreateShowTime)
	app.Get("/api/showtimes", controller.GetShowTimes)
	app.Get("/api/showtimes/:id", controller.GetShowTime)
	app.Put("/api/showtimes/:id", middleware.IsAdmin, controller.UpdateShowTime)
	app.Delete("/api/showtimes/:id", middleware.IsAdmin, controller.DeleteShowTime)

	// Booking routes
	app.Post("/api/bookings", middleware.IsAuthentication, controller.CreateBooking)
	app.Get("/api/bookings", middleware.IsAuthentication, controller.GetUserBookings)
	app.Get("/api/bookings/:id", middleware.IsAuthentication, controller.GetBooking)
	app.Post("/api/bookings/:id/cancel", middleware.IsAuthentication, controller.CancelBooking)
}
