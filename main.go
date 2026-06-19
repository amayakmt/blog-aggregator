package main

import (
	"fmt"
	"os"

	"github.com/amayakmt/blog-aggregator/internal/config"
)

type state struct {
	Config *config.Config
}

func main() {

	mainState := state{}
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	mainState.Config = &cfg

	commandsInit := commands{
		RegisteredCommands: map[string]func(*state, command) error{},
	}

	commandsInit.register("login", handlerLogin)

	args := os.Args
	if len(os.Args) < 2 {
		fmt.Println("no arguments provided")
		os.Exit(1)
	}

	cmd := command{
		Name:      args[1],
		Arguments: args[2:],
	}

	err = commandsInit.run(&mainState, cmd)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

}
