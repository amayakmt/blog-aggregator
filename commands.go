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

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	RegisteredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	f, ok := c.RegisteredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command")
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.RegisteredCommands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("arguments cannot be empty")
	}

	userName := cmd.Arguments[0]
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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("arguments cannot be empty")
	}

	userName := cmd.Arguments[0]

	_, err := s.DB.GetUser(context.Background(), userName)
	if err == nil {
		return fmt.Errorf("this username already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	now := time.Now()

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	createdUser, err := s.DB.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println("======== createdUser ========")
	fmt.Printf("ID: %v\n", createdUser.ID)
	fmt.Printf("Name: %v\n", createdUser.Name)
	fmt.Printf("Created At: %v\n", createdUser.CreatedAt)
	fmt.Println()

	err = s.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User has been registered: %v\n", userName)
	return nil
}
