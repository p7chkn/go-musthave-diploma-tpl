package configurations

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

const (
	ServerAdress        = "localhost:8080"
	DataBaseURI         = "postgresql://postgres:1234@localhost:5432?sslmode=disable"
	AccrualSystemAdress = ""
	// DataBaseURI  = "postgresql://postgres:1234@localhost:5432?sslmode=disable"
	AccessTokenLiveTimeMinutes = 15
	RefreshTokenLiveTimeDays   = 7
	AccessTokenSecret          = "jdnfksdmfksd"
	RefreshTokenSecret         = "mcmvmkmsdnfsdmfdsjf"
)

type Config struct {
	ServerAdress        string `env:"SERVER_ADDRESS"`
	AccrualSystemAdress string `env:"ACCRUAL_SYSTEM_ADDRESS "`
	Token               ConfigToken
	DataBase            ConfigDatabase
}

type ConfigToken struct {
	AccessTokenLiveTimeMinutes int    `env:"ACCESS_TOKEN_LIVE_TIME_MINUTES"`
	RefreshTokenLiveTimeDays   int    `env:"REFRESH_TOKEN_LIVE_TIME_DAYS"`
	AccessTokenSecret          string `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret         string `env:"REFRESH_TOKEN_SECRET"`
}

type ConfigDatabase struct {
	DataBaseURI string `env:"DATABASE_DSN"`
}

func New() *Config {
	dbCfg := ConfigDatabase{
		DataBaseURI: DataBaseURI,
	}

	tokenCfg := NewTokenConfig()

	flagServerAdress := flag.String("a", ServerAdress, "server adress")
	flagDataBaseURI := flag.String("d", DataBaseURI, "URI for database")
	flagAccrualSystemAdress := flag.String("r", AccrualSystemAdress, "URL for accrual system")
	flag.Parse()

	if *flagDataBaseURI != DataBaseURI {
		dbCfg.DataBaseURI = *flagDataBaseURI
	}

	cfg := Config{
		ServerAdress: ServerAdress,
		DataBase:     dbCfg,
		Token:        tokenCfg,
	}

	if *flagServerAdress != ServerAdress {
		cfg.ServerAdress = *flagServerAdress
	}
	if *flagAccrualSystemAdress != AccrualSystemAdress {
		cfg.AccrualSystemAdress = *flagAccrualSystemAdress
	}

	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}

func NewTokenConfig() ConfigToken {
	tokenCfg := ConfigToken{
		AccessTokenLiveTimeMinutes: AccessTokenLiveTimeMinutes,
		RefreshTokenLiveTimeDays:   RefreshTokenLiveTimeDays,
		AccessTokenSecret:          AccessTokenSecret,
		RefreshTokenSecret:         RefreshTokenSecret,
	}
	err := env.Parse(&tokenCfg)
	if err != nil {
		log.Fatal()
	}
	return tokenCfg
}
