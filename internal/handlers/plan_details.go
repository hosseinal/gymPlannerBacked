package handlers

import (
	"github.com/gin-gonic/gin"
	"yourusername/gym-planner/internal/database"
	"yourusername/gym-planner/internal/models"
)

type PlanDetailsHandler struct {
	db *database.DB
}

func NewPlanDetailsHandler(db *database.DB) *PlanDetailsHandler {
	return &PlanDetailsHandler{db: db}
}

func (h *PlanDetailsHandler) AddPlanDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")
	planID := c.Query("plan_id")
	if planID == "" {
		c.JSON(400, gin.H{"error": "Plan ID is required"})
		return
	}

	// Verify plan belongs to user
	var exists bool
	err := h.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM plans WHERE id = $1 AND user_id = $2)", planID, userID)
	if err != nil || !exists {
		c.JSON(404, gin.H{"error": "Plan not found"})
		return
	}

	var detail models.PlanDetails
	if err := c.ShouldBindJSON(&detail); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	_, err = h.db.Exec(
		"INSERT INTO plan_details (plan_id, move, num_of_reps) VALUES ($1, $2, $3)",
		planID, detail.Move, detail.NumOfReps,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error adding plan detail"})
		return
	}

	c.Status(201)
}

func (h *PlanDetailsHandler) GetPlanDetails(c *gin.Context) {
	userID := c.GetInt64("user_id")
	planID := c.Query("plan_id")
	if planID == "" {
		c.JSON(400, gin.H{"error": "Plan ID is required"})
		return
	}

	// Verify plan belongs to user
	var exists bool
	err := h.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM plans WHERE id = $1 AND user_id = $2)", planID, userID)
	if err != nil || !exists {
		c.JSON(404, gin.H{"error": "Plan not found"})
		return
	}

	var details []models.PlanDetails
	err = h.db.Select(&details, "SELECT * FROM plan_details WHERE plan_id = $1", planID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching plan details"})
		return
	}

	c.JSON(200, details)
}

func (h *PlanDetailsHandler) UpdatePlanDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")
	planID := c.Query("plan_id")
	if planID == "" {
		c.JSON(400, gin.H{"error": "Plan ID is required"})
		return
	}

	// Verify plan belongs to user
	var exists bool
	err := h.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM plans WHERE id = $1 AND user_id = $2)", planID, userID)
	if err != nil || !exists {
		c.JSON(404, gin.H{"error": "Plan not found"})
		return
	}

	var detail models.PlanDetails
	if err := c.ShouldBindJSON(&detail); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.db.Exec(
		"UPDATE plan_details SET num_of_reps = $1 WHERE plan_id = $2 AND move = $3",
		detail.NumOfReps, planID, detail.Move,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating plan detail"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error checking update"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Plan detail not found"})
		return
	}

	c.Status(200)
}

func (h *PlanDetailsHandler) DeletePlanDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")
	planID := c.Query("plan_id")
	if planID == "" {
		c.JSON(400, gin.H{"error": "Plan ID is required"})
		return
	}

	move := c.Query("move")
	if move == "" {
		c.JSON(400, gin.H{"error": "Move parameter is required"})
		return
	}

	// Verify plan belongs to user
	var exists bool
	err := h.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM plans WHERE id = $1 AND user_id = $2)", planID, userID)
	if err != nil || !exists {
		c.JSON(404, gin.H{"error": "Plan not found"})
		return
	}

	result, err := h.db.Exec(
		"DELETE FROM plan_details WHERE plan_id = $1 AND move = $2",
		planID, move,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting plan detail"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error checking deletion"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Plan detail not found"})
		return
	}

	c.Status(200)
} 