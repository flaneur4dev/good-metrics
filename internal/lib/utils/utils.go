package utils

import (
	"log"
	"os"
	"strconv"
)

const (
	GaugeName   = "gauge"
	CounterName = "counter"
)

func StringEnv(name, defaultV string) (res string) {
	if val, ok := os.LookupEnv(name); ok {
		res = val
	} else {
		res = defaultV
	}
	return
}

func BoolEnv(name string, defaultV bool) (res bool) {
	if val, ok := os.LookupEnv(name); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatal("Incorrect parameter: ", val)
		}
		res = b
	} else {
		res = defaultV
	}
	return
}
