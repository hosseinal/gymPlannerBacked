package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"yourusername/gym-planner/internal/auth"
	"yourusername/gym-planner/internal/database"
	"yourusername/gym-planner/internal/handlers"
	"yourusername/gym-planner/internal/middleware"
)

func main() {
	// Initialize database connection
	db, err := database.NewDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db.DB); err != nil {
		log.Fatal(err)
	}

	// Initialize JWT manager
	jwtMgr := auth.NewJWTManager(os.Getenv("JWT_SECRET"))

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db, jwtMgr)
	planHandler := handlers.NewPlanHandler(db)
	planDetailsHandler := handlers.NewPlanDetailsHandler(db)

	// Initialize router
	r := gin.Default()

	// Public routes
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(jwtMgr)
	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		// Plan routes
		protected.POST("/plans", planHandler.CreatePlan)
		protected.GET("/plans/list", planHandler.GetPlans)
		protected.GET("/plans/get", planHandler.GetPlan)
		protected.DELETE("/plans/delete", planHandler.DeletePlan)

		// Plan details routes
		protected.POST("/plan-details/add", planDetailsHandler.AddPlanDetail)
		protected.GET("/plan-details/get", planDetailsHandler.GetPlanDetails)
		protected.PUT("/plan-details/update", planDetailsHandler.UpdatePlanDetail)
		protected.DELETE("/plan-details/delete", planDetailsHandler.DeletePlanDetail)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
} 