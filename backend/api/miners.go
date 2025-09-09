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
	Date  string `json:"date"`
	Count int64  `json:"count"`
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
		Select("DATE_FORMAT(FROM_UNIXTIME(timestamp), '%Y-%m-%d') AS date, COUNT(*) AS count").
		Where("timestamp BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Scan(&dbResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resultMap := make(map[string]int64, len(dbResults))
	for _, r := range dbResults {
		resultMap[r.Date] = r.Count
	}

	var results []DailyMinerStat
	for d := 0; d <= days; d++ {
		date := startTime.AddDate(0, 0, d).Format("2006-01-02")
		count := resultMap[date]
		results = append(results, DailyMinerStat{
			Date:  date,
			Count: count,
		})
	}

	if results == nil {
		c.JSON(http.StatusOK, []DailyMinerStat{})
		return
	}

	c.JSON(http.StatusOK, results)
}
