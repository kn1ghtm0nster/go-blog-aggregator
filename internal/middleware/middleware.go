package middleware

import (
	"context"
	"errors"
	"fmt"

	"blog-aggregator/internal/commands"
	"blog-aggregator/internal/database"
	"blog-aggregator/internal/state"
)

func MiddlewareLoggedIn(handler func(s *state.State, cmd commands.Command, user database.User) error) func(*state.State, commands.Command) error {
	// logic goes here

	return func(s *state.State, cmd commands.Command) error {
		currentUser := s.Config.CurrentUserName
		if currentUser == "" {
			return errors.New("no user is currently logged in")
		}

		user, err := s.DB.GetUserByName(context.Background(), currentUser)
		if err != nil {
			return fmt.Errorf("failed to get user %s: %w", currentUser, err)
		}

		return handler(s, cmd, user)
	}
}