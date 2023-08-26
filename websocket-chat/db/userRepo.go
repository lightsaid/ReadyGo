package db

import (
	"context"
	"encoding/json"
	"readygo/wesocket-chat/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRepo interface {
	Save(user model.User) error
	Get(id int) (*model.User, error)
	List() ([]*model.User, error)
}

var _ UserRepo = (*userRepo)(nil)

type userRepo struct {
	rdb *redis.Client
}

func NewUserRepo(rdb *redis.Client) UserRepo {
	return &userRepo{rdb: rdb}
}

func (store *userRepo) Save(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dataBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return store.rdb.HSet(
		ctx,
		Key_Users,
		Key_User(user.ID), string(dataBytes),
	).Err()
}

func (store *userRepo) Get(id int) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user model.User
	err := store.rdb.HGet(ctx, Key_Users, Key_User(id)).Scan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (store *userRepo) List() ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*model.User
	err := store.rdb.HGetAll(ctx, Key_Users).Scan(users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
