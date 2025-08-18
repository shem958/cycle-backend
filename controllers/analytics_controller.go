package controllers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
	"github.com/shem958/cycle-backend/services"
)

// --------- Helpers ---------

// tries to read role from context (set by your AuthMiddleware) otherwise loads from DB
func getRoleFromContextOrDB(c *gin.Context) (string, error) {
	if v, ok := c.Get("role"); ok {
		if role, ok := v.(string); ok && role != "" {
			return role, nil
		}
	}
	// fallback: look up by user_id in context
	var uid uuid.UUID
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				uid = id
			}
		}
	}
	if uid == uuid.Nil {
		return "", nil
	}
	var u models.User
	if err := config.DB.First(&u, "id = ?", uid).Error; err != nil {
		return "", err
	}
	return u.Role, nil
}

func parseRange(c *gin.Context) (*time.Time, *time.Time, bool) {
	var fromPtr, toPtr *time.Time
	if from := c.Query("from"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "from must be RFC3339 (e.g., 2025-01-01T00:00:00Z)"})
			return nil, nil, false
		}
		fromPtr = &t
	}
	if to := c.Query("to"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to must be RFC3339 (e.g., 2025-12-31T23:59:59Z)"})
			return nil, nil, false
		}
		toPtr = &t
	}
	return fromPtr, toPtr, true
}

// --------- User self-analytics ---------
// GET /analytics/user/:user_id/pregnancy-postpartum
func GetPregnancyPostpartumAnalytics(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	fromPtr, toPtr, ok := parseRange(c)
	if !ok {
		return
	}

	result, err := services.GetCombinedAnalytics(userID, fromPtr, toPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load analytics"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// --------- Doctor-only: view patient analytics ---------
// GET /analytics/doctor/patient/:patient_id/pregnancy-postpartum
func GetPatientAnalyticsForDoctor(c *gin.Context) {
	role, err := getRoleFromContextOrDB(c)
	if err != nil || (role != models.RoleDoctor && role != models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "doctor or admin role required"})
		return
	}

	patientID, err := uuid.Parse(c.Param("patient_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	fromPtr, toPtr, ok := parseRange(c)
	if !ok {
		return
	}

	result, err := services.GetCombinedAnalytics(patientID, fromPtr, toPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load analytics"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// --------- CSV export (user self) ---------
// GET /analytics/user/:user_id/pregnancy-postpartum.csv
func ExportPregnancyPostpartumCSV(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	fromPtr, toPtr, ok := parseRange(c)
	if !ok {
		return
	}

	result, err := services.GetCombinedAnalytics(userID, fromPtr, toPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load analytics"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=analytics_pregnancy_postpartum.csv")
	c.Header("Content-Type", "text/csv")
	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// Section 1: Weight trend
	_ = w.Write([]string{"Weight Trend"})
	_ = w.Write([]string{"Time", "Weight"})
	for _, p := range result.WeightTrend {
		_ = w.Write([]string{p.Time.UTC().Format(time.RFC3339), strconvFloat(p.Value)})
	}
	_ = w.Write([]string{}) // blank line

	// Section 2: Blood pressure
	_ = w.Write([]string{"Blood Pressure"})
	_ = w.Write([]string{"Time", "Systolic", "Diastolic", "Raw"})
	for _, bp := range result.BloodPressure {
		var s, d string
		if bp.Systolic != nil {
			s = strconvInt(*bp.Systolic)
		}
		if bp.Diastolic != nil {
			d = strconvInt(*bp.Diastolic)
		}
		_ = w.Write([]string{bp.Time.UTC().Format(time.RFC3339), s, d, bp.Raw})
	}
	_ = w.Write([]string{})

	// Section 3: Timeline
	_ = w.Write([]string{"Timeline"})
	_ = w.Write([]string{"Type", "VisitDate", "Notes", "AttachmentCount"})
	for _, t := range result.Timeline {
		_ = w.Write([]string{t.Type, t.VisitDate.UTC().Format(time.RFC3339), t.Notes, strconvInt(t.AttachmentCount)})
	}
}

// --------- CSV export (doctor-only) ---------
// GET /analytics/doctor/patient/:patient_id/pregnancy-postpartum.csv
func ExportPatientAnalyticsCSVForDoctor(c *gin.Context) {
	role, err := getRoleFromContextOrDB(c)
	if err != nil || (role != models.RoleDoctor && role != models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "doctor or admin role required"})
		return
	}

	patientID, err := uuid.Parse(c.Param("patient_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	fromPtr, toPtr, ok := parseRange(c)
	if !ok {
		return
	}

	result, err := services.GetCombinedAnalytics(patientID, fromPtr, toPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load analytics"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=patient_analytics_pregnancy_postpartum.csv")
	c.Header("Content-Type", "text/csv")
	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	_ = w.Write([]string{"Weight Trend"})
	_ = w.Write([]string{"Time", "Weight"})
	for _, p := range result.WeightTrend {
		_ = w.Write([]string{p.Time.UTC().Format(time.RFC3339), strconvFloat(p.Value)})
	}
	_ = w.Write([]string{})

	_ = w.Write([]string{"Blood Pressure"})
	_ = w.Write([]string{"Time", "Systolic", "Diastolic", "Raw"})
	for _, bp := range result.BloodPressure {
		var s, d string
		if bp.Systolic != nil {
			s = strconvInt(*bp.Systolic)
		}
		if bp.Diastolic != nil {
			d = strconvInt(*bp.Diastolic)
		}
		_ = w.Write([]string{bp.Time.UTC().Format(time.RFC3339), s, d, bp.Raw})
	}
	_ = w.Write([]string{})

	_ = w.Write([]string{"Timeline"})
	_ = w.Write([]string{"Type", "VisitDate", "Notes", "AttachmentCount"})
	for _, t := range result.Timeline {
		_ = w.Write([]string{t.Type, t.VisitDate.UTC().Format(time.RFC3339), t.Notes, strconvInt(t.AttachmentCount)})
	}
}

// small helpers for CSV
func strconvFloat(f float64) string { return fmtFloat(f) }
func strconvInt(i int) string       { return fmtInt(i) }

// local formatting (avoid importing strconv repeatedly in loops)
func fmtFloat(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
func fmtInt(v int) string       { return strconv.Itoa(v) }
