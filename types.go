package main

/// PlaceRecord represents a single Notion database row.
type PlaceRecord struct {
	Title  string `json:"Title"`
	Coords string `json:"Koordynaty"`
}
