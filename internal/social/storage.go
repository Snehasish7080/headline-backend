package social

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type SocialStorage struct {
	db     neo4j.DriverWithContext
	dbName string
}

func NewSocialStorage(db neo4j.DriverWithContext, dbName string) *SocialStorage {
	return &SocialStorage{
		db:     db,
		dbName: dbName,
	}
}

func (s *SocialStorage) follow(userName string, followingUser string, ctx context.Context) (string, error) {
	session := s.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: s.dbName, AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			return tx.Run(ctx,
				"MATCH (u:User {userName:$userName}) MATCH(f:User{userName:$followingUser}) CREATE (u)-[:FOLLOWING]->(f)",
				map[string]interface{}{
					"userName":      userName,
					"followingUser": followingUser,
				},
			)
		})

	if err != nil {
		return "", err
	}

	return "Followed Successfully", nil
}
