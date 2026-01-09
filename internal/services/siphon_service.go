package services

import (
	"net/http"
	"time"
)

type SiphonService struct{}

func NewSiphonService() *SiphonService {
	return &SiphonService{}
}

// CheckHTTP performs a GET request and returns status code and latency in milliseconds
func (s *SiphonService) CheckHTTP(url string) (int, int64, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	latency := time.Since(start).Milliseconds()
	return resp.StatusCode, latency, nil
}
