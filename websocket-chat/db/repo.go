package db

import "github.com/redis/go-redis/v9"

var RepoObj *Repo

type Repo struct {
	UserRepo    UserRepo
	SessionRepo SessionRepo
}

// NewRepo 实例化一个仓库，将实例赋值给RepoObj，暴露给外包使用，同时并返回
func NewRepo(rdb *redis.Client) *Repo {
	store := &Repo{
		UserRepo:    NewUserRepo(rdb),
		SessionRepo: NewSessionRepo(rdb),
	}

	RepoObj = store

	return store
}
