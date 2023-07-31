package opinion

import (
	"context"
	"encoding/json"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zone/headline/internal/models"
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
	session := o.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: o.dbName, AccessMode: neo4j.AccessModeWrite})
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

func (o *OpinionStorage) getOpinions(userName string, ctx context.Context) ([]*models.Opinion, error) {
	session := o.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: o.dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	opinions, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx,
				`
				MATCH (u:User{userName:$userName})-[:CREATED_BY]->(o:Opinion) CALL{WITH o MATCH (o)-[:THREAD]->(t:Thread) RETURN t LIMIT 2} WITH o.uuid AS id,o.image AS image, o.description AS description, o.created_at AS created_at,{id:o.uuid,image:o.image,description:o.description,created_at:o.created_at,threads:collect(t{id:t.uuid,.image,.description,.created_at})} AS Opinion  RETURN Opinion
				UNION MATCH(o:Opinion) CALL{WITH o MATCH(o)-[:THREAD]->(t:Thread)<-[:THREAD_BY]-(u:User{userName:$userName}) WHERE NOT (u)-[:CREATED_BY]->(o) RETURN t LIMIT 2} WITH o.uuid AS id,o.image AS image, o.description AS description, o.created_at AS created_at,  {id:o.uuid,image:o.image,description:o.description,created_at:o.created_at,threads:collect(t{id:t.uuid,.image,.description,.created_at})} AS Opinion  RETURN Opinion
				`,
				map[string]interface{}{
					"userName": userName,
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

	var arr []*models.Opinion
	for _, opinion := range opinions.([]*neo4j.Record) {
		jsonData, _ := json.Marshal(opinion.AsMap()["Opinion"])
		var structData models.Opinion
		json.Unmarshal(jsonData, &structData)

		arr = append(arr, &models.Opinion{
			ID:          structData.ID,
			Description: structData.Description,
			Image:       structData.Image,
			Created_at:  structData.Created_at,
			Threads:     structData.Threads,
		})
	}

	return arr, nil

}
