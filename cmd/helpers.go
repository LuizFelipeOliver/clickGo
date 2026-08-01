package cmd

import "os"

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func getenv(key string) string {
	return os.Getenv(key)
}
