package model

import (
	"encoding/json"
	"io"
	"math/rand"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// NewToken blahblah
func NewToken(c Cert) string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Cert blahblah
// show your Certificate when login
type Cert struct {
	// Login name is the identity
	Login string `json:"login"`
	// TODO NEVER store password plaintext
	Password string `json:"password"`
}

// ReadCert blahblah
func ReadCert(r io.Reader) (Cert, error) {
	var c Cert
	err := json.NewDecoder(r).Decode(&c)
	return c, err
}

// Ticket blahblah
// If pass is okay, we will issue a ticket
type Ticket struct {
	Login string `json:"login"`
	Token string `json:"token"`
}
