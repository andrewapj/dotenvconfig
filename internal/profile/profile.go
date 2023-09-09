package profile

import (
	"fmt"
	"github.com/andrewapj/dotenvconfig/internal/logging"
	"os"
)

const defaultProfile = "default"
const fileExtension = ".env"

// GetProfile determines the current profile that should be used.
// It first checks a profile variable, specified by profileKey, to find the profile to use.
// If not set, it uses the provided fallback profile value.
// If both are missing, it defaults the profile to 'default'.
//
// Parameters:
// - profileKey: Key which specifies the profile to set.
// - profile: Fallback value if the profileKey is not set.
//
// Returns:
// The profile that should be set.
func GetProfile(profileKey, profile string) string {
	if val, ok := os.LookupEnv(profileKey); ok {
		logging.Info(fmt.Sprintf("found a profile variable of %s which will be used to set the current profile",
			profileKey))
		return buildFilename(val)
	}

	if profile != "" {
		logging.Info(fmt.Sprintf("setting current profile to %s", profile))
		return buildFilename(profile)
	}

	logging.Info("no profile set, defaulting to 'default'")
	return buildFilename(defaultProfile)
}

func buildFilename(profile string) string {
	return profile + fileExtension
}
