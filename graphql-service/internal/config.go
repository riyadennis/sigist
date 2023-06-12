package internal

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

// Config is the configuration to run the service
// args are parsed from go-arg, https://github.com/alexflint/go-arg
// Add here service config arguments and add the specific arg tag
type Config struct {
	Env            string `arg:"env:ENVIRONMENT" validate:"required,notblank"`
	Port           string `arg:"env:PORT" validate:"required,hostname_port"`
	LogLevel       string `arg:"env:LOG_LEVEL" validate:"required,notblank"`
	DBFile         string `arg:"env:DB_FILE" default:"../environment/db/user-feedback.sqlite"`
	MigrationsPath string `arg:"env:MIGRATIONS_PATH" default:"migrations"`
	KafkaBroker    string `arg:"env:KAFKA_BROKER" validate:"required,notblank"`
	KafkaTopic     string `arg:"env:KAFKA_TOPIC" validate:"required,notblank"`
}

// NewConfig return a new instance of Config
func NewConfig() (Config, error) {
	var conf Config

	validate := validator.New()
	err := validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return conf, fmt.Errorf("failed to register \"NotBlank\" validator")
	}
	arg.MustParse(&conf)
	err = isValid(validate, conf)

	return conf, err
}

func isValid(validate *validator.Validate, conf Config) error {
	if err := validate.Struct(conf); err != nil {
		return fmt.Errorf("validating struct: %w", err)
	}

	return nil
}
