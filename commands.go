package main

import "fmt"

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
	err := s.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("The user has been set to %v\n", userName)
	return nil
}
