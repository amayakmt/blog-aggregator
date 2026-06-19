package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/amayakmt/blog-aggregator/internal/config"
	"github.com/amayakmt/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	DB     *database.Queries
	Config *config.Config
}

// main ----------------------------------------------------------------
func main() {
	dbURL := "postgres://estyle-163:@localhost:5432/gator?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	dbQueries := database.New(db)

	mainState := state{}
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	mainState.DB = dbQueries
	mainState.Config = &cfg

	commandsInit := commands{
		RegisteredCommands: map[string]func(*state, command) error{},
	}

	commandsInit.register("login", handlerLogin)
	commandsInit.register("register", handlerRegister)

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
