package main

// PlaceRecord represents a single Notion database row.
type PlaceRecord struct {
	Title   string
	Type    string
	Tags    []string
	Town    string
	Section string
	Region  string
	Coords  string
}
