package configurations

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"log"
)

const (
	ServerAdress = "localhost:8000"
	DataBaseURI  = "postgresql://postgres:1234@localhost:5432?sslmode=disable"
	//DataBaseURI = ""
	AccrualSystemAdress        = "http://localhost:8080/"
	AccessTokenLiveTimeMinutes = 15
	RefreshTokenLiveTimeDays   = 7
	AccessTokenSecret          = "jdnfksdmfksd"
	RefreshTokenSecret         = "mcmvmkmsdnfsdmfdsjf"
	NumOfWorkers               = 10
	PoolBuffer                 = 1000
	MaxJobRetryCount           = 5
)

type Config struct {
	ServerAdress        string `env:"RUN_ADDRESS"`
	AccrualSystemAdress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Token               ConfigToken
	DataBase            ConfigDatabase
	WorkerPool          ConfigWorkerPool
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

type ConfigWorkerPool struct {
	NumOfWorkers     int `env:"num_of_workers"`
	PoolBuffer       int `env:"pool_buffer"`
	MaxJobRetryCount int `env:"max_job_retry_count"`
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

	fmt.Println(*flagDataBaseURI)
	if *flagDataBaseURI != DataBaseURI {
		dbCfg.DataBaseURI = "postgresql://" + *flagDataBaseURI
	}

	wpConf := ConfigWorkerPool{
		NumOfWorkers:     NumOfWorkers,
		PoolBuffer:       PoolBuffer,
		MaxJobRetryCount: MaxJobRetryCount,
	}

	cfg := Config{
		ServerAdress:        ServerAdress,
		AccrualSystemAdress: AccrualSystemAdress,
		DataBase:            dbCfg,
		Token:               tokenCfg,
		WorkerPool:          wpConf,
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

	cfg.AccrualSystemAdress += "api/orders/"

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
