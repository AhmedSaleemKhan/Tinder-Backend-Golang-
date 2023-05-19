package util

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBDriver               string        `mapstructure:"DB_DRIVER"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	MigrationURL           string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress      string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress      string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration   time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	AuthenticationFilePath string        `mapstructure:"AUTHENTICATION_FILE_PATH"`
	BucketAccessKey        string        `mapstructure:"BUCKET_ACCESS_KEY"`
	BucketSecretKey        string        `mapstructure:"BUCKET_SECRET_KEY"`
	LinodeRegion           string        `mapstructure:"LINODE_REGION"`
	BucketName             string        `mapstructure:"BUCKET_NAME"`
	BucketEndpoint         string        `mapstructure:"BUCKET_ENDPOINT"`
	BucketACL              string        `mapstructure:"BUCKET_ACL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// export ENV=env.production // To load the production env
	if len(os.Getenv("ENV")) > 0 {
		viper.SetConfigName(os.Getenv("ENV"))
	} else {
		viper.SetConfigName("env.local")
	}
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
