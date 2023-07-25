package thread

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ThreadStorage struct {
	db     neo4j.DriverWithContext
	dbName string
}

func NewThreadStorage(db neo4j.DriverWithContext, dbName string) *ThreadStorage {
	return &ThreadStorage{
		db:     db,
		dbName: dbName,
	}
}

func (t *ThreadStorage) create(userName string, opinionId string, threadField map[string]interface{}, ctx context.Context) (string, error) {
	session := t.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: t.dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			return tx.Run(ctx,
				"MATCH (u:User {userName:$userName}) MATCH (o:Opinion {uuid: $opinionId}) CREATE (t:Thread $fields)<-[:THREAD_BY]-(u) CREATE (t)<-[:THREAD]-(o) ",
				map[string]interface{}{
					"userName":  userName,
					"opinionId": opinionId,
					"fields":    threadField,
				},
			)
		})

	if err != nil {
		return "", err
	}

	return "Created Successfully", nil
}
