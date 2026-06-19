package main

import (
	"fmt"

	"github.com/amayakmt/blog-aggregator/internal/config"
)

func main() {

	// Before ----------------------------------------
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("--- Initial Config ---\n")
	fmt.Printf("DB URL: %s\n", cfg.DBURL)
	fmt.Printf("username: %s\n", cfg.CurrentUserName)

	// Set User Name --------------------------------
	userName := "amayakmt"
	err = cfg.SetUser(userName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	// After ----------------------------------------
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("--- New Config ---\n")
	fmt.Printf("DB URL: %s\n", cfg.DBURL)
	fmt.Printf("username: %s\n", cfg.CurrentUserName)

}
