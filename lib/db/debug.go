package db

import (
	"github.com/yslou/frigo-server/lib/logger"
)

var (
	l         = logger.DefaultLogger.NewFacility("db", "database")
)

func init() {
	l.Debugf("-----------------------------main.init()")
}