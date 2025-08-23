package controllers

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// CycleInsight contains prediction data for user's cycle
type CycleInsight struct {
	AverageLength      float64   `json:"average_length"`
	NextPeriodStart    time.Time `json:"next_period_start"`
	PredictedOvulation time.Time `json:"predicted_ovulation"`
	FertileWindowStart time.Time `json:"fertile_window_start"`
	FertileWindowEnd   time.Time `json:"fertile_window_end"`
	IsIrregular        bool      `json:"is_irregular"`
	CommonMood         string    `json:"common_mood,omitempty"`
	CommonSymptoms     []string  `json:"common_symptoms,omitempty"`
	TrackedCycleCount  int       `json:"tracked_cycle_count"`
}

func GetCycleInsights(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var cycles []models.Cycle
	if err := config.DB.Where("user_id = ?", userID).Order("start_date asc").Find(&cycles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cycle data"})
		return
	}

	if len(cycles) < 2 {
		// Return empty insight with default values when not enough data
		insight := CycleInsight{
			AverageLength:      0,
			NextPeriodStart:    time.Now(),
			PredictedOvulation: time.Now(),
			FertileWindowStart: time.Now(),
			FertileWindowEnd:   time.Now(),
			IsIrregular:        false,
			CommonMood:         "",
			CommonSymptoms:     []string{},
			TrackedCycleCount:  len(cycles),
		}
		c.JSON(http.StatusOK, gin.H{
			"insight": insight,
			"message": "Not enough cycle data to calculate insights. Need at least 2 cycles.",
		})
		return
	}

	// Calculate average cycle length
	var lengths []int
	var total int
	for _, cycle := range cycles {
		lengths = append(lengths, cycle.Length)
		total += cycle.Length
	}
	avgLength := float64(total) / float64(len(lengths))

	// Calculate standard deviation to flag irregularity
	variance := 0.0
	for _, l := range lengths {
		diff := float64(l) - avgLength
		variance += diff * diff
	}
	stddev := variance / float64(len(lengths))
	isIrregular := stddev > 4 // You can tune this threshold

	// Determine last cycle start
	lastCycle := cycles[len(cycles)-1]
	nextStart := lastCycle.StartDate.AddDate(0, 0, int(avgLength))
	ovulation := nextStart.AddDate(0, 0, -14)
	fertileStart := ovulation.AddDate(0, 0, -5)
	fertileEnd := ovulation.AddDate(0, 0, 1)

	// Analyze mood/symptom patterns
	moodCounts := map[string]int{}
	symptomCounts := map[string]int{}

	for _, c := range cycles {
		if c.Mood != "" {
			moodCounts[c.Mood]++
		}
		if c.Symptoms != "" {
			symptoms := parseSymptoms(c.Symptoms)
			for _, s := range symptoms {
				symptomCounts[s]++
			}
		}
	}

	commonMood := mostCommon(moodCounts)
	commonSymptoms := topSymptoms(symptomCounts)

	insight := CycleInsight{
		AverageLength:      avgLength,
		NextPeriodStart:    nextStart,
		PredictedOvulation: ovulation,
		FertileWindowStart: fertileStart,
		FertileWindowEnd:   fertileEnd,
		IsIrregular:        isIrregular,
		CommonMood:         commonMood,
		CommonSymptoms:     commonSymptoms,
		TrackedCycleCount:  len(cycles),
	}

	c.JSON(http.StatusOK, insight)
}

func parseSymptoms(s string) []string {
	// Treat as comma-separated string (or later JSON)
	var result []string
	for _, part := range splitAndTrim(s, ",") {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func splitAndTrim(s string, sep string) []string {
	var result []string
	parts := strings.Split(s, sep)
	for _, p := range parts {
		result = append(result, strings.TrimSpace(p))
	}
	return result
}

func mostCommon(m map[string]int) string {
	max := 0
	var key string
	for k, v := range m {
		if v > max {
			key = k
			max = v
		}
	}
	return key
}

func topSymptoms(m map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range m {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	top := []string{}
	for i := 0; i < len(sorted) && i < 3; i++ {
		top = append(top, sorted[i].Key)
	}
	return top
}
