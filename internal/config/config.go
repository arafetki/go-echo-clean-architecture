package config

type Config struct {
	App struct {
		Env     string
		BaseURL string
		Version string
		Debug   bool
	}
	Server struct {
		Port                 int
		ReadTimeoutDuration  int
		WriteTimeoutDuration int
		IdleTimeoutDuration  int
		ShutdownPeriod       int
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

	JWT struct {
		SecretKey string
	}
}
