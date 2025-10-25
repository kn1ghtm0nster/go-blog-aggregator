package main

import (
	"fmt"
	"log"

	"blog-aggregator/internal/config"
)

func main() {
	// read config file
	configFile, err := config.Read()

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// set user to my actual name
	err = configFile.SetUser("Diego")

	if err != nil {
		log.Fatalf("Error setting user in config file: %v", err)
	}

	// read config file again and print contents
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	fmt.Print(cfg.DBUrl, "\n")
	fmt.Print(cfg.CurrentUserName, "\n")
}