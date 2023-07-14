package user

import (
	"context"
	"errors"

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

	isMobileExist := u.mobileExists(mobile, ctx)

	if isMobileExist {
		return "", errors.New("mobile already exists")
	}
	isUserNameExist := u.userNameExists(userName, ctx)

	if isUserNameExist {
		return "", errors.New("username already exists")
	}

	session := u.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: u.dbName, AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	generatedOtp := otp.EncodeToString(6)
	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			return tx.Run(ctx,
				"CREATE (:User {firstName: $firstName, lastName: $lastName, userName: $userName, mobile: $mobile, otp:$otp, isVerified:$isVerified, isComplete:$isComplete})",
				map[string]any{"firstName": firstName, "lastName": lastName, "userName": userName, "mobile": mobile, "otp": generatedOtp, "isVerified": false, "isComplete": false})
		})

	if err != nil {
		return "", err
	}

	return "Otp sent", nil

}

func (u *UserStorage) mobileExists(mobile string, ctx context.Context) bool {
	session := u.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: u.dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, _ := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (interface{}, error) {
			result, err := tx.Run(ctx,
				"MATCH (u:User {mobile:$mobile}) RETURN u.mobile AS mobile",
				map[string]interface{}{
					"mobile": mobile,
				},
			)
			if err != nil {
				return nil, err
			}
			record, err := result.Single(ctx)
			if err != nil {
				return nil, err
			}
			mobile, _ := record.Get("mobile")
			return mobile.(string), nil
		})

	return result != nil

}
func (u *UserStorage) userNameExists(userName string, ctx context.Context) bool {
	session := u.db.NewSession(ctx, neo4j.SessionConfig{DatabaseName: u.dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, _ := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (interface{}, error) {
			result, err := tx.Run(ctx,
				"MATCH (u:User {userName:$userName}) RETURN u.userName AS userName",
				map[string]interface{}{
					"userName": userName,
				},
			)
			if err != nil {
				return nil, err
			}
			record, err := result.Single(ctx)
			if err != nil {
				return nil, err
			}
			userName, _ := record.Get("userName")
			return userName.(string), nil
		})

	return result != nil

}
