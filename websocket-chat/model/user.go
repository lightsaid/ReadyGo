package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPswdNotHashed = errors.New("密码无法加密")

// User 用户
type User struct {
	ID       int    `redis:"id" json:"id,omitempty"`
	Nickname string `redis:"nickname" json:"nickname"`
	Password string `redis:"password" json:"password"`
	Avatar   string `redis:"avatar" json:"avatar"`
}

// NewUser 创建一个User，创建 uuid 赋值给ID，密码哈希化
func NewUser(id int, nickname, password, avatar string) (*User, error) {
	user := &User{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
	}

	err := user.GenHashedPswd()
	if err != nil {
		return nil, ErrPswdNotHashed
	}

	return user, nil
}

func (user *User) GenHashedPswd() error {
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
