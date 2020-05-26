package feedado

import (
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/user"
)

func initUserService(db *mongo.Client) *user.Service {
	repository := initMongoUserRepository(db)

	return user.NewService(repository)
}
