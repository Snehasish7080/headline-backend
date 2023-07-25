package opinion

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type OpinionStorage struct {
	db     neo4j.DriverWithContext
	dbName string
}

func NewOpinionStorage(db neo4j.DriverWithContext, dbName string) *OpinionStorage {
	return &OpinionStorage{
		db:     db,
		dbName: dbName,
	}
}

func (o *OpinionStorage) create(userName string, opinionField map[string]interface{}, ctx context.Context) (string, error) {
	session := o.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: o.dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			return tx.Run(ctx,
				"MATCH (u:User {userName:$userName}) CREATE (:Opinion $fields)<-[:CREATED_BY]-(u)",
				map[string]interface{}{
					"userName": userName,
					"fields":   opinionField,
				},
			)
		})

	if err != nil {
		return "", err
	}

	return "Created Successfully", nil

}
