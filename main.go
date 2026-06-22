package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/amayakmt/blog-aggregator/internal/config"
	"github.com/amayakmt/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

// state is threaded through every command handler so each handler has
// access to both the database and the on-disk config without globals.
type state struct {
	DB     *database.Queries
	Config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	dbQueries := database.New(db)

	mainState := state{}
	mainState.DB = dbQueries
	mainState.Config = &cfg

	commandsInit := commands{
		RegisteredCommands: map[string]func(*state, command) error{},
	}

	commandsInit.register("login", handlerLogin)
	commandsInit.register("register", handlerRegister)
	commandsInit.register("reset", handlerReset)
	commandsInit.register("users", handlerGetUsers)

	args := os.Args
	// os.Args[0] is the binary name; the command name must be os.Args[1].
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
