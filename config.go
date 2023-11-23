package dotenvconfig

import (
	"errors"
	"github.com/andrewapj/dotenvconfig/internal/logging"
	"github.com/andrewapj/dotenvconfig/internal/parser"
	"github.com/andrewapj/dotenvconfig/internal/profile"
	"io/fs"
	"os"
	"strconv"
)

// Options contains optional ways to configure an application.
type Options struct {
	// Profile represents the specific Profile to use.
	Profile string

	// ProfileKey is the key used to specify the Profile to use. Takes precedence over other
	//configuration options.
	ProfileKey string

	// JsonLogging indicates whether logging should be done in JSON format.
	JsonLogging bool

	// LoggingEnabled determines if logging is enabled for the application.
	LoggingEnabled bool
}

var ErrFsIsNil = errors.New("error, the FS was nil")
var ErrMissingKey = errors.New("error, could not find key")

// Load reads the config from an env file and adds it to the environment.
// If an environment variable already exists then it is not overwritten.
func Load(fSys fs.FS, opts Options) error {

	logging.SetupLogging(opts.LoggingEnabled, opts.JsonLogging)

	if fSys == nil {
		logging.Error("error, fs was nil")
		return ErrFsIsNil
	}

	p := profile.GetProfile(opts.ProfileKey, opts.Profile)

	bytes, err := fs.ReadFile(fSys, p)
	if err != nil {
		logging.Error("error reading config file " + p)
		return err
	}

	cfg, err := parser.Parse(bytes)
	if err != nil {
		logging.Error("error parsing config file " + p)
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
func GetKey(key string) string {
	return os.Getenv(key)
}

// GetKeyMust retrieves a value from the config by key. It panics if the key does not exist.
func GetKeyMust(key string) string {

	v, ok := os.LookupEnv(key)
	if !ok {
		panic(ErrMissingKey.Error())
	}
	return v
}

// GetKeyAsInt retrieves a value from the config by key and converts the value to an int. If the key could not be found
// or the value can not be converted to an int, a zero is returned.
func GetKeyAsInt(key string) int {
	val := GetKey(key)

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

// GetKeyAsIntMust retrieves a value from the config by key and converts the value to an int. If the key could not be found
// or the value can not be converted to an int, the function will panic.
func GetKeyAsIntMust(key string) int {
	val := GetKey(key)

	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err.Error())
	}
	return i
}
