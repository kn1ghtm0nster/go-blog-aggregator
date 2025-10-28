package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"blog-aggregator/internal/database"
	"blog-aggregator/internal/state"

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