package model

import (
	"encoding/json"
	"io"
)

// User blahblah
type User struct {
	// TODO NEVER store password plaintext
	Login       string   `json:"login,omitempty"`
	Password    string   `json:"password,omitempty"`
	FullName    string   `json:"fullName,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Friends     []string `json:"friends,omitempty"`
}

// JSONUser blahblah
func JSONUser(s string) (User, error) {
	var u User
	err := json.Unmarshal([]byte(s), &u)
	return u, err
}

// ReadUser blahblah
func ReadUser(r io.Reader) (User, error) {
	var u User
	err := json.NewDecoder(r).Decode(&u)
	return u, err
}

// LatLng is latidute and longtitue
type LatLng struct {
	Login string  `json:"login,omitempty"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}

// ReadLatLng read latidute and longtitue from io stream
func ReadLatLng(r io.Reader) (LatLng, error) {
	var u LatLng
	err := json.NewDecoder(r).Decode(&u)
	return u, err
}

// JSONLatLng is shortcut
func JSONLatLng(lat, lng float64) (ll LatLng) {
	ll.Lat = lat
	ll.Lng = lng
	return
}
