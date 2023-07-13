package storage

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func BootstrapNeo4j(dbUri string, dbName string, dbUser string, dbPassword string, timeout time.Duration) (neo4j.DriverWithContext, neo4j.SessionWithContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	if err != nil {
		return nil, nil, err
	}
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, nil, err
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: dbName})
	println("Connected to Neo4j ... :)")
	return driver, session, nil
}

func CloseNeo4j(driver neo4j.DriverWithContext) error {
	return driver.Close(context.Background())
}
func CloseNeo4jSession(session neo4j.SessionWithContext) error {
	return session.Close(context.Background())
}
