package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
)

var (
	ApiKey     string
	DatabaseID string
)

var (
	recordCount int
)

func init() {
	log.SetFlags(0)
	flag.IntVar(&recordCount, "record-count", 10, "how many records to display")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("exporter: failed to load .env file")
	}

	ApiKey = os.Getenv("NOTION_API_KEY")
	DatabaseID = os.Getenv("DATABASE_ID")
}

func main() {
	flag.Parse()
	client := notion.NewClient(ApiKey)

	filter := notion.DatabaseQuery{Filter: notion.DatabaseQueryFilter{
		Property: "Koordynaty",
		Text:     notion.TextDatabaseQueryFilter{IsNotEmpty: true},
	}}

	res, err := client.QueryDatabase(context.Background(), DatabaseID, filter)
	if err != nil {
		log.Fatalln("exporter: failed to query Notion database:", err)
	}

	for i := 0; i < 3; i++ {
		page := res.Results[i]

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

	if res.HasMore {
		fmt.Println("exporter: there are more records")
	}
}
