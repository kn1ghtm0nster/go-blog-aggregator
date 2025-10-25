package commands

import (
	"errors"
	"fmt"

	"blog-aggregator/internal/state"
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
	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("username set to %s\n", username)

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