package common

import "syscall"

func EnvString(key, fallback string) string {
	if res, ok := syscall.Getenv(key); ok {
		return res
	}
	return fallback
}
