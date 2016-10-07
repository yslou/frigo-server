package main

import (
	"flag"
	"path/filepath"

	"github.com/yslou/frigo-server/lib/model"
	"github.com/yslou/frigo-server/lib/config"
	"github.com/thejerf/suture"
	"os"
)

const (
	exitSuccess            = 0
	exitError              = 1
	exitRestarting         = 2
)

var (
	stop        = make(chan int)
)

type RuntimeOptions struct {
	Host        string
	RedisAddr   string
	Key         string
	Cert        string
}

func parseRuntimeOptions() RuntimeOptions {
	opts := RuntimeOptions {
		Host: ":80",
		RedisAddr: ":6379",
	}

	flag.StringVar(&opts.Host,      "host",         "", "host server address")
	flag.StringVar(&opts.RedisAddr, "redis-addr",   "", "Redis server address")
	flag.StringVar(&opts.Key,       "key",          "", "HTTPS private key")
	flag.StringVar(&opts.Cert,      "cert",         "", "HTTPS certificate")
	flag.Parse()
	return opts
}

func main() {
	mainService := suture.New("main", suture.Spec{
		Log: func(line string) {
			l.Debugln(line)
		},
	})
	mainService.ServeBackground()

	l.SetPrefix("[start] ")
	l.Infoln("Working directroy: ", workingDir())

	opts := parseRuntimeOptions()

	cfg := config.DefaultConfig()

	if opts.Host != "" {
		cfg.Api.RawAddr = opts.Host
	}
	if opts.RedisAddr != "" {
		cfg.Db.RedisAddr = opts.RedisAddr
	}
	if opts.Key != "" {
		cfg.Api.KeyFile = opts.Key
	}
	if opts.RedisAddr != "" {
		cfg.Api.CertFile = opts.Key
	}

	m, e := model.NewModel(&cfg)
	if e != nil {
		l.Fatalln(e)
	}

	mainService.Add(m)

	api, e := NewApiService(m, &cfg)
	if e != nil {
		l.Fatalln(e)
	}

	mainService.Add(api)

	ret := <- stop

	mainService.Stop()

	l.Okln("Quit")

	os.Exit(ret)
}

func workingDir() string {
	dir, _ := filepath.Abs(".")
	return dir
}