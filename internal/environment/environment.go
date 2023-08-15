package environment

import (
	"fmt"
	"github.com/andrewapj/dotenvconfig/internal/logging"
	"os"
)

const defaultEnvironment = "default"
const fileExtension = ".env"

// GetEnvironment calculates the current environment that should be used.
// If first checks an environment variable specified by environmentKey to find the environment to use.
// If not set, it uses the provided fallback environment value. If both are missing,
// it defaults the environment to 'default'.
//
// Parameters:
// - environmentKey: Key which specifies the environment to set.
// - environment: Fallback value if the environmentKey is not set.
//
// Returns:
// The environment that should be set.
func GetEnvironment(environmentKey, environment string) string {
	if val, ok := os.LookupEnv(environmentKey); ok {
		logging.Info(fmt.Sprintf("found an environment variable of %s which will be used to set the current environment",
			environmentKey))
		return buildFilename(val)
	}

	if environment != "" {
		logging.Info(fmt.Sprintf("setting current environment to %s", environment))
		return buildFilename(environment)
	}

	logging.Info("no environment set, defaulting to 'default'")
	return buildFilename(defaultEnvironment)
}

func buildFilename(env string) string {
	return env + fileExtension
}
