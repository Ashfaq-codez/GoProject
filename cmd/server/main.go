package main

import (
	"database/sql"
	"log"
	"os" // <-- 1. Add the os package

	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/middleware"
	"user-api/internal/repository"
	"user-api/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	// 2. Read the database URL from Docker, with a fallback for local testing
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Your local Arch database credentials
		dsn = "postgres://ash:274327@localhost:5432/ashdb?sslmode=disable"
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	// ... the rest of the file stays exactly the same ...
	defer dbConn.Close()

	// 3. Initialize Repository & Handlers
	userRepo := repository.NewUserRepository(dbConn)
	userHandler := handler.NewUserHandler(userRepo)

	// 4. Setup Fiber App
	app := fiber.New()

	// 2. Add the CORS middleware right after creating the app
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allows any frontend to connect
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// 5. Apply Middleware
	app.Use(middleware.RequestLogger())

	// 6. Setup Routes
	routes.SetupRoutes(app, userHandler)

	// 7. Start Server
	logger.Log.Info("Server starting on port 3000")
	if err := app.Listen(":3000"); err != nil {
		logger.Log.Fatal("Server failed to start", zap.Error(err))
	}
}
