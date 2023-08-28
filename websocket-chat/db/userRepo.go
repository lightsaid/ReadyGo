package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"readygo/wesocket-chat/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRepo interface {
	Save(user *model.User) error
	Get(nickname string) (*model.User, error)
	List() ([]*model.User, error)
}

var _ UserRepo = (*userRepo)(nil)

type userRepo struct {
	rdb *redis.Client
}

func NewUserRepo(rdb *redis.Client) UserRepo {
	return &userRepo{rdb: rdb}
}

func (store *userRepo) Save(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dataBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return store.rdb.HSet(
		ctx,
		Key_Users,
		Key_User(user.Nickname), string(dataBytes),
	).Err()
}

func (store *userRepo) Get(nickname string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user model.User
	result, err := store.rdb.HGet(ctx, Key_Users, Key_User(nickname)).Result()
	if err != nil {
		log.Println("HGet err: ", err.Error())
		switch {
		case err == redis.Nil:
			return nil, fmt.Errorf("key不存在")
		case err != nil:
			return nil, fmt.Errorf("错误：%w", err)
		case result == "":
			return nil, fmt.Errorf("值是空字符串")
		}
	}

	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		log.Println("json unmarshal err: ", err.Error())
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &user, nil
}

func (store *userRepo) List() ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*model.User
	result, err := store.rdb.HGetAll(ctx, Key_Users).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key不存在: %w", err)
		}
		return nil, fmt.Errorf("HGetAll error: %w", err)
	}

	for _, val := range result {
		var user model.User
		err = json.Unmarshal([]byte(val), &user)
		if err != nil {
			log.Println("json unmarshal err: ", err.Error())
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}
