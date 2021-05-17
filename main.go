package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
)

var (
	API_KEY string
	DB_ID   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("exporter: failed to load .env file")
	}

	API_KEY = os.Getenv("NOTION_API_KEY")
	DB_ID = os.Getenv("DATABASE_ID")
}

func main() {
	client := notion.NewClient(API_KEY)

	res, err := client.QueryDatabase(context.Background(), DB_ID, notion.DatabaseQuery{Filter: notion.DatabaseQueryFilter{
		Property: "Koordynaty",
		Text:     notion.TextDatabaseQueryFilter{IsNotEmpty: true},
	}})
	if err != nil {
		log.Fatalln("exporter: failed to query Notion database:", err)
	}

	fmt.Printf("exporter: got results! hasMore: %t, nextCursor: %p\n", res.HasMore, res.NextCursor)
	fmt.Println("exporter: page ids")
	for _, page := range res.Results {
		fmt.Println(page.ID)
	}
}
