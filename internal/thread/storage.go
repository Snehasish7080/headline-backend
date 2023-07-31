package thread

import (
	"context"
	"encoding/json"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zone/headline/internal/models"
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

func (t *ThreadStorage) get(id string, ctx context.Context) ([]*models.Thread, error) {
	session := t.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: t.dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	threads, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx,
				"MATCH(t:Thread)<-[:THREAD]-(o:Opinion) WHERE o.uuid=$id RETURN t{id:t.uuid,.description,.image,.created_at} AS Threads",
				map[string]interface{}{
					"id": id,
				},
			)
			if err != nil {
				return nil, err
			}
			record, err := result.Collect(ctx)

			if err != nil {
				return nil, err
			}

			return record, nil
		})

	if err != nil {
		return nil, err
	}

	var arr []*models.Thread
	for _, opinion := range threads.([]*neo4j.Record) {
		jsonData, _ := json.Marshal(opinion.AsMap()["Threads"])

		var structData models.Thread
		json.Unmarshal(jsonData, &structData)

		arr = append(arr, &models.Thread{
			ID:          structData.ID,
			Description: structData.Description,
			Image:       structData.Image,
			Created_at:  structData.Created_at,
		})
	}

	return arr, nil
}
