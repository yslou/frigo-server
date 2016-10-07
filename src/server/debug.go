package main

import (
	"github.com/yslou/frigo-server/lib/logger"
	"strings"
	"os"
)

var (
	l       = logger.DefaultLogger.NewFacility("main", "main")
	apil    = logger.DefaultLogger.NewFacility("api", "API")
)

func init() {
	l.SetDebug("main", strings.Contains(os.Getenv("TRACE"), "main") || os.Getenv("TRACE") == "all")
	l.SetDebug("http", strings.Contains(os.Getenv("TRACE"), "http") || os.Getenv("TRACE") == "all")
}