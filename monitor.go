package goat

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

//Monit struct
type Monit struct {
	mu                  sync.RWMutex
	UpTime              time.Time
	ResponseCounts      map[string]int
	TotalResponseCounts map[string]int
	TotalResponseTime   time.Time
	Pid                 int
}

//MonitData struct
type MonitData struct {
	Pid                    int
	UpTime                 string
	UpTimeSec              float64
	Time                   string
	TimeUnix               int64
	StatusCodeCount        map[string]int
	TotalStatusCodeCount   map[string]int
	Count                  int
	TotalCount             int
	TotalResponseTime      string
	TotalResponseTimeSec   float64
	AverageResponseTime    string
	AverageResponseTimeSec float64
	Memory                 string
}

//Get func to get the Monit Data
func (m *Monit) Get() *MonitData {
	m.mu.RLock()
	responseCounts := make(map[string]int, len(m.ResponseCounts))
	totalResponseCounts := make(map[string]int, len(m.TotalResponseCounts))

	upTime := time.Since(m.UpTime)
	totalCount := 0
	count := 0
	for code, count := range m.TotalResponseCounts {
		totalResponseCounts[code] = count
		totalCount += count
	}
	for code, current := range m.ResponseCounts {
		responseCounts[code] = current
		count += current
	}

	totalResponseTime := m.TotalResponseTime.Sub(time.Time{})
	averageResponseTime := time.Duration(0)

	if totalCount > 0 {
		avg := int64(totalResponseTime) / (int64)(totalCount)
		averageResponseTime = time.Duration(avg)
	}
	m.mu.RUnlock()

	data := &MonitData{
		Pid:                    m.Pid,
		UpTime:                 upTime.String(),
		UpTimeSec:              upTime.Seconds(),
		Time:                   time.Now().String(),
		TimeUnix:               time.Now().Unix(),
		StatusCodeCount:        responseCounts,
		TotalStatusCodeCount:   totalResponseCounts,
		TotalResponseTime:      totalResponseTime.String(),
		TotalResponseTimeSec:   totalResponseTime.Seconds(),
		AverageResponseTimeSec: averageResponseTime.Seconds(),
		AverageResponseTime:    averageResponseTime.String(),
	}

	return data
}

//ResetResponseCounts every second
func (m *Monit) ResetResponseCounts() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponseCounts = map[string]int{}
}

//NewMonitor to get new monit object
func NewMonitor() *Monit {
	monit := &Monit{
		UpTime:              time.Now(),
		Pid:                 os.Getpid(),
		ResponseCounts:      map[string]int{},
		TotalResponseCounts: map[string]int{},
		TotalResponseTime:   time.Time{},
	}

	go func() {
		monit.ResetResponseCounts()

		time.Sleep(time.Second)
	}()

	return monit
}

//Monitor middleware to update the monit data
func (m *Monit) Monitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		nrw := NewResponseWriter(w)
		next.ServeHTTP(w, r)
		responseTime := time.Since(start)
		m.mu.Lock()
		defer m.mu.Unlock()
		statusCode := fmt.Sprintf("%d", nrw.Status())
		m.ResponseCounts[statusCode]++
		m.TotalResponseCounts[statusCode]++
		m.TotalResponseTime = m.TotalResponseTime.Add(responseTime)
	})
}
