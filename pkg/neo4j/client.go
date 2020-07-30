package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/spf13/viper"
)

type Neo4j struct {
	Client  neo4j.Driver
	Session neo4j.Session
}

func NewClient() (*Neo4j, error) {
	config := func(conf *neo4j.Config) { conf.Encrypted = false }

	var err error
	n := new(Neo4j)
	n.Client, err = neo4j.NewDriver(
		fmt.Sprintf(
			"%s://%s:%s",
			viper.GetString("NEO4J_PROTOCOL"),
			viper.GetString("NEO4J_HOST"),
			viper.GetString("NEO4J_PORT"),
		),
		neo4j.BasicAuth(
			viper.GetString("NEO4J_USER"),
			viper.GetString("NEO4J_PASSWORD"),
			"",
		),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to neo4j: %w", err)
	}

	// For multidatabase support, set sessionConfig.DatabaseName to requested database
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	n.Session, err = n.Client.NewSession(sessionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start neo4j session: %w", err)
	}

	return n, nil
}
