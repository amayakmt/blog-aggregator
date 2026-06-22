package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/amayakmt/blog-aggregator/internal/database"
	"github.com/google/uuid"
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

// handlerLogin switches the active user in config without creating an account.
// It verifies the user exists in the DB first so we never set a dangling name.
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("arguments cannot be empty")
	}

	userName := cmd.Arguments[0]
	// GetUser returns sql.ErrNoRows if the name is not found, which becomes an
	// error here — login is rejected for users who haven't registered yet.
	_, err := s.DB.GetUser(context.Background(), userName)
	if err != nil {
		return err
	}

	err = s.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("The user has been set to %v\n", userName)
	return nil
}

// handlerRegister creates a new user and immediately logs them in.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("arguments cannot be empty")
	}

	userName := cmd.Arguments[0]

	// Probe for an existing user. A nil error means the name is taken; we only
	// proceed when GetUser returns sql.ErrNoRows (user does not exist yet).
	_, err := s.DB.GetUser(context.Background(), userName)
	if err == nil {
		return fmt.Errorf("this username already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	now := time.Now()

	params := database.CreateUserParams{
		ID:        uuid.New(), // generate a fresh UUID client-side
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	_, err = s.DB.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	// Log the new user in immediately after registration.
	err = s.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User has been registered: %v\n", userName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.DB.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database reset successfully")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for i := 0; i < len(users); i++ {
		if s.Config.CurrentUserName == users[i] {
			fmt.Printf("* %v (current)\n", users[i])
		} else {
			fmt.Printf("* %v\n", users[i])
		}

	}
	return nil
}
