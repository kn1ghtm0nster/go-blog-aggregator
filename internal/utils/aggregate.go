package utils

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"blog-aggregator/internal/database"
	"blog-aggregator/rss"

	"github.com/google/uuid"
)


func ScrapeFeeds(ctx context.Context, db *database.Queries) error {

	feed, err := db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	err = db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("failed to mark feed %s as fetched: %w", feed.ID, err)
	}

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed from url %s: %w", feed.Url, err)
	}

	
	for _, item := range rssFeed.Channel.Item {
		var publishedAt sql.NullTime
		if item.PubDate != "" {
			parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				parsedTime, err = time.Parse(time.RFC822Z, item.PubDate)
			}
			if err == nil {
				publishedAt = sql.NullTime{Time: parsedTime, Valid: true}
			} 
		}

		newPost := database.CreatePostParams{
			ID:          uuid.New().String(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:      sql.NullString{String: item.Title, Valid: item.Title != ""},
			Url:        item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID:     feed.ID,
		}

		_, err := db.CreatePost(ctx, newPost)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				continue
			}

			fmt.Printf("Error saving post %s: %v\n", item.Link, err)
			continue
		}

	}

	return nil
}