package user

import (
	"context"
	"fmt"

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
	// u.userExists(mobile, ctx)

	return "Otp sent", nil

}

func (u *UserStorage) userExists(mobile string, ctx context.Context) bool {
	session := u.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: u.dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, _ := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx,
				"MATCH (u:User ${mobile:$mobile}) RETURN u.mobile AS mobile",
				map[string]interface{}{
					"mobile": mobile,
				},
			)
			if err != nil {
				return nil, err
			}
			record, err := result.Single(ctx)
			fmt.Printf("record %v", record)
			if err != nil {
				return nil, err
			}
			mobile, _ := record.Get("mobile")

			return mobile.(string), nil
		})

	fmt.Printf("result %v", result)
	// if err != nil {
	// 	return false
	// }

	return true
}
