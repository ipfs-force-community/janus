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

	end := time.Now().Unix()
	start := time.Now().AddDate(0, 0, -days).Unix()

	var results []DailyMinerStat
	if err := s.db.Model(&orm.Miner{}).
		Select("DATE_FORMAT(FROM_UNIXTIME(timestamp), '%Y-%m-%d') AS date, COUNT(*) AS count").
		Where("timestamp BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if results == nil {
		c.JSON(http.StatusOK, []DailyMinerStat{})
		return
	}

	c.JSON(http.StatusOK, results)
}
