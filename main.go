package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dstotijn/go-notion"
)

var (
	APIKey     string
	DatabaseID string
)

var recordCount int

func init() {
	log.SetFlags(0)
	flag.IntVar(&recordCount, "record-count", 10, "how many records to display")

	APIKey = os.Getenv("NOTION_API_KEY")
	DatabaseID = os.Getenv("NOTION_DATABASE_ID")
}

func main() {
	flag.Parse()
	client := notion.NewClient(APIKey)

	filter := notion.DatabaseQuery{Filter: &notion.DatabaseQueryFilter{
		Property: "Koordynaty",
		Text:     &notion.TextDatabaseQueryFilter{IsNotEmpty: true},
	}}

	res, err := client.QueryDatabase(context.Background(), DatabaseID, &filter)
	if err != nil {
		log.Fatalln("failed to query Notion database:", err)
	}

	for i := 0; i < 1; i++ {
		page := res.Results[i]

		switch props := page.Properties.(type) {
		case notion.DatabasePageProperties:

			placeRecord := PlaceRecord{
				Title:   props["Nazwa"].Title[0].PlainText,
				Type:    props["Typ"].Select.Name,
				Towns:   parseSelectOptions(props["Miejscowość"].MultiSelect),
				Section: props["Sekcja"].Select.Name,
				Region:  props["Region"].RichText[0].Text.Content,
				Tags:    parseSelectOptions(props["Tagi"].MultiSelect),
				Coords:  props["Koordynaty"].RichText[0].PlainText,
			}

			fmt.Printf("%+v\n", placeRecord)

			blk, err := client.FindBlockChildrenByID(context.Background(), page.ID, &notion.PaginationQuery{PageSize: 100})
			if err != nil {
				log.Fatalln("failed to find block children by id:", err)
			}

			fmt.Println("- - - blocks - - -")
			for _, v := range blk.Results {
				if v.Heading1 != nil && len(v.Heading1.Text) > 0 {
					fmt.Printf("heading: %s\n", v.Heading1.Text[0].PlainText)
				}

				if v.Paragraph != nil && len(v.Paragraph.Text) > 0 {
					fmt.Printf("paragraph: %v\n", v.Paragraph.Text[0].PlainText)
				}

				if v.Image != nil {
					fmt.Printf("image url: %s \n", v.Image.File.URL)
				}

				if v.Column != nil {
					log.Printf("column %v: \n", v.Column)
				}

				cllist, err := client.FindBlockChildrenByID(context.Background(), v.ID, &notion.PaginationQuery{PageSize: 100})
				if err != nil {
					log.Fatalln("failed to find block children for column_list by id:", err)
				}

				for j, blk2 := range cllist.Results {
					if blk2.Column != nil {
						log.Println("column", j, "len", len(blk2.Column.Children))
						for i, w := range blk2.Column.Children {
							log.Println("xd", i)
							if w.Image != nil {
								fmt.Printf("image url %s: \n", w.Image.File.URL)
							}
						}
					}
				}

				if v.ColumnList != nil {
					log.Println("column_list", len(v.ColumnList.Children))
					for i, w := range v.ColumnList.Children {
						log.Println("xd", i)
						if w.Image != nil {
							fmt.Printf("image url %s: \n", w.Image.File.URL)
						}
					}
				}
			}
		}
	}

	if res.HasMore {
		fmt.Println("exporter: there are more records")
	}
}

func parseSelectOptions(options []notion.SelectOptions) []string {
	towns := make([]string, 0)
	for _, v := range options {
		towns = append(towns, v.Name)
	}

	return towns
}
