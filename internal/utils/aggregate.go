package utils

import (
	"context"
	"fmt"

	"blog-aggregator/internal/database"
	"blog-aggregator/rss"
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

	fmt.Printf("found %d posts in feed %q:\n", len(rssFeed.Channel.Item), feed.Name)
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("- %s\n", item.Title)
	}

	return nil
}