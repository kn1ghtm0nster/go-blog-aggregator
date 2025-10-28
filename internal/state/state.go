package state

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}