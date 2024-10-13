package config

type Config struct {
	App struct {
		Env     string
		BaseURL string
		Version string
		Debug   bool
	}
	Server struct {
		Port int
	}
	Database struct {
		Dsn                 string
		AutoMigrate         bool
		MaxOpenConns        int
		MaxIdleConns        int
		ConnTimeoutDuration int
		ConnMaxIdleDuration int
		ConnMaxLifeDuration int
	}
}
