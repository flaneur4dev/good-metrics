package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	GaugeName   = "gauge"
	CounterName = "counter"
)

func EnvVar(name string, defaultV any) (res any) {
	if val, ok := os.LookupEnv(name); ok {
		switch fmt.Sprintf("%T", defaultV) {
		case "string":
			res = val
		case "int":
			i, err := strconv.Atoi(val)
			if err != nil {
				log.Fatal("Incorrect parameter: ", val)
			}
			res = i
		case "bool":
			b, err := strconv.ParseBool(val)
			if err != nil {
				log.Fatal("Incorrect parameter: ", val)
			}
			res = b
		default:
			log.Fatal("Incorrect parameter: ", val)
		}
	} else {
		res = defaultV
	}
	return
}
