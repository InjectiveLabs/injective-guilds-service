package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/xlab/suplog"
)

const (
	apiEnvPrefix     = "GUILDS"
	processEnvPrefix = "GUILDS_PROCESS"
	envFile          = ".env"
)

func panicIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LoadEnvString(name, defaultValue string) string {
	value, exist := os.LookupEnv(name)
	if !exist {
		return defaultValue
	}
	return value
}

func LoadEnvInt(name string, defaultValue int) int {
	value, exist := os.LookupEnv(name)
	if !exist {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	panicIf(err)
	return n
}

func LoadEnvBool(name string, defaultValue bool) bool {
	value, exist := os.LookupEnv(name)
	if !exist {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	panicIf(err)
	return b
}

func LoadEnvDuration(name string, defaultValue time.Duration) time.Duration {
	value, exist := os.LookupEnv(name)
	if !exist {
		return defaultValue
	}

	n, err := time.ParseDuration(value)
	panicIf(err)

	return n
}

type StatsdConfig struct {
	Agent    string
	Disabled bool
	Prefix   string
	Addr     string
	StuckDur time.Duration
	Mocking  bool
}

func loadStatsdConfig(envPrefix string) StatsdConfig {
	return StatsdConfig{
		Agent:    LoadEnvString(fmt.Sprintf("%s_STATSD_AGENT", envPrefix), "telegraf"),
		Disabled: LoadEnvBool(fmt.Sprintf("%s_STATSD_DISABLED", envPrefix), false),
		Prefix:   LoadEnvString(fmt.Sprintf("%s_STATSD_PREFIX", envPrefix), ""),
		Addr:     LoadEnvString(fmt.Sprintf("%s_STATSD_ADDR", envPrefix), ""),
		StuckDur: LoadEnvDuration(fmt.Sprintf("%s_STATSD_STUCK_DUR", envPrefix), 5*time.Minute),
		Mocking:  LoadEnvBool(fmt.Sprintf("%s_STATSD_MOCKING", envPrefix), false),
	}
}

type GuildsAPIServerConfig struct {
	EnvName  string
	LogLevel string

	ListenAddress   string
	TLSCertFilePath string
	TLSKeyFilePath  string

	DBName          string
	DBConnectionURL string

	ExchangeGRPCURL string
	LcdURL          string
	AssetPriceURL   string

	StatsdConfig StatsdConfig
}

func (c GuildsAPIServerConfig) Validate() error {
	return nil
}

func LoadGuildsAPIServerConfig() GuildsAPIServerConfig {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Warningln("load file '.env' failed, going to use env vars")
	}

	return GuildsAPIServerConfig{
		EnvName:  LoadEnvString(fmt.Sprintf("%s_ENV", apiEnvPrefix), "local"),
		LogLevel: LoadEnvString(fmt.Sprintf("%s_LOG_LEVEL", apiEnvPrefix), "DEBUG"),

		ListenAddress:   LoadEnvString(fmt.Sprintf("%s_LISTEN_ADDRESS", apiEnvPrefix), "http://localhost:9900"),
		DBName:          LoadEnvString(fmt.Sprintf("%s_DB_NAME", apiEnvPrefix), "asset_price"),
		DBConnectionURL: LoadEnvString(fmt.Sprintf("%s_DB_CONNECTION_URL", apiEnvPrefix), ""),

		ExchangeGRPCURL: LoadEnvString(fmt.Sprintf("%s_EXCHANGE_GRPC_URL", apiEnvPrefix), "http://localhost:9910"),
		LcdURL:          LoadEnvString(fmt.Sprintf("%s_LCD_URL", apiEnvPrefix), ""),
		AssetPriceURL:   LoadEnvString(fmt.Sprintf("%s_ASSET_PRICE_URL", apiEnvPrefix), ""),

		StatsdConfig: loadStatsdConfig(apiEnvPrefix),
	}
}

type GuildProcessConfig struct {
	EnvName  string
	LogLevel string

	DBName          string
	DBConnectionURL string

	PortfolioUpdateInterval time.Duration
	DisqualifyInterval      time.Duration

	ExchangeGRPCURL string
	AssetPriceURL   string
	LcdURL          string

	StatsdConfig StatsdConfig
}

func (c GuildProcessConfig) Validate() error {
	return nil
}

func LoadGuildsProcessConfig() GuildProcessConfig {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Warningln("load file '.env' failed, going to use env vars")
	}

	return GuildProcessConfig{
		EnvName:  LoadEnvString(fmt.Sprintf("%s_ENV", processEnvPrefix), "local"),
		LogLevel: LoadEnvString(fmt.Sprintf("%s_LOG_LEVEL", processEnvPrefix), "DEBUG"),

		DBName:          LoadEnvString(fmt.Sprintf("%s_DB_NAME", processEnvPrefix), "asset_price"),
		DBConnectionURL: LoadEnvString(fmt.Sprintf("%s_DB_CONNECTION_URL", processEnvPrefix), ""),

		// TODO: Discuss + Update interval
		PortfolioUpdateInterval: LoadEnvDuration(fmt.Sprintf("%s_PORTFOLIO_UPDATE_INTERVAL", processEnvPrefix), time.Hour),
		DisqualifyInterval:      LoadEnvDuration(fmt.Sprintf("%s_DISQUALIFY_INTERVAL", processEnvPrefix), 6*time.Hour),
		StatsdConfig:            loadStatsdConfig(processEnvPrefix),

		ExchangeGRPCURL: LoadEnvString(fmt.Sprintf("%s_EXCHANGE_GRPC_URL", processEnvPrefix), "http://localhost:9910"),
		LcdURL:          LoadEnvString(fmt.Sprintf("%s_LCD_URL", processEnvPrefix), ""),
		AssetPriceURL:   LoadEnvString(fmt.Sprintf("%s_ASSET_PRICE_URL", processEnvPrefix), "https://k8s.mainnet.asset.injective.network"),
	}
}
