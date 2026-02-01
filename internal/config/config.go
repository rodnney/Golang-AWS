package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"SERVER_HOST"`
		Port string `mapstructure:"SERVER_PORT"`
	}
	AWS struct {
		Region          string `mapstructure:"AWS_REGION"`
		Endpoint        string `mapstructure:"AWS_ENDPOINT"`
		AccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	}
	SNS struct {
		TopicARN string `mapstructure:"SNS_TOPIC_ARN"`
	}
	SQS struct {
		ValidationURL string `mapstructure:"SQS_VALIDATION_URL"`
		EnrichmentURL string `mapstructure:"SQS_ENRICHMENT_URL"`
		AuditURL      string `mapstructure:"SQS_AUDIT_URL"`
	}
	DynamoDB struct {
		Table string `mapstructure:"DYNAMODB_TABLE"`
	}
	S3 struct {
		Bucket string `mapstructure:"S3_BUCKET"`
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	var cfg Config
	// Manually binding nested structures because viper/mapstructure can be tricky with flat env vars
	cfg.Server.Host = viper.GetString("SERVER_HOST")
	cfg.Server.Port = viper.GetString("SERVER_PORT")
	cfg.AWS.Region = viper.GetString("AWS_REGION")
	cfg.AWS.Endpoint = viper.GetString("AWS_ENDPOINT")
	cfg.AWS.AccessKeyID = viper.GetString("AWS_ACCESS_KEY_ID")
	cfg.AWS.SecretAccessKey = viper.GetString("AWS_SECRET_ACCESS_KEY")
	cfg.SNS.TopicARN = viper.GetString("SNS_TOPIC_ARN")
	cfg.SQS.ValidationURL = viper.GetString("SQS_VALIDATION_URL")
	cfg.SQS.EnrichmentURL = viper.GetString("SQS_ENRICHMENT_URL")
	cfg.SQS.AuditURL = viper.GetString("SQS_AUDIT_URL")
	cfg.DynamoDB.Table = viper.GetString("DYNAMODB_TABLE")
	cfg.S3.Bucket = viper.GetString("S3_BUCKET")

	return &cfg, nil
}
