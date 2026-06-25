package main

import (
	"context"
	"fmt"

	"github.com/amayakmt/blog-aggregator/internal/database"
)

// command holds the parsed CLI input: the verb (e.g. "login") and any
// positional arguments that follow it.
type command struct {
	Name      string
	Arguments []string
}

// commands is a registry that maps command names to their handler functions.
type commands struct {
	RegisteredCommands map[string]func(*state, command) error
}

// run looks up the handler for cmd.Name and executes it.
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.RegisteredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return f(s, cmd)
}

// register adds or replaces the handler for the given command name.
func (c *commands) register(name string, f func(*state, command) error) {
	c.RegisteredCommands[name] = f
}

// middlewareLoggedIn wraps a handler that requires an authenticated user.
// It fetches the current user from the DB and passes it to the inner handler,
// so the handler itself doesn't need to repeat that boilerplate.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("could not get current user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
