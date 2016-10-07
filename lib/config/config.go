package config

type Config struct {
	Db          DbConfig    `json:"db"`
	Api         ApiConfig   `json:"api"`
}

type DbConfig struct {
	RedisAddr   string      `json:"redisAddr" default:":6379"`
}

type ApiConfig struct {
	RawAddr     string      `json:"address" default:":443"`
	KeyFile     string      `json:"key" default:"https-key.pem"`
	CertFile    string      `json:"cert" default:"https-cert.pem"`
}

func DefaultConfig() Config {
	c :=
	 Config {
		Db:     DbConfig{
			":6379",
		},
		Api:    ApiConfig{
			":443",
			"https-key.pem",
			"https-cert.pem",
		},
	}
	return c
}