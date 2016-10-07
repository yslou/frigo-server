package model

import (
	"github.com/yslou/frigo-server/lib/logger"
)

var (
	l         = logger.DefaultLogger.NewFacility("model", "model")
)

func init() {
	l.Debugf("-----------------------------main.init()")
}