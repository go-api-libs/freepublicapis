// Package freepublicapis provides types for the API.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.
package freepublicapis

import "net/url"

type ListApisParams struct {
	Limit int
	Sort  Sort
}

type Sort string

// SimpleAPIInfo defines a model
type SimpleAPIInfo struct {
	ID            int     `json:"id"`
	Emoji         string  `json:"emoji"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Documentation url.URL `json:"documentation"`
	Methods       int     `json:"methods"`
	Health        int     `json:"health"`
	Source        url.URL `json:"source"`
}

// APIInfos defines a model
type APIInfos []APIInfo

// APIInfo defines a model
type APIInfo struct {
	ID             int     `json:"id"`
	Emoji          string  `json:"emoji"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Documentation  url.URL `json:"documentation"`
	Methods        int     `json:"methods"`
	Health         int     `json:"health"`
	Popularity     int     `json:"popularity"`
	AvgReliability int     `json:"avg_reliability"`
	AvgError       int     `json:"avg_error"`
	AvgLatency     int     `json:"avg_latency"`
	Source         url.URL `json:"source"`
}
