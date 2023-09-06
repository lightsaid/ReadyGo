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

// SessionRepo db 操作接口定义
type SessionRepo interface {
	Save(sess *model.Session) error
	Get(id string) (*model.Session, error)
}

var _ SessionRepo = (*sessionRepo)(nil)

type sessionRepo struct {
	rdb *redis.Client
}

func NewSessionRepo(rdb *redis.Client) SessionRepo {
	return &sessionRepo{rdb: rdb}
}

func (store *sessionRepo) Save(sess *model.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dataBytes, err := json.Marshal(sess)
	if err != nil {
		return err
	}

	return store.rdb.HSet(
		ctx,
		Key_Session,
		sess.ID.String(), string(dataBytes),
	).Err()
}

func (store *sessionRepo) Get(id string) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var sess model.Session
	result, err := store.rdb.HGet(ctx, Key_Session, id).Result()
	if err != nil {
		log.Println("HGet err: ", err.Error())
		switch {
		case err == redis.Nil:
			return nil, fmt.Errorf("key不存在: %w", redis.Nil)
		case err != nil:
			return nil, fmt.Errorf("错误：%w", err)
		case result == "":
			return nil, fmt.Errorf("值是空字符串: %w", err)
		}
	}

	err = json.Unmarshal([]byte(result), &sess)
	if err != nil {
		log.Println("json unmarshal err: ", err.Error())
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &sess, nil
}
