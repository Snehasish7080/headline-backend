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

}
