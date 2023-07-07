package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/TurniXXD/go-rss/internal/database"
	"github.com/TurniXXD/go-rss/utils"
	"github.com/google/uuid"
)

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := utils.UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}

		if _, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			}); err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to create post: ", err)
			continue
		}
	}

	log.Printf("\nFeed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}

func Start(
	db *database.Queries,
	concurrencyUnits int,
	timeBetweenReq time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrencyUnits, timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)

	// Empty fields so it runs immediately the first time and then it waits for the interval
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrencyUnits),
		)

		if err != nil {
			log.Println("error fetching feeds: ", err)
			// There is no time for this function to stop because, if we put the "return" keyword here it would stop scraping completely
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			// Add one item to the wait group
			wg.Add(1)

			// Now we can scrape number of concurrencyUnits at the same time
			go scrapeFeed(db, wg, feed)
		}
		// Wait for number of concurrencyUnits goroutines to be Done
		wg.Wait()
	}
}
