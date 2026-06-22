package main

import "fmt"

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
