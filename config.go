package dotenvconfig

import (
	"errors"
	"github.com/andrewapj/dotenvconfig/internal/logging"
	"github.com/andrewapj/dotenvconfig/internal/parser"
	"io/fs"
	"os"
	"strconv"
)

// Options contains optional ways to configure an application.
type Options struct {
	// JsonLogging indicates whether logging should be done in JSON format.
	JsonLogging bool

	// LoggingEnabled determines if logging is enabled for the library.
	LoggingEnabled bool
}

var ErrFsIsNil = errors.New("error, the FS was nil")
var ErrMissingKey = errors.New("error, could not find key")
var ErrConversion = errors.New("error, unable to convert value from string")

// Load reads the config from an env file and adds it to the environment.
// If an environment variable already exists then it is not overwritten.
func Load(fSys fs.FS, path string, opts Options) error {

	logging.SetupLogging(opts.LoggingEnabled, opts.JsonLogging)

	if fSys == nil {
		logging.Error("error, fs was nil")
		return ErrFsIsNil
	}

	bytes, err := fs.ReadFile(fSys, path)
	if err != nil {
		logging.Error("error reading config file " + path)
		return err
	}

	cfg, err := parser.Parse(bytes)
	if err != nil {
		logging.Error("error parsing config file " + path)
		return err
	}

	for k, v := range cfg {
		if _, ok := os.LookupEnv(k); !ok {
			err = os.Setenv(k, v)
			if err != nil {
				logging.Error("unable to set environment with key " + k + " and value " + v)
				return err
			}
		}
	}

	return nil
}

// GetKey retrieves a value from the config by key.
func GetKey(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", ErrMissingKey
	}

	return v, nil
}

// GetKeyMust retrieves a value from the config by key. It panics if the key does not exist.
func GetKeyMust(key string) string {

	v, err := GetKey(key)
	if err != nil {
		panic(err)
	}

	return v
}

// GetKeyAsInt retrieves a value from the config by key and converts the value to an int. If the key could not be found
// or the value can not be converted to an int, a zero is returned.
func GetKeyAsInt(key string) (int, error) {

	v, err := GetKey(key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, ErrConversion
	}
	return i, nil
}

// GetKeyAsIntMust retrieves a value from the config by key and converts the value to an int. If the key could not be found
// or the value can not be converted to an int, the function will panic.
func GetKeyAsIntMust(key string) int {
	v, err := GetKeyAsInt(key)
	if err != nil {
		panic(err.Error())
	}

	return v
}
