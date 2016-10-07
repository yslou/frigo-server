package model

import (
	"github.com/yslou/frigo-server/lib/db"
	"github.com/yslou/frigo-server/lib/config"
	"github.com/thejerf/suture"
)

type Model struct {
	*suture.Supervisor

	cfg     *config.Config
	db      db.Interface
}

func NewModel(cfg *config.Config) (*Model, error) {
	m := &Model {
		Supervisor: suture.New("model", suture.Spec{
			Log: func(line string) {
				l.Debugln(line)
			},
		}),
		cfg:    cfg,
	}
	var e error
	m.db, e = db.NewInstance(cfg)
	if e != nil {
		l.Fatalln(e)
		return nil, e
	}
	return m, nil
}