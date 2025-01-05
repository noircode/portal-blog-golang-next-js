package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type CloudflareR2 struct {
	Name string `json:"name"`
	ApiKey string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Token string `json:"token"`
	AccountID string `json:"account_id"`
	PublicURL string `json:"public_url"`
}

type Config struct {
	App  App
	Psql PsqlDB
	R2 CloudflareR2
}

// NewConfig creates and returns a new Config instance.
// It initializes the configuration by reading values from environment variables
// using the viper package. The function populates both the App and PsqlDB
// structs within the Config.
//
// Returns:
//   - *Config: A pointer to a new Config instance with all fields populated
//     from the corresponding environment variables.
func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),

			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),
		},

		Psql: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetString("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USER"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			DBName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_OPEN_CONNECTION"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
		},
		
		R2: CloudflareR2{
			Name:       viper.GetString("CLOUDFLARE_R2_BUCKET_NAME"),
      ApiKey:      viper.GetString("CLOUDFLARE_R2_API_KEY"),
      ApiSecret: viper.GetString("CLOUDFLARE_R2_API_SECRET"),
      Token:      viper.GetString("CLOUDFLARE_R2_TOKEN"),
      AccountID:  viper.GetString("CLOUDFLARE_R2_ACCOUNT_ID"),
      PublicURL: viper.GetString("CLOUDFLARE_R2_PUBLIC_URL"),
		},
	}
}
