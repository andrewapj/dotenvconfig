package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"io/fs"
	"os"
	"strconv"
)

// Config contains the configuration data for an application
type Config struct {
	configMap  map[string]string
	fs         fs.FS
	currentEnv string
}

const configKey = "config"

// GetKey gets a key from the current config.
// It first checks to see if an environment variable is available, otherwise it returns a value from the map.
func (c Config) GetKey(key string) (string, error) {

	if val := os.Getenv(key); val != "" {
		return val, nil
	}

	if v, ok := c.configMap[key]; ok {
		return v, nil
	} else {
		return "", errors.New("missing value in config: " + key)
	}
}

// GetKeyAsInt gets a key from the config and converts it to an int.
func (c Config) GetKeyAsInt(key string) (int, error) {
	val, err := c.GetKey(key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New("error converting config value to int with key: " + key)
	}
	return i, nil
}

// NewConfig builds a new config. It sets the fs.FS and creates an empty config map.
func NewConfig(fSys fs.FS) Config {

	c := Config{
		configMap: make(map[string]string),
		fs:        fSys,
	}
	return c
}

// WithEnvironment set the current environment for this config.
func (c Config) WithEnvironment(currentEnv string) Config {
	c.currentEnv = currentEnv
	return c
}

// Load loads the config from a configuration file. If an environment has been specified then the config is loaded from a file called {environment}.env.
// Otherwise the config is loaded from a file called 'default.env'
func (c Config) Load() (Config, error) {

	if c.fs == nil {
		return c, errors.New("error loading config, fs was nil")
	}

	var configFile string
	if c.currentEnv != "" {
		slog.Info("found environment: " + c.currentEnv)
		configFile = fmt.Sprintf("%s.env", c.currentEnv)
	} else {
		slog.Info("environment not set, setting environment to 'default'")
		configFile = "default.env"
	}

	b, err := fs.ReadFile(c.fs, configFile)
	if err != nil {
		return c, fmt.Errorf("error loading config: %w", err)
	}

	c.configMap, err = godotenv.UnmarshalBytes(b)
	if err != nil {
		return c, fmt.Errorf("error loading config: %w", err)
	}
	slog.Info("loading config from: " + configFile)
	return c, nil
}

// ToContext adds the Config to the given context.Context
func ToContext(ctx context.Context, cfg Config) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("error, nil context")
	}

	return context.WithValue(ctx, configKey, cfg), nil
}

// FromContext gets a Config from the given context.Context
func FromContext(ctx context.Context) (Config, error) {
	if ctx == nil {
		return Config{}, errors.New("error, nil context")
	}

	val := ctx.Value(configKey)
	if val == nil {
		return Config{}, fmt.Errorf("no value found")
	}

	config, ok := val.(Config)
	if !ok {
		return Config{}, fmt.Errorf("value in context is not of type Config")
	}
	return config, nil
}
