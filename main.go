package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"blog-aggregator/internal/commands"
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"blog-aggregator/internal/state"
)

func main() {
	// read config file
	configFile, err := config.Read()

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// open database connection
	db, err := sql.Open("postgres", configFile.DBUrl)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state.State{
		DB:     dbQueries,
		Config: &configFile,
	}

	cmds := &commands.Commands{}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerResetUsers)
	cmds.Register("users", commands.HandlerGetUsers)
	cmds.Register("agg", commands.HandlerAgg)

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