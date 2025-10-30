package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"blog-aggregator/internal/database"
	"blog-aggregator/internal/state"
	"blog-aggregator/rss"

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
	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Println(item)
	}

	return nil
}

func HandlerAddFeed(s *state.State, cmd Command) error {
	currentLoggedInUser := s.Config.CurrentUserName

	if len(cmd.Args) < 2 {
		return errors.New("name and feed URL are required")
	}

	if currentLoggedInUser == "" {
		return errors.New("no user is currently logged in")
	}

	user, err := s.DB.GetUserByName(context.Background(), currentLoggedInUser)
	if err != nil {
		return fmt.Errorf("failed to get user %s: %w", currentLoggedInUser, err)
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	_, err = s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		UserID:    user.ID,
		Url:       feedURL,
	})

	if err != nil {
		return fmt.Errorf("failed to add feed %s for user %s: %w", feedURL, currentLoggedInUser, err)
	}

	fmt.Printf("feed %s added successfully!\n", feedName)

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