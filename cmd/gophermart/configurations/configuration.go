package configurations

const (
	ServerAdress = "localhost:8080"
	BaseURL      = "http://localhost:8080/"
	DataBaseURI  = ""
	// DataBaseURI  = "postgresql://postgres:1234@localhost:5432?sslmode=disable"
)

type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
	DataBase     ConfigDatabase
}

type ConfigDatabase struct {
	DataBaseURI string `env:"DATABASE_DSN"`
}
