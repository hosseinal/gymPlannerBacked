package handlers

import (
	"github.com/gin-gonic/gin"
	"yourusername/gym-planner/internal/database"
	"yourusername/gym-planner/internal/models"
)

type PlanHandler struct {
	db *database.DB
}

func NewPlanHandler(db *database.DB) *PlanHandler {
	return &PlanHandler{db: db}
}

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var plan models.Plan
	plan.UserID = userID

	var id int64
	err := h.db.QueryRow(
		"INSERT INTO plans (user_id) VALUES ($1) RETURNING id",
		plan.UserID,
	).Scan(&id)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating plan"})
		return
	}

	c.JSON(201, gin.H{"id": id})
}

func (h *PlanHandler) GetPlans(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var plans []models.Plan
	err := h.db.Select(&plans, "SELECT * FROM plans WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching plans"})
		return
	}

	c.JSON(200, plans)
}

func (h *PlanHandler) GetPlan(c *gin.Context) {
	userID := c.GetInt64("user_id")
	planID := c.Query("id")
	if planID == "" {
		c.JSON(400, gin.H{"error": "Plan ID is required"})
		return
	}

	var plan models.Plan
	err := h.db.Get(&plan, "SELECT * FROM plans WHERE id = $1 AND user_id = $2", planID, userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Plan not found"})
		return
	}

	c.JSON(200, plan)
}

func (h *PlanHandler) DeletePlan(c *gin.Context) {
	userID := c.GetInt64("user_id")
	planID := c.Query("id")
	if planID == "" {
		c.JSON(400, gin.H{"error": "Plan ID is required"})
		return
	}

	result, err := h.db.Exec("DELETE FROM plans WHERE id = $1 AND user_id = $2", planID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting plan"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error checking deletion"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Plan not found"})
		return
	}

	c.Status(200)
} 