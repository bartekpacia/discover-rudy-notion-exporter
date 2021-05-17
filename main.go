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
	log.SetFlags(0)

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("exporter: failed to load .env file")
	}

	API_KEY = os.Getenv("NOTION_API_KEY")
	DB_ID = os.Getenv("DATABASE_ID")
}

func main() {
	client := notion.NewClient(API_KEY)

	filter := notion.DatabaseQuery{Filter: notion.DatabaseQueryFilter{
		Property: "Koordynaty",
		Text:     notion.TextDatabaseQueryFilter{IsNotEmpty: true},
	}}

	res, err := client.QueryDatabase(context.Background(), DB_ID, filter)
	if err != nil {
		log.Fatalln("exporter: failed to query Notion database:", err)
	}

	fmt.Printf("exporter: got results! hasMore: %t, nextCursor: %p\n", res.HasMore, res.NextCursor)
	for _, page := range res.Results {
		switch props := page.Properties.(type) {
		case notion.DatabasePageProperties:

			var tags []string
			for _, v := range props["Tagi"].MultiSelect {
				tags = append(tags, v.Name)
			}

			placeRecord := PlaceRecord{
				Title:   props["Nazwa"].ID,
				Type:    props["Typ"].Select.Name,
				Town:    props["Miejscowość"].Select.Name,
				Section: props["Sekcja"].Select.Name,
				Region:  props["Region"].Select.Name,
				Tags:    tags,
				Coords:  props["Koordynaty"].RichText[0].PlainText,
			}

			fmt.Printf("%v\n", placeRecord)
		}
	}
}
