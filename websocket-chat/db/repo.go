package db

import "github.com/redis/go-redis/v9"

type Repo struct {
	UserRepo    UserRepo
	SessionRepo SessionRepo
}

func NewRepo(rdb *redis.Client) Repo {
	return Repo{
		UserRepo:    NewUserRepo(rdb),
		SessionRepo: NewSessionRepo(rdb),
	}
}
