package handlers

import (
	"daily-tracker/database"
	"daily-tracker/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateActivity creates a new activity entry
func CreateActivity(c *gin.Context) {
	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// GetActivities retrieves all activities
func GetActivities(c *gin.Context) {
	var activities []models.Activity
	if err := database.DB.Order("date desc").Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// GetActivityByDate retrieves activity for a specific date
func GetActivityByDate(c *gin.Context) {
	date := c.Param("date")
	var activity models.Activity

	if err := database.DB.Where("date = ?", date).First(&activity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// UpdateActivity updates an existing activity
func UpdateActivity(c *gin.Context) {
	id := c.Param("id")
	var activity models.Activity

	if err := database.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// DeleteActivity deletes an activity
func DeleteActivity(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Activity{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

// GetWeeklyAnalysis provides weekly analysis
func GetWeeklyAnalysis(c *gin.Context) {
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	var activities []models.Activity

	if err := database.DB.Where("date >= ?", weekAgo).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	analysis := calculateAnalysis(activities)
	c.JSON(http.StatusOK, gin.H{
		"period":   "weekly",
		"analysis": analysis,
		"data":     activities,
	})
}

// GetMonthlyAnalysis provides monthly analysis
func GetMonthlyAnalysis(c *gin.Context) {
	monthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	var activities []models.Activity

	if err := database.DB.Where("date >= ?", monthAgo).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	analysis := calculateAnalysis(activities)
	c.JSON(http.StatusOK, gin.H{
		"period":   "monthly",
		"analysis": analysis,
		"data":     activities,
	})
}

// GetYearlyAnalysis provides yearly analysis
func GetYearlyAnalysis(c *gin.Context) {
	yearAgo := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	var activities []models.Activity

	if err := database.DB.Where("date >= ?", yearAgo).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	analysis := calculateAnalysis(activities)
	c.JSON(http.StatusOK, gin.H{
		"period":   "yearly",
		"analysis": analysis,
		"data":     activities,
	})
}

// Helper function to calculate analysis
func calculateAnalysis(activities []models.Activity) map[string]interface{} {
	if len(activities) == 0 {
		return map[string]interface{}{
			"total_days":         0,
			"avg_learning_hours": 0,
			"avg_sleep_hours":    0,
			"avg_office_hours":   0,
			"avg_wasted_hours":   0,
			"total_learning":     0,
			"total_sleep":        0,
			"total_office":       0,
			"total_wasted":       0,
		}
	}

	var totalLearning, totalSleep, totalOffice, totalWasted float64
	const maxSleepHours = 9.0
	const hoursInDay = 24.0

	for _, activity := range activities {
		totalLearning += activity.LearningHours

		// Cap sleep hours at 9
		sleepHours := activity.SleepHours
		if sleepHours > maxSleepHours {
			sleepHours = maxSleepHours
		}
		totalSleep += sleepHours

		totalOffice += activity.OfficeHours

		// Calculate wasted hours: 24 - (learning + sleep(max 9) + office)
		usedHours := activity.LearningHours + sleepHours + activity.OfficeHours
		wastedHours := hoursInDay - usedHours
		if wastedHours < 0 {
			wastedHours = 0
		}
		totalWasted += wastedHours
	}

	days := float64(len(activities))
	return map[string]interface{}{
		"total_days":         len(activities),
		"avg_learning_hours": totalLearning / days,
		"avg_sleep_hours":    totalSleep / days,
		"avg_office_hours":   totalOffice / days,
		"avg_wasted_hours":   totalWasted / days,
		"total_learning":     totalLearning,
		"total_sleep":        totalSleep,
		"total_office":       totalOffice,
		"total_wasted":       totalWasted,
	}
}

// Made with Bob
