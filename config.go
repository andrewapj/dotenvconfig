package dotenvconfig

import (
	"context"
	"errors"
	"fmt"
	"github.com/andrewapj/dotenvconfig/internal/environment"
	"github.com/andrewapj/dotenvconfig/internal/logging"
	"github.com/andrewapj/dotenvconfig/internal/parser"
	"io/fs"
	"os"
	"strconv"
)

// Config contains the configuration data for an application.
type Config struct {
	configMap map[string]string
}

// Options contains optional ways to configure an application.
type Options struct {
	// contextKey represents the key used to store/retrieve the config from a context.
	contextKey string

	// environment represents the specific environment to use.
	environment string

	// environmentKey is the key used to specify the environment to use. Takes precedence over other
	//configuration options.
	environmentKey string

	// jsonLogging indicates whether logging should be done in JSON format.
	jsonLogging bool

	// loggingEnabled determines if logging is enabled for the application.
	loggingEnabled bool
}

var contextKey = "config"
var ErrFsIsNil = errors.New("error, the FS was nil")

// NewConfig builds a new config.
func NewConfig(fSys fs.FS, opts Options) (Config, error) {

	logging.SetupLogging(opts.loggingEnabled, opts.jsonLogging)
	if opts.contextKey != "" {
		logging.Info("setting context key to " + opts.contextKey)
		contextKey = opts.contextKey
	} else {
		logging.Info("using default context key of " + contextKey)
	}

	if fSys == nil {
		logging.Error("error =, fs was nil")
		return Config{}, ErrFsIsNil
	}

	env := environment.GetEnvironment(opts.environmentKey, opts.environment)

	bytes, err := fs.ReadFile(fSys, env)
	if err != nil {
		logging.Error("error reading config file " + env)
		return Config{}, err
	}

	cfg, err := parser.Parse(bytes)
	if err != nil {
		logging.Error("error parsing config file " + env)
		return Config{}, err
	}

	return Config{cfg}, nil
}

// GetKey gets a key from the current config.
// It first checks to see if an environment variable is available, otherwise it returns a value from the map.
func (c Config) GetKey(key string) string {

	if val := os.Getenv(key); val != "" {
		return val
	}

	if v, ok := c.configMap[key]; ok {
		return v
	} else {
		return ""
	}
}

// GetKeyAsInt gets a key from the config and converts it to an int.
func (c Config) GetKeyAsInt(key string) int {
	val := c.GetKey(key)

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

// ToContext adds the Config to the given context.Context
func ToContext(ctx context.Context, cfg Config) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("error, nil context")
	}

	return context.WithValue(ctx, contextKey, cfg), nil
}

// FromContext gets a Config from the given context.Context
func FromContext(ctx context.Context) (Config, error) {
	if ctx == nil {
		return Config{}, errors.New("error, nil context")
	}

	val := ctx.Value(contextKey)
	if val == nil {
		return Config{}, fmt.Errorf("no value found")
	}

	config, ok := val.(Config)
	if !ok {
		return Config{}, fmt.Errorf("value in context is not of type Config")
	}
	return config, nil
}
