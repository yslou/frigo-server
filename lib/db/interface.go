package db

import (
	"github.com/yslou/frigo-server/lib/config"
	"github.com/yslou/frigo-server/lib/model/def"
)

// Interface defines
type Interface interface {
	AddUser(login, pwd string) error
	UserExists(login string) bool
	GetUser(login string) (model.User, error)
	UpdateUser(user model.User) error
	UpdateLatLng(login string, ll model.LatLng) error
	GetLatLng(login string) model.LatLng
	GetFriends(login string) []string
	UpdateFriends(login string, friends []string)
	AddFriend(login string, friend string)
	HasFriend(login string, friend string) bool
}

// NewInstance returns db instance
func NewInstance(cfg *config.Config) (Interface, error) {
	return NewRedisDbInstance(cfg.Db.RedisAddr)
}

