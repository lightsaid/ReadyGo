package db_test

import (
	"os"
	"readygo/wesocket-chat/db"
	"testing"

	"github.com/redis/go-redis/v9"
)

var testDB *redis.Client
var testUserRepo db.UserRepo
var testSessRepo db.SessionRepo

func TestMain(m *testing.M) {
	testDB = db.NewRedisClient()

	testUserRepo = db.NewUserRepo(testDB)

	testSessRepo = db.NewSessionRepo(testDB)

	os.Exit(m.Run())
}
