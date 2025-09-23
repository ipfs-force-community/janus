package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ipfs-force-community/janus/database/orm"
)

// DailyMinerStat represents the daily statistics of new miners
type DailyMinerStat struct {
	Date  string  `json:"date"`
	Count int64   `json:"count"`
	Cost  float64 `json:"cost"`
}

// GetDailyMinerStats handles the GET /miners endpoint to retrieve daily new miner statistics
func (s *Server) GetDailyMinerStats(c *gin.Context) {
	intervalParam := c.Query("interval")
	if intervalParam == "" {
		intervalParam = "7d"
	}

	days := 7
	if strings.HasSuffix(intervalParam, "d") {
		if n, err := strconv.Atoi(strings.TrimSuffix(intervalParam, "d")); err == nil {
			days = n
		}
	}

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	end := endTime.Unix()
	start := startTime.Unix()

	var dbResults []DailyMinerStat
	if err := s.db.Model(&orm.Miner{}).
		Select(`
			DATE_FORMAT(FROM_UNIXTIME(timestamp), '%Y-%m-%d') AS date, 
			COUNT(*) AS count,
			AVG(CAST(cost AS DECIMAL(32,0))) / 1e18 AS cost
		`).
		Where("timestamp BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Scan(&dbResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	countMap := make(map[string]int64, len(dbResults))
	costMap := make(map[string]float64, len(dbResults))
	for _, r := range dbResults {
		countMap[r.Date] = r.Count
		costMap[r.Date] = r.Cost
	}

	var results []DailyMinerStat
	var lastNonZeroCost float64

	for d := 0; d <= days; d++ {
		date := startTime.AddDate(0, 0, d).Format("2006-01-02")
		count := countMap[date]
		cost := costMap[date]

		if lastNonZeroCost == 0 && cost != 0 {
			lastNonZeroCost = cost
		}

		results = append(results, DailyMinerStat{
			Date:  date,
			Count: count,
			Cost:  cost,
		})
	}

	for i := 0; i < len(results); i++ {
		if results[i].Cost == 0 {
			results[i].Cost = lastNonZeroCost
		}
	}

	if results == nil {
		c.JSON(http.StatusOK, []DailyMinerStat{})
		return
	}

	c.JSON(http.StatusOK, results)
}
