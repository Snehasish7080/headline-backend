package user

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zone/headline/pkg/otp"
)

type UserStorage struct {
	db     neo4j.DriverWithContext
	dbName string
}

func NewUserStorage(db neo4j.DriverWithContext, dbName string) *UserStorage {
	return &UserStorage{
		db:     db,
		dbName: dbName,
	}
}

func (u *UserStorage) signUp(firstName string, lastName string, userName string, mobile string, ctx context.Context) (string, error) {
	session := u.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: u.dbName, AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	generatedOtp := otp.EncodeToString(6)
	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			return tx.Run(ctx,
				"CREATE (:User {firstName: $firstName, lastName: $lastName, userName: $userName, mobile: $mobile, otp:$otp, isVerified:$isVerified})",
				map[string]any{"firstName": firstName, "lastName": lastName, "userName": userName, "mobile": mobile, "otp": generatedOtp, "isVerified": false})
		})

	if err != nil {
		return "", err
	}

	return "Otp sent", nil

}
