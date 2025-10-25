package main

import (
	"fmt"
	"log"
	"os"

	"blog-aggregator/internal/commands"
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/state"
)

func main() {
	// read config file
	configFile, err := config.Read()

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	programState := &state.State{
		Config: &configFile,
	}

	cmds := &commands.Commands{}
	cmds.Register("login", commands.HandlerLogin)

	// ensure we have at least one command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> [args...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}

	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}

	cmd := commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(programState, cmd)

	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}