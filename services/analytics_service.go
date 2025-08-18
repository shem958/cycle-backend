package services

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// ====== PUBLIC TYPES ======

type TimeValue struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

type BloodPressurePoint struct {
	Time      time.Time `json:"time"`
	Systolic  *int      `json:"systolic,omitempty"`
	Diastolic *int      `json:"diastolic,omitempty"`
	Raw       string    `json:"raw"`
}

type CheckupItem struct {
	ID              uuid.UUID `json:"id"`
	Type            string    `json:"type"` // "pregnancy" | "postpartum"
	VisitDate       time.Time `json:"visit_date"`
	Notes           string    `json:"notes,omitempty"`
	AttachmentCount int       `json:"attachment_count"`
}

type CombinedAnalytics struct {
	UserID              uuid.UUID            `json:"user_id"`
	From                *time.Time           `json:"from,omitempty"`
	To                  *time.Time           `json:"to,omitempty"`
	PregnancyCount      int                  `json:"pregnancy_count"`
	PostpartumCount     int                  `json:"postpartum_count"`
	UpcomingNextCheckup *time.Time           `json:"upcoming_next_checkup,omitempty"`
	WeightTrend         []TimeValue          `json:"weight_trend"`
	BloodPressure       []BloodPressurePoint `json:"blood_pressure"`
	Timeline            []CheckupItem        `json:"timeline"`
}

// ====== SIMPLE IN-MEMORY CACHE ======

type cacheEntry struct {
	data      *CombinedAnalytics
	expiresAt time.Time
}

var (
	analyticsCache   = map[string]cacheEntry{}
	analyticsCacheMu sync.RWMutex
	cacheTTL         = 10 * time.Minute // adjust as needed (5–15 min recommended)
)

func cacheKey(userID uuid.UUID, from, to *time.Time) string {
	var f, t string
	if from != nil {
		f = from.UTC().Format(time.RFC3339Nano)
	}
	if to != nil {
		t = to.UTC().Format(time.RFC3339Nano)
	}
	return userID.String() + "|" + f + "|" + t
}

func getFromCache(userID uuid.UUID, from, to *time.Time) (*CombinedAnalytics, bool) {
	analyticsCacheMu.RLock()
	defer analyticsCacheMu.RUnlock()
	k := cacheKey(userID, from, to)
	entry, ok := analyticsCache[k]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.data, true
}

func putInCache(userID uuid.UUID, from, to *time.Time, data *CombinedAnalytics) {
	analyticsCacheMu.Lock()
	defer analyticsCacheMu.Unlock()
	k := cacheKey(userID, from, to)
	analyticsCache[k] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(cacheTTL),
	}
}

// Optional: for invalidation if you add write operations that change analytics sources
func InvalidateAnalyticsCacheForUser(userID uuid.UUID) {
	analyticsCacheMu.Lock()
	defer analyticsCacheMu.Unlock()
	prefix := userID.String() + "|"
	for k := range analyticsCache {
		if strings.HasPrefix(k, prefix) {
			delete(analyticsCache, k)
		}
	}
}

// ====== CORE AGGREGATION ======

var bpRe = regexp.MustCompile(`^\s*(\d{2,3})\s*/\s*(\d{2,3})\s*$`)

func parseBP(raw string) (*int, *int) {
	m := bpRe.FindStringSubmatch(strings.TrimSpace(raw))
	if len(m) != 3 {
		return nil, nil
	}
	systolic, err1 := strconv.Atoi(m[1])
	diastolic, err2 := strconv.Atoi(m[2])
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	return &systolic, &diastolic
}

func firstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	return ""
}

// GetCombinedAnalytics aggregates pregnancy + postpartum checkups for a user within optional date range.
// Uses an in-memory cache to reduce DB load.
func GetCombinedAnalytics(userID uuid.UUID, from, to *time.Time) (*CombinedAnalytics, error) {
	// Try cache
	if data, ok := getFromCache(userID, from, to); ok {
		return data, nil
	}

	var preg []models.PregnancyCheckup
	var post []models.PostpartumCheckup

	qp := config.DB.Preload("Attachments").Where("user_id = ?", userID)
	qpp := config.DB.Preload("Attachments").Where("user_id = ?", userID)

	if from != nil {
		qp = qp.Where("visit_date >= ?", *from)
		qpp = qpp.Where("visit_date >= ?", *from)
	}
	if to != nil {
		qp = qp.Where("visit_date <= ?", *to)
		qpp = qpp.Where("visit_date <= ?", *to)
	}

	if err := qp.Find(&preg).Error; err != nil {
		return nil, err
	}
	if err := qpp.Find(&post).Error; err != nil {
		return nil, err
	}

	analytics := &CombinedAnalytics{
		UserID:          userID,
		From:            from,
		To:              to,
		PregnancyCount:  len(preg),
		PostpartumCount: len(post),
	}

	// Build trends and unified timeline
	var weights []TimeValue
	var bps []BloodPressurePoint
	var timeline []CheckupItem

	for _, c := range preg {
		if c.Weight > 0 {
			weights = append(weights, TimeValue{Time: c.VisitDate, Value: c.Weight})
		}
		if strings.TrimSpace(c.BloodPressure) != "" {
			sys, dia := parseBP(c.BloodPressure)
			bps = append(bps, BloodPressurePoint{
				Time: c.VisitDate, Systolic: sys, Diastolic: dia, Raw: c.BloodPressure,
			})
		}
		timeline = append(timeline, CheckupItem{
			ID: c.ID, Type: "pregnancy", VisitDate: c.VisitDate,
			Notes:           strings.TrimSpace(c.DoctorNotes),
			AttachmentCount: len(c.Attachments),
		})
	}

	for _, c := range post {
		timeline = append(timeline, CheckupItem{
			ID: c.ID, Type: "postpartum", VisitDate: c.VisitDate,
			Notes:           strings.TrimSpace(firstNonEmpty(c.MotherHealthNotes, c.BabyHealthNotes, c.Complications, c.MentalHealth)),
			AttachmentCount: len(c.Attachments),
		})
	}

	sort.Slice(weights, func(i, j int) bool { return weights[i].Time.Before(weights[j].Time) })
	sort.Slice(bps, func(i, j int) bool { return bps[i].Time.Before(bps[j].Time) })
	sort.Slice(timeline, func(i, j int) bool { return timeline[i].VisitDate.Before(timeline[j].VisitDate) })

	analytics.WeightTrend = weights
	analytics.BloodPressure = bps
	analytics.Timeline = timeline

	// Upcoming next checkup
	now := time.Now()
	var candidates []time.Time
	for _, c := range preg {
		if !c.NextCheckupAt.IsZero() && c.NextCheckupAt.After(now) {
			candidates = append(candidates, c.NextCheckupAt)
		}
	}
	for _, c := range post {
		if !c.NextCheckupAt.IsZero() && c.NextCheckupAt.After(now) {
			candidates = append(candidates, c.NextCheckupAt)
		}
	}
	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Before(candidates[j]) })
		analytics.UpcomingNextCheckup = &candidates[0]
	}

	// Cache result
	putInCache(userID, from, to, analytics)
	return analytics, nil
}
