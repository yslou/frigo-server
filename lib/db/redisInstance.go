package db

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/yslou/frigo-server/lib/model/def"
)


// Instance xx
type Instance struct {
	conn redis.Conn
}

// NewRedisDbInstance returns RedisDbInstance
func NewRedisDbInstance(addr string) (*Instance, error) {
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		fmt.Println("failed connect to redis server")
		return nil, err
	}
	ins := Instance {
		conn: conn,
	}
	return &ins, nil
}

func userKey(login string) string {
	return login + "_profile"
}

func llKey(login string) string {
	return login + "_latlng"
}

// AddUser add user
func (db *Instance) AddUser(login, pwd string) error {
	return db.UpdateUser(model.User{
		Login:       login,
		Password:    pwd,
		DisplayName: login,
	})
}

// AddUser add user
func (db *Instance) UpdateUser(user model.User) error {
	js, _ := json.Marshal(user)
	_, err := db.conn.Do("SET", userKey(user.Login), string(js))
	return err
}

// UserExists check if user exists
func (db *Instance) UserExists(login string) bool {
	res, err := redis.Bool(db.conn.Do("EXISTS", userKey(login)))
	return err == nil && res
}

// GetUser pulls user profile from db
func (db *Instance) GetUser(login string) (user model.User, err error) {
	if !db.UserExists(login) {
		return user, fmt.Errorf("ErrUserNotExist")
	}
	rec, err := redis.String(db.conn.Do("GET", userKey(login)))
	user, err = model.JSONUser(rec)
	if err != nil {
		return user, fmt.Errorf("ErrFormat")
	}
	return
}

// UpdateLatLng update user's position
func (db *Instance) UpdateLatLng(login string, ll model.LatLng) error {
	if !db.UserExists(login) {
		return fmt.Errorf("ErrUserNotExist")
	}

	_, err := db.conn.Do("GEOADD", llKey(login), ll.Lng, ll.Lat, login)
	if err != nil {
		fmt.Println("db: GEOADD error ", err)
	}
	return err
}

// GetLatLng get user's position
func (db *Instance) GetLatLng(login string) model.LatLng {
	// TODO error handling
	// if !db.UserExists(login) {
	// 	return fmt.Errorf("ErrUserNotExist")
	// }

	ret, err := db.conn.Do("GEOPOS", llKey(login), login)
	if err != nil {
		return model.LatLng{}
	}
	arr, ok := ret.([]interface{})
	if !ok || len(arr) <= 0 {
		return model.JSONLatLng(0.0, 0.0)
	}
	pos := arr[0].([]interface{})
	lat := 0.0
	lng := 0.0
	if len(pos) >= 2 {
		lat, _ = redis.Float64(pos[1], err)
	}
	if len(pos) >= 1 {
		lng, _ = redis.Float64(pos[0], err)
	}
	return model.JSONLatLng(lat, lng)
}

// GetFriends returns array of login names
func (db *Instance) GetFriends(login string) []string {
	user, err := db.GetUser(login)
	if err != nil {
		return []string{}
	}
	return user.Friends
}

// UpdateFriends accept array of login names
func (db *Instance) UpdateFriends(login string, friends []string) {
	user, err := db.GetUser(login)
	if err != nil {
		return
	}
	user.Friends = friends
	db.UpdateUser(user)
}

// AddFriend add friend
func (db *Instance) AddFriend(login string, friend string) {
	if !db.UserExists(login) || db.HasFriend(login, friend) {
		return
	}
	user, _ := db.GetUser(login)
	user.Friends = append(user.Friends, friend)
	db.UpdateUser(user)
}

// HasFriend check friendship
func (db *Instance) HasFriend(login string, friend string) bool {
	user, err := db.GetUser(login)
	if err != nil {
		return false
	}
	for _, f := range user.Friends {
		if f == friend {
			return true
		}
	}
	return false
}
