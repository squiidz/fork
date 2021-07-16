package main

import (
	"time"
)

type Link struct {
	URL         string `json:"url"`
	Short       string `json:"short"`
	Click       int64  `json:"click"`
	UpdateCount int64  `json:"updateCount"`
	CreatedAt   int64  `json:"createdAt"`
	LastViewed  int64  `json:"lastViewed"`
	LastUpdated int64  `json:"lastUpdated"`
}

func NewLink(url string, count int64) (*Link, error) {
	enc, err := EncodeURL(url, count)
	if err != nil {
		return nil, err
	}
	return &Link{
		URL:         url,
		Short:       enc,
		Click:       0,
		UpdateCount: 0,
		CreatedAt:   time.Now().Unix(),
		LastViewed:  time.Now().Unix(),
		LastUpdated: time.Now().Unix(),
	}, nil
}
