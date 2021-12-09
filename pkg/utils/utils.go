package utils

import "os"

// This helper function check if the env variable is empty, and if it's empty
// then assigns the fallback / default variable defined for each variable
func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = fallback
	}
	return value
}
