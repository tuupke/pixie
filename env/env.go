package env

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func DurationFb(key string, fb time.Duration) time.Duration {
	durStr := String(key)
	d, err := time.ParseDuration(durStr)
	if err != nil {
		return fb
	}

	return d
}

func String(key string) string {
	return os.Getenv(key)
}

func StringFb(key, fb string) (v string) {
	v = String(key)
	if v == "" {
		v = fb
	}

	return
}

func Bool(key string) bool {
	val, err := strconv.ParseBool(os.Getenv(key))
	return val && err == nil
}

func BoolFb(key string, fb bool) bool {
	if val, err := strconv.ParseBool(os.Getenv(key)); err != nil {
		return val
	}

	return fb
}

func IntFb(key string, fb int) int {
	str := String(key)
	i, err := strconv.Atoi(str)
	if err != nil {
		return fb
	}

	return i
}

func UInt64Fb(key string, fb uint64) uint64 {
	str := String(key)

	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fb
	}

	return i
}
