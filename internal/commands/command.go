package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"blog-aggregator/internal/database"
	"blog-aggregator/internal/state"
	"blog-aggregator/internal/utils"

	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(s *state.State, c Command) error
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("username argument is required")
	}

	username := cmd.Args[0]

	_, err := s.DB.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user %s does not exist", username)
	}

	if err := s.Config.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user %s: %w", username, err)
	}

	fmt.Printf("username set to %s\n", username)

	return nil
}

func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("username argument is required")
	}

	username := cmd.Args[0]
	_, err := s.DB.GetUserByName(context.Background(), username)

	if err == nil {
		return fmt.Errorf("user %s already exists", username)
	} else if err != sql.ErrNoRows {
		return err
	} else {
		_, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		})

		if err != nil {
			return err
		}

		if err := s.Config.SetUser(username); err != nil {
			return fmt.Errorf("failed to set user %s: %w", username, err)
		}

		fmt.Printf("user %s registered and set as current user\n", username)
	}

	return nil
}

func HandlerResetUsers(s *state.State, cmd Command) error {
	err := s.DB.ResetUserTable(context.Background())

	if err != nil {
		return fmt.Errorf("failed to reset users table: %w", err)
	}

	fmt.Println("users table reset successfully")

	return nil
}

func HandlerGetUsers(s *state.State, cmd Command) error{
	users, err := s.DB.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	currentUser := s.Config.CurrentUserName

	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func HandlerAgg(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("duration argument is required (e.g., '1m', '30s')")
	}

	duration := cmd.Args[0]
	timeBetweenRequests, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	fmt.Printf("collecting feeds every %s\nPress Ctrl+C to stop gracefully...\n\n", timeBetweenRequests)
	fmt.Println()

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		err := utils.ScrapeFeeds(context.Background(), s.DB)
		if err != nil {
			fmt.Printf("error fetching feeds: %v\n", err)
			// continue the loop instead of returning
		}
	}

}

func HandlerAddFeed(s *state.State, cmd Command, user database.User) error {

	if len(cmd.Args) < 2 {
		return errors.New("name and feed URL are required")
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		UserID:    user.ID,
		Url:       feedURL,
	})

	if err != nil {
		return fmt.Errorf("failed to add feed %s for user %s: %w", feedURL, user.Name, err)
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to create feed follow for user %s and feed %s: %w", user.Name, feedURL, err)
	}

	fmt.Printf("feed %s added successfully!\n", feedName)

	return nil
}

func HandlerListFeeds(s *state.State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}

	for _, feed := range feeds {
		username := "unknown"
		if feed.UserName.Valid {
			username = feed.UserName.String
		}
		fmt.Println("Feed Name: ", feed.Name)
		fmt.Println("Feed URL: ", feed.Url)
		fmt.Println("Author: ", username)
		fmt.Println("-----")
	}

	return nil
}

func HandlerFollowFeed(s *state.State, cmd Command, user database.User) error {
	// takes a single url arg and creates a new feed follow record for the current user.

	if len(cmd.Args) < 1 {
		return errors.New("feed URL argument is required")
	}

	feedURL := cmd.Args[0]

	feed, err := s.DB.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to get feed by URL %s: %w", feedURL, err)
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow for user %s and feed %s: %w", user.Name, feedURL, err)
	}

	fmt.Printf("user %s is now following feed %s\n", user.Name, feedURL)

	return nil
}

func HandlerListFollowedFeeds(s *state.State, cmd Command, user database.User) error {

	followedFeeds, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get followed feeds for user %s: %w", user.Name, err)
	}

	for _, feed := range followedFeeds {
		fmt.Println(feed.FeedName)
	}

	return nil
}

func HandlerUnfollowFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL argument is required")
	}

	feedURL := cmd.Args[0]

	feed, err := s.DB.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to get feed by URL %s: %w", feedURL, err)
	}

	err = s.DB.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed %s for user %s: %w", feedURL, user.Name, err)
	}

	fmt.Println("feed unfollowed.")
	return nil
}

func HandlerBrowse(s *state.State, cmd Command, user database.User) error {
	limit := 2

	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil || parsedLimit <= 0 {
			return fmt.Errorf("invalid limit value, please use a positive number")
		}
		limit = parsedLimit
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to get posts for user %s: %w", user.Name, err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title.String)
		fmt.Printf("URL: %s\n", post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("Published At: %s\n", post.PublishedAt.Time.Format(time.RFC1123))
		}
		fmt.Printf("Feed: %s\n", post.FeedName)
		fmt.Println("-----")
	}
	return nil
}


func (c *Commands) Run(s *state.State, cmd Command) error {
	// runs a given command with the proivided state IF it exists

	if handler, exists := c.handlers[cmd.Name]; exists {
		return handler(s, cmd)
	}

	return fmt.Errorf("Command %s not found", cmd.Name)
	
}


func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	// registers a new handler function for a command name

	if c.handlers == nil {
		c.handlers = make(map[string]func(s *state.State, c Command) error)
	}

	c.handlers[name] = f
}