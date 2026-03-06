package config

const (
	BaseConfigPath   = "../../config/"
	ConfigFileSuffix = "Config.yml"
)

type App struct {
	AppName string `yaml:"appName"`
	Port    int    `yaml:"port"`
	Env     string `yaml:"env,omitempty"`
}

type Logger struct {
	ConsoleLoggingEnabled bool   `yaml:"consoleLoggingEnabled"`
	FileName              string `yaml:"fileName"`
	ErrorFileName         string `yaml:"errorFileName"`
	MaxSizeInMB           int    `yaml:"maxSizeInMB"`
	MaxBackups            int    `yaml:"maxBackups"`
	MaxAgeInDays          int    `yaml:"maxAgeInDays"`
	Compress              bool   `yaml:"compress"`
}

type RateLimit struct {
	LoginService int `yaml:"loginService"`
	OrderService int `yaml:"orderService"`
}

type Session struct {
	ExpiryTime          int    `yaml:"expiryTime"`
	Secret              string `yaml:"secret"`
	SessionsPerPlatform int    `yaml:"sessionsPerPlatform"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Username string `yaml:"user"`
}

type MongoDb struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DBName   string `yaml:"db_name"`
}

type Config struct {
	AppConfig       App       `yaml:"app"`
	LoggerConfig    Logger    `yaml:"logger"`
	RateLimitConfig RateLimit `yaml:"rateLimit"`
	SessionConfig   Session   `yaml:"session"`
	SMRedisConfig   Redis     `yaml:"smRedis"`
	MongoDbConfig   MongoDb   `yaml:"mongoDb"`
}
