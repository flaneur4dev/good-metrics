package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func EnvVar(k string, d any) (res any) {
	if val, ok := os.LookupEnv(k); ok {
		switch fmt.Sprintf("%T", d) {
		case "string":
			res = val
		case "int":
			i, err := strconv.Atoi(val)
			if err != nil {
				log.Fatal("Incorrect parameter!")
			}
			res = i
		case "bool":
			b, err := strconv.ParseBool(val)
			if err != nil {
				log.Fatal("Incorrect parameter!")
			}
			res = b
		default:
			log.Fatal("Incorrect parameter!")
		}
	} else {
		res = d
	}
	return
}
