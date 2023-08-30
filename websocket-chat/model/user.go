package model

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrPswdNotHashed = errors.New("密码无法加密")

// User 用户
type User struct {
	ID       uuid.UUID `redis:"id" json:"id,omitempty"`
	Nickname string    `redis:"nickname" json:"nickname"`
	Password string    `redis:"password" json:"password"`
	Avatar   string    `redis:"avatar" json:"avatar"`
}

// NewUser 创建一个User，创建 uuid 赋值给ID，密码哈希化
func NewUser(nickname, password, avatar string) (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("uuid.NewRandom failed: ", err)
		return nil, err
	}
	user := &User{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
		Password: password,
	}

	err = user.GenHashedPswd()
	if err != nil {
		return nil, ErrPswdNotHashed
	}

	return user, nil
}

func (user *User) GenHashedPswd() error {
	if user.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	dataBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(dataBytes)

	return nil
}

func (user *User) CheckedPswd(pswd string, hashedPswd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPswd), []byte(pswd))
	return err == nil
}
