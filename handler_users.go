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

	if err = s.Config.SetUser(userName); err != nil {
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

	if _, err = s.DB.CreateUser(context.Background(), params); err != nil {
		return err
	}

	// Log the new user in immediately after registration.
	if err = s.Config.SetUser(userName); err != nil {
		return err
	}

	fmt.Printf("User has been registered: %v\n", userName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.DB.ResetUsers(context.Background()); err != nil {
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

	for _, user := range users {
		if s.Config.CurrentUserName == user {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Printf("* %v\n", user)
		}
	}
	return nil
}
