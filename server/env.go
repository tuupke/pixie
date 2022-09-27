package main

import (
	"os"
	"strconv"
	"time"
)

func envDurationFb(key string, fb time.Duration) time.Duration {
	durStr := envString(key)
	d, err := time.ParseDuration(durStr)
	if err != nil {
		return fb
	}

	return d
}

func envString(key string) string {
	return os.Getenv(key)
}

func envStringFb(key, fb string) (v string) {
	v = envString(key)
	if v == "" {
		v = fb
	}

	return
}

func envBool(key string) bool {
	val, err := strconv.ParseBool(os.Getenv(key))
	return val && err != nil
}

func envBoolFb(key string, fb bool) bool {
	if val, err := strconv.ParseBool(os.Getenv(key)); err == nil {
		return val
	}

	return fb
}

func envIntFb(key string, fb int) int {
	str := envString(key)
	i, err := strconv.Atoi(str)
	if err != nil {
		return fb
	}

	return i
}

func envUInt64Fb(key string, fb uint64) uint64 {
	str := envString(key)

	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fb
	}

	return i
}
